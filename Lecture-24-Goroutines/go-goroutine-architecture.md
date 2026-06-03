# Go Goroutine Architecture — Reference Notes

A bottom-to-top reference on how Go runs goroutines: the runtime scheduler, the
G-M-P model, and the mechanisms that let one program run millions of goroutines
on a handful of OS threads.

> **The one idea everything hangs on:** `go f()` does *not* ask the OS for a
> thread. It hands work to Go's own scheduler, which runs in user space inside
> your process. Go multiplexes many cheap goroutines onto a few expensive OS
> threads. This is **M:N scheduling** — M goroutines onto N threads.

**Analogy used throughout:** a restaurant kitchen.
`G` = an order ticket / dish (the work) · `M` = a chef (an OS thread, the only
thing that can actually cook) · `P` = a cooking station with its own ticket rail
and prep stock (there are only `GOMAXPROCS` of them).

---

## 1. Prerequisites (the layer below Go)

| Concept | What it is | Why it matters for Go |
|---|---|---|
| **CPU core** | The hardware that physically executes instructions, one stream at a time. | True parallelism is capped at the number of cores. |
| **Process** | A running program plus its isolated address space, memory, and file descriptors. | All Gs, Ms, Ps, and queues live inside one process. |
| **OS thread** | The unit the OS scheduler places on a core. Heavy: ~1 MB stack + kernel bookkeeping. | An `M` *is* an OS thread. |
| **Context switch** | OS saving one thread's state and loading another's to time-slice them. | Thread switches go through the kernel and are expensive. Goroutine switches stay in user space and are cheap. |
| **User vs kernel mode** | Two CPU privilege levels; a **syscall** traps from user to kernel and back. | Go's scheduler runs in user mode, so scheduling never pays the syscall toll. |
| **Concurrency vs parallelism** | Concurrency = tasks *in progress* at overlapping times. Parallelism = tasks running at the *same instant*. | Goroutines give concurrency always; parallelism only up to `GOMAXPROCS`. |
| **Blocking vs non-blocking I/O** | Blocking = thread sits idle until data is ready. Non-blocking = "tell me when ready," do other work meanwhile. | The basis of the network poller. |

---

## 2. Why goroutines exist

Creating OS threads is expensive — roughly **1 MB** of reserved stack each, plus
kernel scheduling overhead and costly context switches. Ten thousand threads is
already gigabytes of memory and heavy switching cost.

A goroutine starts at about **2 KB** of stack and switches in user space without
a kernel trip. That is the entire reason you can run **hundreds of thousands to
millions** of goroutines where you could never have that many threads.

---

## 3. The three core structures: G, M, P

### G — Goroutine (the work)
A `G` is *just a struct*. It does nothing on its own; it is a description of work
plus the saved state of that work. It holds:

- **Stack** — where its local variables and call frames live (starts at 2 KB).
- **Program counter (PC)** — which instruction it is paused at, so it can resume
  exactly where it stopped.
- **Status** — runnable / running / waiting / syscall / dead.
- **Goroutine ID + scheduler metadata.**

### M — Machine (the worker)
An `M` is an **actual OS thread** — the only thing that can physically execute
code, because the CPU only ever runs threads. An idle M parks (sleeps).

### P — Processor (the permission + toolkit)
A `P` is **not** a CPU and **not** a thread. It is a *scheduling context*: a
license to run Go code, bundled with private, lock-free resources. Go creates
exactly `GOMAXPROCS` of them.

> **The binding rule:** to run a goroutine, a chef must stand at a station.
> **M + P + G = execution.** No free P → no execution, even with idle Ms and
> pending Gs. This is why `GOMAXPROCS` caps simultaneous execution.

**What a P owns:**

- **Local run queue** — a 256-slot ring buffer of runnable Gs. Lock-free for its
  owning M; only work-stealing from another P uses atomics.
- **`runnext` slot** — a single-G fast-path slot holding the *most recently
  readied* goroutine, run before the main queue. Gives cache locality for
  channel hand-offs (a just-woken receiver runs immediately, while the sender's
  data is still warm in cache).
- **mcache** — a per-P memory allocator cache. Small allocations are served
  here without locks, so goroutines on different Ps allocate in parallel.
- **Per-P timer heap** (since Go 1.14) and **per-P GC work buffers.**

The recurring pattern: *take something that used to be global and contended, and
give every P its own private copy.*

---

## 4. The queues

- **Local run queue** — each P's own 256-slot queue. The common case
  ("run my next G") touches only this, no global lock.
- **Global run queue** — overflow store. When a local queue fills, ~half is
  flushed here. With a million goroutines, most live here. Protected by a lock,
  so it's the slower path.

---

## 5. The scheduler loop

Every M, after finishing or parking a G, runs this loop forever (simplified):

```
for {
    find a runnable goroutine
    run it
    if it blocks            -> save its state, park it
    if my queue is empty    -> steal work from another P
    if network I/O is ready -> wake the waiting goroutine
    if a goroutine ran too long -> preempt it
}
```

**Where it looks for work, in order:**
1. Its own local run queue (cheapest).
2. **Every 61st scheduling tick**, check the global queue first instead — this
   forced peek prevents the global queue from starving when local queues are
   always busy.
3. The network poller.
4. **Steal** half of a random other P's local queue.

---

## 6. Work stealing

Keeps every core busy. When a P's local queue empties, its M doesn't idle — it
steals **roughly half** of a busy P's queue.

```
Before:  P1: G1 G2 G3 G4 G5 G6 G7 G8   |  P2: (empty)
After:   P1: G1 G2 G3 G4               |  P2: G5 G6 G7 G8
```

Both stations now cook in parallel instead of one chef working while another
stands idle.

---

## 7. Context switching (goroutine-level)

To swap one goroutine out for another, the scheduler saves only the **PC, stack
pointer, and a few registers** into the G struct — all in user space, no kernel
trip. This is an order of magnitude cheaper than an OS thread switch, which is
exactly what makes constant switching among millions of goroutines affordable.

---

## 8. Preemption

**Problem:** a goroutine that never yields (`for {}` with no function calls,
channel ops, or allocations) would hog its station forever, and could stall GC
(which needs every goroutine to reach a safe point).

- **Before Go 1.14 — cooperative preemption:** the runtime could only preempt at
  function-call boundaries (a check in each function prologue). A tight loop with
  no calls had no such check → it could freeze the program.
- **Since Go 1.14 — asynchronous preemption:** the supervisor notices a G has run
  too long (~10 ms), sends the M an OS signal (`SIGURG` on Unix), and the signal
  handler safely parks the G and reschedules. This is why `for {}` no longer
  freezes a modern Go program.

---

## 9. sysmon (the supervisor)

A special M that runs **without a P** (it supervises, it doesn't cook). It loops
forever, sleeping ~20 µs to ~10 ms between rounds, and handles the edge cases:
triggering async preemption, retaking Ps from long syscalls, poking the network
poller if neglected, and forcing periodic GC.

---

## 10. Network poller (netpoll)

The piece that makes Go feel different from Node.js. When a goroutine does
network I/O with no data ready:

1. Go registers the socket's file descriptor with the OS readiness mechanism —
   **epoll** (Linux), **kqueue** (macOS/BSD), **IOCP** (Windows).
2. The G is **parked** (status: waiting), and **the M and P are freed** to run
   other goroutines.
3. When the OS reports the socket readable, the poller flips the G back to
   **runnable** and puts it on a run queue.
4. The G resumes **on the exact line after the read**, with all locals intact.

A waiting connection therefore costs a parked goroutine (a few KB), not an OS
thread. This is how one Go server holds **a million idle connections with a
handful of threads** — the engine under `net/http`, gRPC, DB drivers, etc.

---

## 11. Blocking syscalls (different from network I/O!)

The poller only helps with things the OS can watch for *readiness* (mostly
sockets). A genuinely blocking syscall (`os.Open`, a synchronous file read, a CGo
call) **pins the OS thread**. So Go does a **handoff**:

- On entering the syscall, the runtime **detaches the P** from that M.
- The M dives into the syscall carrying its G (both stuck for the duration).
- The freed P is handed to another M (woken or newly created), so that station
  keeps cooking.
- When the syscall returns, the M tries to reacquire a P; if none is free, it
  parks its G on the global queue and goes to sleep.

**Side effect:** many simultaneous blocking syscalls force Go to spin up many OS
threads — fine usually, but it's why heavy synchronous file/CGo work can quietly
grow the thread count.

---

## 12. Goroutine lifecycle

```
go worker()
   │
   ▼
runnable ──(M+P picks it up)──▶ running
   ▲                              │
   │                              ├──(channel / sleep / netpoll / mutex)──▶ waiting
   │                              └──(blocking syscall)────────────────────▶ syscall
   │                                                                          │
   └──────────────(unblocked: poller / timer / sender wakes it)──────────────┘
                                  │
                                  ▼
                                 dead  (G struct recycled into a free pool, not freed)
```

Most production debugging is the question: **"why is this goroutine stuck in
*waiting*, and who was supposed to wake it?"** (a receiver that never came, a
lock never released, a context never cancelled).

---

## 13. Channel internals

A channel is a struct holding a ring buffer, a lock, a queue of waiting senders,
and a queue of waiting receivers.

- **Unbuffered `ch <- x`:** if a receiver is already parked, the value is copied
  straight to it and that receiver's G is woken — no buffer involved. If no
  receiver, the sender's G parks on the send-wait queue until one appears.
- **Buffered channel:** send copies into the buffer and continues while there's
  room; blocks only when full.

Key point: "blocking on a channel" never means a busy-wait or a wasted thread —
it's the **same parking mechanism** as sleep and netpoll. The G goes to
*waiting*, the M is freed.

Questions you should be able to answer: why does an unbuffered channel block? why
can buffering improve throughput? why can channels become a bottleneck (lock
contention on a hot channel)?

---

## 14. Stack growth

Each goroutine starts with a 2 KB stack. When a call would overflow it, a check
in the function prologue catches it; the runtime allocates a **larger stack**
(doubling: 2 → 4 → 8 → 16 KB …), copies the contents over, fixes up pointers, and
continues. Stacks can shrink during GC. These **contiguous, growable stacks** are
the final ingredient that makes millions of goroutines feasible — you only pay
for a big stack if a goroutine actually goes deep.

---

## 15. Garbage collector (brief)

Go's GC is concurrent — it runs alongside your goroutines using special worker
goroutines (mark workers, background sweepers, and GC-assist done by allocating
goroutines). Deeper internals worth knowing later: tri-color marking, write
barriers, stop-the-world (STW) pauses, and GC pacing.

---

## 16. GOMAXPROCS = the number of Ps

`GOMAXPROCS` does **not** set the number of threads — it sets the number of
**Ps**, the real ceiling on how many goroutines execute Go code at the same
instant. Read it with `runtime.GOMAXPROCS(0)`.

- **Historic default:** the host's logical CPU count.
- **Since Go 1.25 (Linux):** **container-aware** — if running in a container with
  a CPU limit, it defaults to that limit when it's lower than the core count, and
  re-checks roughly every 30 seconds. A fractional limit is rounded down (4.5
  cores → 4). This avoids the old bug where a 2-CPU pod on a 16-core node assumed
  16 CPUs, over-scheduled, and got throttled (hurting tail latency).
- **Overrides / opt-outs:** setting the `GOMAXPROCS` env var or calling
  `runtime.GOMAXPROCS(n)` disables the cgroup-aware default; so does
  `GODEBUG=containermaxprocs=0`. **Watch out:** a stale hardcoded
  `GOMAXPROCS=8` in a deployment manifest silently disables the new behavior.

**Mental model:** Ms = how many chefs (can balloon during blocking I/O) ·
Gs = how many dishes ordered (millions, mostly parked) · **Ps = how many stoves**
— and that number bounds CPU-bound throughput.

---

## 17. The complete walkthrough

When you write `go worker()`:

1. Runtime grabs a G struct (reused from a free list — usually no allocation).
2. Gives it a 2 KB stack, sets PC to `worker`'s start, stashes arguments.
3. Marks it runnable, pushes it onto the current P's local run queue.
4. Returns immediately (the goroutine does **not** run yet).
5. Later, an M holding a P pops it and runs it.
6. If it blocks on network I/O → parked, M and P freed, socket registered with
   the poller.
7. Data arrives → poller marks it runnable, back onto a run queue.
8. An M runs it again; it resumes from the exact paused line.
9. `worker` returns → G becomes dead, struct recycled.

---

## 18. Go vs Node.js

| | Node.js | Go |
|---|---|---|
| Model | 1 JS thread + 1 event loop | many goroutines, many threads, many Ps |
| On I/O completion | a **callback** is pushed to the event queue | the **goroutine** is woken and resumed |
| Your code style | callbacks / `async`-`await` (control inverted) | straight-line code that *looks* synchronous |
| Resume behavior | run a callback later | continue from the exact next instruction, locals intact |

Node says *"when ready, run this callback."* Go says *"when ready, wake this
goroutine and continue mid-sentence."* That resume-where-you-paused behavior is
why goroutine code reads more naturally while scaling just as well.

---

## 19. Further layers (next study targets)

Sitting on top of the scheduler foundation above:

- **Memory model** — happens-before, visibility, why `x` can read as stale
  without synchronization.
- **Synchronization primitives** — `sync.Mutex`, `RWMutex`, `WaitGroup`, `Once`,
  `Cond`, `Map`; when to use a channel vs a mutex vs an atomic.
- **Escape analysis** — when the compiler moves a value from stack to heap.
- **Heap vs stack** — allocation cost and GC pressure.
- **`select` internals** — how Go chooses among ready cases (randomized for
  fairness), what happens when none are ready.
- **Mutex internals** — fast path, slow path, spinning, parking, starvation mode.
- **Atomics** — CAS (compare-and-swap), memory barriers, lock-free patterns.
- **GC internals** — tri-color marking, write barriers, STW, pacing.
- **Goroutine leaks** — a goroutine blocked forever on a channel nobody sends to;
  a common production failure.
- **`context` package** — cancellation, deadlines, request scoping, propagation.
- **Diagnostics & tracing** — `runtime.NumGoroutine()`, `pprof`, `go tool trace`,
  `GODEBUG=schedtrace`.

---

## 20. Tools to see it live

```bash
# Print scheduler state (P count, run-queue depths) every 1000 ms:
GODEBUG=schedtrace=1000 ./your-program

# Add per-P detail:
GODEBUG=schedtrace=1000,scheddetail=1 ./your-program

# Visual timeline of goroutines moving across Ps, blocking, and GC:
go tool trace trace.out

# CPU / memory / blocking profiles:
go tool pprof
```

```bash
# Sanity-check parallelism vs the CPU your container actually granted:
cat /proc/cpuinfo | grep -c processor   # logical cores the host sees
# then compare against runtime.GOMAXPROCS(0) inside the app
```

---

## Suggested learning order

1. Process vs thread → core → context switching → user/kernel mode
2. Goroutine lifecycle
3. G-M-P
4. Channel internals
5. `select`
6. `Mutex` / `RWMutex`
7. `context`
8. Atomics
9. Escape analysis
10. Heap vs stack
11. Network poller
12. Garbage collector
13. Scheduler tracing
14. Runtime profiling

---

*Notes compiled as a personal reference. Version-specific details (async
preemption, container-aware `GOMAXPROCS`) reflect Go 1.14 and Go 1.25
respectively.*
