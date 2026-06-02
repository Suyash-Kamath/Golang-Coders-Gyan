const scenarios = [
  {
    id: "basic",
    title: "Basic goroutine scheduling",
    code: [
      "package main",
      "",
      "import (",
      '    "fmt"',
      '    "time"',
      ")",
      "",
      "func worker() {",
      '    fmt.Println("worker: start")',
      "    time.Sleep(10 * time.Millisecond)",
      '    fmt.Println("worker: done")',
      "}",
      "",
      "func main() {",
      '    fmt.Println("main: start")',
      "    go worker()",
      '    fmt.Println("main: continues")',
      "    time.Sleep(30 * time.Millisecond)",
      "}",
    ],
    steps: [
      {
        title: "Program starts",
        time: 0,
        line: 14,
        concept: "main is also a goroutine",
        explanation:
          "Go starts the program by creating the main goroutine. A goroutine is a lightweight user-space task managed by the Go runtime, not a full operating-system thread.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [],
        chips: ["G = goroutine", "M = OS thread", "P = processor token"],
      },
      {
        title: "go statement creates G2",
        time: 1,
        line: 16,
        concept: "go schedules work",
        explanation:
          "The go keyword does not call worker immediately on the current stack. It creates a new goroutine, gives it a small stack, and puts it in a run queue so the scheduler can run it.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running"] },
          { id: "G2", role: "worker", state: "ready", history: ["", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [],
        chips: ["go worker()", "local run queue", "no blocking yet"],
      },
      {
        title: "main continues",
        time: 2,
        line: 17,
        concept: "concurrency, not always parallelism",
        explanation:
          "The main goroutine keeps executing. G2 is ready, but it only runs when a P and M pick it up. With multiple Ps it may run in parallel; with one P it is interleaved.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running"] },
          { id: "G2", role: "worker", state: "ready", history: ["", "ready", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [],
        chips: ["ready does not mean running", "scheduler decides"],
      },
      {
        title: "G2 runs on another processor",
        time: 3,
        line: 9,
        concept: "P lets M execute G",
        explanation:
          "A processor token P owns runnable work. An OS thread M must hold a P to execute Go code. Here M1 with P1 starts running the worker goroutine.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "running"] },
          { id: "G2", role: "worker", state: "running", history: ["", "ready", "ready", "running"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "G2", queue: [] }, global: [] },
        channels: [],
        chips: ["parallel when multiple Ps exist", "GOMAXPROCS controls Ps"],
      },
      {
        title: "worker sleeps",
        time: 10,
        line: 10,
        concept: "blocking parks goroutines",
        explanation:
          "time.Sleep parks G2 in the runtime timer system. The OS thread is not wasted waiting for the sleep; the scheduler can run other goroutines on that M.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "running", "running"] },
          { id: "G2", role: "worker", state: "blocked", history: ["", "ready", "ready", "running", "blocked"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "timer heap", kind: "timer", slots: ["G2 wakes at 20ms"], senders: [], receivers: [] }],
        chips: ["blocked", "timer wait", "M can do other work"],
      },
      {
        title: "main sleeps, worker wakes",
        time: 20,
        line: 18,
        concept: "wake puts G back in queue",
        explanation:
          "When the timer fires, G2 becomes runnable again. It re-enters a run queue and waits for the scheduler to assign it to an available M/P pair.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "running", "running", "blocked"] },
          { id: "G2", role: "worker", state: "ready", history: ["", "ready", "ready", "running", "blocked", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "timer heap", kind: "timer", slots: ["G1 wakes at 30ms"], senders: [], receivers: [] }],
        chips: ["wakeup", "runnable again"],
      },
      {
        title: "worker finishes",
        time: 21,
        line: 11,
        concept: "goroutine exits",
        explanation:
          "After worker prints its final line, G2 returns from its function and exits. Its stack and bookkeeping can be cleaned up by the runtime.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "running", "running", "blocked", "blocked"] },
          { id: "G2", role: "worker", state: "done", history: ["", "ready", "ready", "running", "blocked", "ready", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "timer heap", kind: "timer", slots: ["G1 wakes at 30ms"], senders: [], receivers: [] }],
        chips: ["return exits goroutine", "no join automatically"],
      },
      {
        title: "main exits",
        time: 30,
        line: 19,
        concept: "process ends with main",
        explanation:
          "When the main goroutine returns, the Go process exits. Any still-running goroutines would be stopped, which is why real programs use channels, WaitGroup, or context for coordination.",
        goroutines: [
          { id: "G1", role: "main", state: "done", history: ["running", "running", "running", "running", "running", "blocked", "blocked", "done"] },
          { id: "G2", role: "worker", state: "done", history: ["", "ready", "ready", "running", "blocked", "ready", "done", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [],
        chips: ["main controls process lifetime", "coordinate before exit"],
      },
    ],
  },
  {
    id: "channel",
    title: "Unbuffered channel handoff",
    code: [
      "package main",
      "",
      "import \"fmt\"",
      "",
      "func main() {",
      "    ch := make(chan string)",
      "",
      "    go func() {",
      '        ch <- "done"',
      "    }()",
      "",
      "    msg := <-ch",
      "    fmt.Println(msg)",
      "}",
    ],
    steps: [
      {
        title: "main creates channel",
        time: 0,
        line: 6,
        concept: "channel object",
        explanation:
          "A channel is a runtime object. An unbuffered channel has no storage slot for values; send and receive must meet each other.",
        goroutines: [{ id: "G1", role: "main", state: "running", history: ["running"] }],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "ch", kind: "unbuffered", slots: [], senders: [], receivers: [] }],
        chips: ["make(chan string)", "capacity 0"],
      },
      {
        title: "sender goroutine is created",
        time: 1,
        line: 8,
        concept: "new goroutine becomes ready",
        explanation:
          "The anonymous function becomes G2 and is placed in a run queue. It will attempt to send when scheduled.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running"] },
          { id: "G2", role: "sender", state: "ready", history: ["", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "ch", kind: "unbuffered", slots: [], senders: [], receivers: [] }],
        chips: ["go func", "ready queue"],
      },
      {
        title: "main tries to receive",
        time: 2,
        line: 12,
        concept: "receive blocks",
        explanation:
          "G1 reaches <-ch before a sender has completed the send. Because the channel is unbuffered, G1 parks in the channel receiver wait queue.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "blocked"] },
          { id: "G2", role: "sender", state: "ready", history: ["", "ready", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "ch", kind: "unbuffered", slots: [], senders: [], receivers: ["G1 waiting for value"] }],
        chips: ["park receiver", "no busy waiting"],
      },
      {
        title: "sender runs",
        time: 3,
        line: 9,
        concept: "send finds receiver",
        explanation:
          "G2 runs and sends the string. The runtime sees G1 already waiting, copies the value directly to G1, and makes G1 runnable again.",
        goroutines: [
          { id: "G1", role: "main", state: "ready", history: ["running", "running", "blocked", "ready"] },
          { id: "G2", role: "sender", state: "running", history: ["", "ready", "ready", "running"] },
        ],
        scheduler: { p0: { m: "M0", running: "G2", queue: ["G1"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "ch", kind: "unbuffered", slots: ['"done" handed off'], senders: [], receivers: [] }],
        chips: ["direct handoff", "receiver wakes"],
      },
      {
        title: "sender exits",
        time: 4,
        line: 10,
        concept: "send complete",
        explanation:
          "After the value is handed off, G2 returns and exits. The receiver is ready to continue from the receive expression.",
        goroutines: [
          { id: "G1", role: "main", state: "ready", history: ["running", "running", "blocked", "ready", "ready"] },
          { id: "G2", role: "sender", state: "done", history: ["", "ready", "ready", "running", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: ["G1"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "ch", kind: "unbuffered", slots: [], senders: [], receivers: [] }],
        chips: ["sender done", "receiver queued"],
      },
      {
        title: "main prints value",
        time: 5,
        line: 13,
        concept: "coordination completed",
        explanation:
          "G1 resumes with msg set to the received value. Channels coordinate both data movement and timing between goroutines.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "blocked", "ready", "ready", "running"] },
          { id: "G2", role: "sender", state: "done", history: ["", "ready", "ready", "running", "done", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "ch", kind: "unbuffered", slots: [], senders: [], receivers: [] }],
        chips: ["data + synchronization", "fmt.Println"],
      },
    ],
  },
  {
    id: "buffered",
    title: "Buffered channel and backpressure",
    code: [
      "package main",
      "",
      "func main() {",
      "    jobs := make(chan int, 2)",
      "",
      "    jobs <- 1",
      "    jobs <- 2",
      "",
      "    go func() {",
      "        jobs <- 3",
      "    }()",
      "",
      "    <-jobs",
      "    <-jobs",
      "    <-jobs",
      "}",
    ],
    steps: [
      {
        title: "buffered channel created",
        time: 0,
        line: 4,
        concept: "channel capacity",
        explanation:
          "A buffered channel has storage inside the channel object. Sends can complete without a receiver until the buffer is full.",
        goroutines: [{ id: "G1", role: "main", state: "running", history: ["running"] }],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "jobs", kind: "buffered cap 2", slots: ["empty", "empty"], senders: [], receivers: [] }],
        chips: ["capacity 2", "buffer storage"],
      },
      {
        title: "first send fills one slot",
        time: 1,
        line: 6,
        concept: "send without receiver",
        explanation:
          "jobs <- 1 completes immediately because the channel has free buffer space. The value is copied into the channel buffer.",
        goroutines: [{ id: "G1", role: "main", state: "running", history: ["running", "running"] }],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "jobs", kind: "buffered cap 2", slots: ["1", "empty"], senders: [], receivers: [] }],
        chips: ["non-blocking send", "space available"],
      },
      {
        title: "second send fills buffer",
        time: 2,
        line: 7,
        concept: "buffer is full",
        explanation:
          "The second send also completes. Now the channel buffer is full, so any future send must wait until a receive removes a value.",
        goroutines: [{ id: "G1", role: "main", state: "running", history: ["running", "running", "running"] }],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "jobs", kind: "buffered cap 2", slots: ["1", "2"], senders: [], receivers: [] }],
        chips: ["full buffer", "backpressure begins"],
      },
      {
        title: "third sender blocks",
        time: 3,
        line: 10,
        concept: "backpressure",
        explanation:
          "G2 attempts jobs <- 3 while the buffer is full. The runtime parks G2 in the channel sender wait queue instead of spinning.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "running"] },
          { id: "G2", role: "producer", state: "blocked", history: ["", "", "", "blocked"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "jobs", kind: "buffered cap 2", slots: ["1", "2"], senders: ["G2 wants to send 3"], receivers: [] }],
        chips: ["sender wait queue", "buffer full"],
      },
      {
        title: "receive frees a slot",
        time: 4,
        line: 13,
        concept: "receive drains buffer",
        explanation:
          "G1 receives the oldest buffered value. That frees a slot, so the waiting send from G2 can complete and G2 becomes runnable again.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "running", "running"] },
          { id: "G2", role: "producer", state: "ready", history: ["", "", "", "blocked", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "jobs", kind: "buffered cap 2", slots: ["2", "3"], senders: [], receivers: [] }],
        chips: ["FIFO receive", "blocked sender wakes"],
      },
      {
        title: "remaining values are received",
        time: 5,
        line: 15,
        concept: "buffer drains",
        explanation:
          "The main goroutine receives the remaining values. Buffered channels smooth bursts, but they do not remove the need for coordination.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "running", "running", "running"] },
          { id: "G2", role: "producer", state: "done", history: ["", "", "", "blocked", "ready", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "jobs", kind: "buffered cap 2", slots: ["empty", "empty"], senders: [], receivers: [] }],
        chips: ["buffer empty", "producer done"],
      },
    ],
  },
  {
    id: "waitgroup",
    title: "sync.WaitGroup coordination",
    code: [
      "package main",
      "",
      "import (",
      '    "fmt"',
      '    "sync"',
      ")",
      "",
      "func main() {",
      "    var wg sync.WaitGroup",
      "",
      "    for i := 1; i <= 2; i++ {",
      "        wg.Add(1)",
      "        go func(id int) {",
      "            defer wg.Done()",
      '            fmt.Println("worker", id)',
      "        }(i)",
      "    }",
      "",
      "    wg.Wait()",
      '    fmt.Println("all done")',
      "}",
    ],
    steps: [
      {
        title: "WaitGroup created",
        time: 0,
        line: 9,
        concept: "counting semaphore",
        explanation:
          "A WaitGroup is a tiny counter inside the runtime. Add increments it, Done decrements it, and Wait blocks the caller until the counter reaches zero.",
        goroutines: [{ id: "G1", role: "main", state: "running", history: ["running"] }],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "wg", kind: "WaitGroup", slots: ["counter: 0"], senders: [], receivers: [] }],
        chips: ["sync.WaitGroup", "counter starts at 0"],
      },
      {
        title: "Add(1) then spawn worker 1",
        time: 1,
        line: 13,
        concept: "Add before go",
        explanation:
          "wg.Add(1) raises the counter to 1, then go func creates G2. Add must run before go to avoid a race where Wait sees zero before the goroutine has registered itself.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running"] },
          { id: "G2", role: "worker 1", state: "ready", history: ["", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "wg", kind: "WaitGroup", slots: ["counter: 1"], senders: [], receivers: [] }],
        chips: ["wg.Add(1)", "go func()"],
      },
      {
        title: "Add(1) then spawn worker 2",
        time: 2,
        line: 13,
        concept: "two pending workers",
        explanation:
          "Second iteration of the loop. Counter rises to 2 and G3 enters a run queue. Both workers are eligible to run on any P.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running"] },
          { id: "G2", role: "worker 1", state: "ready", history: ["", "ready", "ready"] },
          { id: "G3", role: "worker 2", state: "ready", history: ["", "", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2", "G3"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "wg", kind: "WaitGroup", slots: ["counter: 2"], senders: [], receivers: [] }],
        chips: ["counter: 2", "two ready workers"],
      },
      {
        title: "main blocks on Wait",
        time: 3,
        line: 19,
        concept: "Wait parks the caller",
        explanation:
          "wg.Wait sees counter > 0 and parks G1 on the WaitGroup. Now the scheduler is free to pick up G2 and G3 on idle Ms.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "blocked"] },
          { id: "G2", role: "worker 1", state: "running", history: ["", "ready", "ready", "running"] },
          { id: "G3", role: "worker 2", state: "ready", history: ["", "", "ready", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G2", queue: [] }, p1: { m: "M1", running: "G3", queue: [] }, global: [] },
        channels: [{ name: "wg", kind: "WaitGroup", slots: ["counter: 2"], senders: [], receivers: ["G1 waiting"] }],
        chips: ["main parked", "workers picked up"],
      },
      {
        title: "worker 1 prints and calls Done",
        time: 4,
        line: 14,
        concept: "defer wg.Done",
        explanation:
          "G2 finishes its print and the deferred Done runs. Counter drops to 1. Not zero yet, so G1 stays parked.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "blocked", "blocked"] },
          { id: "G2", role: "worker 1", state: "done", history: ["", "ready", "ready", "running", "done"] },
          { id: "G3", role: "worker 2", state: "running", history: ["", "", "ready", "ready", "running"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: [] }, p1: { m: "M1", running: "G3", queue: [] }, global: [] },
        channels: [{ name: "wg", kind: "WaitGroup", slots: ["counter: 1"], senders: [], receivers: ["G1 waiting"] }],
        chips: ["wg.Done()", "counter: 1"],
      },
      {
        title: "worker 2 finishes, main wakes",
        time: 5,
        line: 14,
        concept: "counter hits zero",
        explanation:
          "G3 prints and calls Done. Counter hits zero, so the runtime releases every goroutine parked on this WaitGroup. G1 becomes runnable again.",
        goroutines: [
          { id: "G1", role: "main", state: "ready", history: ["running", "running", "running", "blocked", "blocked", "ready"] },
          { id: "G2", role: "worker 1", state: "done", history: ["", "ready", "ready", "running", "done", "done"] },
          { id: "G3", role: "worker 2", state: "done", history: ["", "", "ready", "ready", "running", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: ["G1"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "wg", kind: "WaitGroup", slots: ["counter: 0"], senders: [], receivers: [] }],
        chips: ["all Done", "Wait returns"],
      },
      {
        title: "main resumes and finishes",
        time: 6,
        line: 20,
        concept: "coordination complete",
        explanation:
          "G1 picks up where Wait returned and prints the final line. WaitGroup is the right tool when you only need to know that N goroutines are finished, not the values they produced.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "blocked", "blocked", "ready", "running"] },
          { id: "G2", role: "worker 1", state: "done", history: ["", "ready", "ready", "running", "done", "done", "done"] },
          { id: "G3", role: "worker 2", state: "done", history: ["", "", "ready", "ready", "running", "done", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "wg", kind: "WaitGroup", slots: ["counter: 0"], senders: [], receivers: [] }],
        chips: ["main resumes", "use channels for values"],
      },
    ],
  },
  {
    id: "select",
    title: "select with multiple channels",
    code: [
      "package main",
      "",
      "import (",
      '    "fmt"',
      '    "time"',
      ")",
      "",
      "func main() {",
      "    a := make(chan string)",
      "    b := make(chan string)",
      "",
      "    go func() {",
      "        time.Sleep(20 * time.Millisecond)",
      '        a <- "from a"',
      "    }()",
      "    go func() {",
      "        time.Sleep(10 * time.Millisecond)",
      '        b <- "from b"',
      "    }()",
      "",
      "    select {",
      "    case msg := <-a:",
      "        fmt.Println(msg)",
      "    case msg := <-b:",
      "        fmt.Println(msg)",
      "    }",
      "}",
    ],
    steps: [
      {
        title: "channels created",
        time: 0,
        line: 10,
        concept: "two unbuffered channels",
        explanation:
          "main makes two channels. No senders or receivers exist yet. Both will be used as inputs to a single select.",
        goroutines: [{ id: "G1", role: "main", state: "running", history: ["running"] }],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "a", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "b", kind: "unbuffered", slots: [], senders: [], receivers: [] },
        ],
        chips: ["chan string", "no buffer"],
      },
      {
        title: "two producers spawned",
        time: 1,
        line: 16,
        concept: "racing producers",
        explanation:
          "G2 will eventually send on a after 20ms. G3 will send on b after 10ms. G3 will reach its send first.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running"] },
          { id: "G2", role: "a sender", state: "ready", history: ["", "ready"] },
          { id: "G3", role: "b sender", state: "ready", history: ["", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2", "G3"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "a", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "b", kind: "unbuffered", slots: [], senders: [], receivers: [] },
        ],
        chips: ["two goroutines", "different timers"],
      },
      {
        title: "producers go to sleep",
        time: 2,
        line: 17,
        concept: "timer parks",
        explanation:
          "Both senders call time.Sleep and park in the runtime timer system. The Ms are free for other work.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running"] },
          { id: "G2", role: "a sender", state: "blocked", history: ["", "ready", "blocked"] },
          { id: "G3", role: "b sender", state: "blocked", history: ["", "ready", "blocked"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "a", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "b", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "timer heap", kind: "timer", slots: ["G3 wakes at 10ms", "G2 wakes at 20ms"], senders: [], receivers: [] },
        ],
        chips: ["time.Sleep", "park on timer"],
      },
      {
        title: "main enters select",
        time: 3,
        line: 21,
        concept: "select parks on every case",
        explanation:
          "select with no ready case parks the caller on every channel listed. G1 is registered as a receiver on both a and b at the same time.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "blocked"] },
          { id: "G2", role: "a sender", state: "blocked", history: ["", "ready", "blocked", "blocked"] },
          { id: "G3", role: "b sender", state: "blocked", history: ["", "ready", "blocked", "blocked"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "a", kind: "unbuffered", slots: [], senders: [], receivers: ["G1 select"] },
          { name: "b", kind: "unbuffered", slots: [], senders: [], receivers: ["G1 select"] },
          { name: "timer heap", kind: "timer", slots: ["G3 wakes at 10ms", "G2 wakes at 20ms"], senders: [], receivers: [] },
        ],
        chips: ["queued on a and b", "no busy wait"],
      },
      {
        title: "G3 timer fires",
        time: 10,
        line: 18,
        concept: "first sender wins",
        explanation:
          "G3 wakes after 10ms and runs its send on b. The runtime sees G1 already waiting on b through the select, dequeues G1 from BOTH a and b, and hands the value over.",
        goroutines: [
          { id: "G1", role: "main", state: "ready", history: ["running", "running", "running", "blocked", "ready"] },
          { id: "G2", role: "a sender", state: "blocked", history: ["", "ready", "blocked", "blocked", "blocked"] },
          { id: "G3", role: "b sender", state: "running", history: ["", "ready", "blocked", "blocked", "running"] },
        ],
        scheduler: { p0: { m: "M0", running: "G3", queue: ["G1"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "a", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "b", kind: "unbuffered", slots: ['"from b" handed off'], senders: [], receivers: [] },
          { name: "timer heap", kind: "timer", slots: ["G2 wakes at 20ms"], senders: [], receivers: [] },
        ],
        chips: ["b case wins", "G1 removed from a queue too"],
      },
      {
        title: "main prints b value",
        time: 11,
        line: 25,
        concept: "only one case runs",
        explanation:
          "G1 returns from select via the b case and prints. The a case is not executed even if G2 eventually wakes. If multiple cases were ready at once, the runtime picks one at random.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "blocked", "ready", "running"] },
          { id: "G2", role: "a sender", state: "blocked", history: ["", "ready", "blocked", "blocked", "blocked", "blocked"] },
          { id: "G3", role: "b sender", state: "done", history: ["", "ready", "blocked", "blocked", "running", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "a", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "b", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "timer heap", kind: "timer", slots: ["G2 wakes at 20ms"], senders: [], receivers: [] },
        ],
        chips: ["one branch runs", "random when tied"],
      },
      {
        title: "main returns, G2 left behind",
        time: 12,
        line: 27,
        concept: "main exit kills siblings",
        explanation:
          "main returns and the process exits. G2 is still parked in the timer heap and never gets to send. Forgetting to coordinate sleeping goroutines is how leaks happen in real code.",
        goroutines: [
          { id: "G1", role: "main", state: "done", history: ["running", "running", "running", "blocked", "ready", "running", "done"] },
          { id: "G2", role: "a sender", state: "blocked", history: ["", "ready", "blocked", "blocked", "blocked", "blocked", "blocked"] },
          { id: "G3", role: "b sender", state: "done", history: ["", "ready", "blocked", "blocked", "running", "done", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "a", kind: "unbuffered", slots: [], senders: [], receivers: [] },
          { name: "b", kind: "unbuffered", slots: [], senders: [], receivers: [] },
        ],
        chips: ["use context for cancel", "process exits"],
      },
    ],
  },
  {
    id: "forloop",
    title: "for loop spawning goroutines (closure capture)",
    code: [
      "package main",
      "",
      "import (",
      '    "fmt"',
      '    "time"',
      ")",
      "",
      "func main() {",
      "    for i := 0; i <= 3; i++ {",
      "        go func(i int) {",
      "            fmt.Println(i)",
      "        }(i)",
      "    }",
      "",
      "    time.Sleep(time.Second * 2)",
      "}",
    ],
    steps: [
      {
        title: "for loop starts",
        time: 0,
        line: 9,
        concept: "loop drives spawns",
        explanation:
          "main enters the for loop. Each iteration will spawn one goroutine. The loop is using i <= 3 here to keep the diagram readable; the same shape holds for i <= 10 in your original code.",
        goroutines: [{ id: "G1", role: "main", state: "running", history: ["running"] }],
        scheduler: { p0: { m: "M0", running: "G1", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [{ name: "loop var i", kind: "stack value", slots: ["i = 0"], senders: [], receivers: [] }],
        chips: ["for init", "i := 0"],
      },
      {
        title: "iter 0 spawns G2 with i=0",
        time: 1,
        line: 12,
        concept: "pass i as parameter",
        explanation:
          "go func(i int){...}(i) passes the current loop value as a parameter, so G2 gets its own copy of 0. This is the safe pattern: if you closed over i without passing it, every goroutine would share the same variable (a classic pre-Go-1.22 bug).",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running"] },
          { id: "G2", role: "print 0", state: "ready", history: ["", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "loop var i", kind: "stack value", slots: ["i = 0"], senders: [], receivers: [] },
          { name: "G2 captured i", kind: "param copy", slots: ["i = 0"], senders: [], receivers: [] },
        ],
        chips: ["pass by value", "no shared i"],
      },
      {
        title: "iter 1 spawns G3 with i=1",
        time: 2,
        line: 12,
        concept: "queue grows",
        explanation:
          "Loop increments i to 1, spawns G3 with its own captured 1. Both G2 and G3 are now sitting ready to run. The scheduler will hand them to any free P.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running"] },
          { id: "G2", role: "print 0", state: "ready", history: ["", "ready", "ready"] },
          { id: "G3", role: "print 1", state: "ready", history: ["", "", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G2", "G3"] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "loop var i", kind: "stack value", slots: ["i = 1"], senders: [], receivers: [] },
          { name: "captured i values", kind: "param copies", slots: ["G2: 0", "G3: 1"], senders: [], receivers: [] },
        ],
        chips: ["i = 1", "two ready goroutines"],
      },
      {
        title: "iter 2 and 3 spawn G4 and G5",
        time: 3,
        line: 12,
        concept: "spawn is cheap",
        explanation:
          "Creating goroutines is cheap, roughly a few KB of stack and a runtime struct. The loop blasts through four spawns. P1 has been picking up work too, so G2 has already started running on M1.",
        goroutines: [
          { id: "G1", role: "main", state: "running", history: ["running", "running", "running", "running"] },
          { id: "G2", role: "print 0", state: "running", history: ["", "ready", "ready", "running"] },
          { id: "G3", role: "print 1", state: "ready", history: ["", "", "ready", "ready"] },
          { id: "G4", role: "print 2", state: "ready", history: ["", "", "", "ready"] },
          { id: "G5", role: "print 3", state: "ready", history: ["", "", "", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G1", queue: ["G3", "G4", "G5"] }, p1: { m: "M1", running: "G2", queue: [] }, global: [] },
        channels: [
          { name: "loop var i", kind: "stack value", slots: ["i = 3"], senders: [], receivers: [] },
          { name: "captured i values", kind: "param copies", slots: ["G2: 0", "G3: 1", "G4: 2", "G5: 3"], senders: [], receivers: [] },
        ],
        chips: ["lightweight spawn", "queue fills"],
      },
      {
        title: "loop ends, main calls Sleep",
        time: 4,
        line: 15,
        concept: "Sleep keeps main alive",
        explanation:
          "i becomes 4, the condition fails, the loop ends. main hits time.Sleep and parks on the timer heap. Without this sleep, main would return immediately and kill the still-pending goroutines.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "running", "blocked"] },
          { id: "G2", role: "print 0", state: "running", history: ["", "ready", "ready", "running", "running"] },
          { id: "G3", role: "print 1", state: "running", history: ["", "", "ready", "ready", "running"] },
          { id: "G4", role: "print 2", state: "ready", history: ["", "", "", "ready", "ready"] },
          { id: "G5", role: "print 3", state: "ready", history: ["", "", "", "ready", "ready"] },
        ],
        scheduler: { p0: { m: "M0", running: "G3", queue: ["G4", "G5"] }, p1: { m: "M1", running: "G2", queue: [] }, global: [] },
        channels: [
          { name: "timer heap", kind: "timer", slots: ["G1 wakes at 2s"], senders: [], receivers: [] },
          { name: "captured i values", kind: "param copies", slots: ["G2: 0", "G3: 1", "G4: 2", "G5: 3"], senders: [], receivers: [] },
        ],
        chips: ["main parked", "workers free to run"],
      },
      {
        title: "workers print in unpredictable order",
        time: 5,
        line: 11,
        concept: "no ordering guarantee",
        explanation:
          "G2 and G3 finish their prints. G4 takes M0, G5 takes M1. The print order is NOT 0,1,2,3 in general. With multiple Ps, the order depends on which goroutine the scheduler picked up first.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "running", "blocked", "blocked"] },
          { id: "G2", role: "print 0", state: "done", history: ["", "ready", "ready", "running", "running", "done"] },
          { id: "G3", role: "print 1", state: "done", history: ["", "", "ready", "ready", "running", "done"] },
          { id: "G4", role: "print 2", state: "running", history: ["", "", "", "ready", "ready", "running"] },
          { id: "G5", role: "print 3", state: "running", history: ["", "", "", "ready", "ready", "running"] },
        ],
        scheduler: { p0: { m: "M0", running: "G4", queue: [] }, p1: { m: "M1", running: "G5", queue: [] }, global: [] },
        channels: [
          { name: "timer heap", kind: "timer", slots: ["G1 wakes at 2s"], senders: [], receivers: [] },
          { name: "stdout (recent prints)", kind: "output", slots: ["0", "1"], senders: [], receivers: [] },
        ],
        chips: ["order is not 0,1,2,3", "GOMAXPROCS matters"],
      },
      {
        title: "all workers done, main still sleeping",
        time: 6,
        line: 15,
        concept: "we are overshooting Sleep",
        explanation:
          "All four workers exit but main is still parked until its 2-second timer fires. time.Sleep is a blunt tool; in real code use a WaitGroup or a done channel so main wakes the moment work finishes.",
        goroutines: [
          { id: "G1", role: "main", state: "blocked", history: ["running", "running", "running", "running", "blocked", "blocked", "blocked"] },
          { id: "G2", role: "print 0", state: "done", history: ["", "ready", "ready", "running", "running", "done", "done"] },
          { id: "G3", role: "print 1", state: "done", history: ["", "", "ready", "ready", "running", "done", "done"] },
          { id: "G4", role: "print 2", state: "done", history: ["", "", "", "ready", "ready", "running", "done"] },
          { id: "G5", role: "print 3", state: "done", history: ["", "", "", "ready", "ready", "running", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [
          { name: "timer heap", kind: "timer", slots: ["G1 wakes at 2s"], senders: [], receivers: [] },
          { name: "stdout (full)", kind: "output", slots: ["0", "1", "2", "3"], senders: [], receivers: [] },
        ],
        chips: ["wasted wait time", "prefer sync.WaitGroup"],
      },
      {
        title: "timer fires, main exits",
        time: 7,
        line: 16,
        concept: "process ends",
        explanation:
          "After 2 seconds the timer fires, G1 becomes runnable, executes the return from main, and the process exits. Compare this scenario to the WaitGroup one to see why Sleep is a fragile substitute for real coordination.",
        goroutines: [
          { id: "G1", role: "main", state: "done", history: ["running", "running", "running", "running", "blocked", "blocked", "blocked", "done"] },
          { id: "G2", role: "print 0", state: "done", history: ["", "ready", "ready", "running", "running", "done", "done", "done"] },
          { id: "G3", role: "print 1", state: "done", history: ["", "", "ready", "ready", "running", "done", "done", "done"] },
          { id: "G4", role: "print 2", state: "done", history: ["", "", "", "ready", "ready", "running", "done", "done"] },
          { id: "G5", role: "print 3", state: "done", history: ["", "", "", "ready", "ready", "running", "done", "done"] },
        ],
        scheduler: { p0: { m: "M0", running: "", queue: [] }, p1: { m: "M1", running: "", queue: [] }, global: [] },
        channels: [],
        chips: ["main returns", "all goroutines accounted for"],
      },
    ],
  },
];

const stateClass = {
  new: "state-new",
  ready: "state-ready",
  running: "state-running",
  blocked: "state-blocked",
  done: "state-done",
};

const scenarioSelect = document.querySelector("#scenarioSelect");
const resetBtn = document.querySelector("#resetBtn");
const prevBtn = document.querySelector("#prevBtn");
const stepBtn = document.querySelector("#stepBtn");
const playBtn = document.querySelector("#playBtn");
const speedRange = document.querySelector("#speedRange");
const codeBlock = document.querySelector("#codeBlock");
const goroutineLanes = document.querySelector("#goroutineLanes");
const schedulerBoard = document.querySelector("#schedulerBoard");
const channelBoard = document.querySelector("#channelBoard");
const stepCounter = document.querySelector("#stepCounter");
const runtimeClock = document.querySelector("#runtimeClock");
const activeScenario = document.querySelector("#activeScenario");
const stepTitle = document.querySelector("#stepTitle");
const currentConcept = document.querySelector("#currentConcept");
const stepExplanation = document.querySelector("#stepExplanation");
const conceptChips = document.querySelector("#conceptChips");

let activeScenarioIndex = 0;
let activeStepIndex = 0;
let playTimer = null;

function html(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

function getScenario() {
  return scenarios[activeScenarioIndex];
}

function getStep() {
  return getScenario().steps[activeStepIndex];
}

function setupScenarioOptions() {
  scenarioSelect.innerHTML = scenarios
    .map((scenario, index) => `<option value="${index}">${html(scenario.title)}</option>`)
    .join("");
}

function renderCode() {
  const scenario = getScenario();
  const activeLine = getStep().line;
  codeBlock.innerHTML = scenario.code
    .map((line, index) => {
      const lineNumber = index + 1;
      const active = lineNumber === activeLine ? " active" : "";
      return `<span class="code-line${active}"><span class="num">${lineNumber}</span><span>${html(line) || " "}</span></span>`;
    })
    .join("");
}

function renderGoroutines() {
  const step = getStep();
  const maxHistory = Math.max(...step.goroutines.map((g) => g.history.length), 1);
  const tickCount = Math.max(maxHistory, 6);

  goroutineLanes.innerHTML = step.goroutines
    .map((g) => {
      const ticks = Array.from({ length: tickCount }, (_, index) => {
        const status = g.history[index] || "";
        const filled = status ? " filled" : "";
        const title = status ? `${g.id} ${status}` : `${g.id} not created`;
        return `<div class="tick ${html(status)}${filled}" title="${html(title)}"></div>`;
      }).join("");

      return `
        <div class="g-row">
          <div class="g-label">
            <div class="g-name">${html(g.id)}</div>
            <div class="g-role">${html(g.role)}</div>
            <span class="state-pill ${stateClass[g.state] || "state-ready"}">${html(g.state)}</span>
          </div>
          <div class="timeline" style="grid-template-columns: repeat(${tickCount}, minmax(28px, 1fr));">${ticks}</div>
        </div>
      `;
    })
    .join("");
}

function renderScheduler() {
  const scheduler = getStep().scheduler;
  const renderMachine = (name, machine) => `
    <div class="machine">
      <div class="machine-head">
        <span>${html(name)}</span>
        <span>${html(machine.m || "no M")}</span>
      </div>
      <div class="processor">
        ${machine.running ? `<span class="g-token running">${html(machine.running)}</span>` : "idle"}
      </div>
      <div class="run-queue" aria-label="${html(name)} local run queue">
        ${machine.queue.length ? machine.queue.map((g) => `<span class="queue-token">${html(g)}</span>`).join("") : "<span class=\"empty-text\">local queue empty</span>"}
      </div>
    </div>
  `;

  schedulerBoard.innerHTML = `
    <div class="machine-grid">
      ${renderMachine("P0", scheduler.p0)}
      ${renderMachine("P1", scheduler.p1)}
    </div>
    <div class="global-queue">
      <strong>Global run queue</strong>
      <div class="run-queue">
        ${scheduler.global.length ? scheduler.global.map((g) => `<span class="queue-token global">${html(g)}</span>`).join("") : "<span class=\"empty-text\">empty</span>"}
      </div>
    </div>
    <div class="note-box">
      <strong>Scheduler idea</strong>
      G is the goroutine. M is an OS thread. P is the runtime processor token that lets an M execute Go code.
    </div>
  `;
}

function renderChannels() {
  const channels = getStep().channels;
  if (!channels.length) {
    channelBoard.innerHTML = `
      <div class="note-box">
        <strong>No channel wait right now</strong>
        The current step is about goroutine scheduling and timers.
      </div>
    `;
    return;
  }

  channelBoard.innerHTML = channels
    .map((channel) => {
      const slots =
        channel.slots.length > 0
          ? channel.slots
              .map((slot) => {
                const filled = slot !== "empty" ? " filled" : "";
                return `<div class="slot${filled}">${html(slot)}</div>`;
              })
              .join("")
          : `<div class="slot">no buffer</div>`;

      const senders = channel.senders.length ? channel.senders.join(", ") : "none";
      const receivers = channel.receivers.length ? channel.receivers.join(", ") : "none";

      return `
        <div class="channel-object">
          <div class="channel-name">
            ${html(channel.name)}
            <span>${html(channel.kind)}</span>
          </div>
          <div class="channel-slots">${slots}</div>
          <div class="waiters">
            <span>senders: ${html(senders)}</span>
            <span>receivers: ${html(receivers)}</span>
          </div>
        </div>
      `;
    })
    .join("");
}

function renderExplanation() {
  const scenario = getScenario();
  const step = getStep();
  stepCounter.textContent = `Step ${activeStepIndex + 1} of ${scenario.steps.length}`;
  runtimeClock.textContent = `t = ${step.time}ms`;
  activeScenario.textContent = scenario.title;
  stepTitle.textContent = step.title;
  currentConcept.textContent = step.concept;
  stepExplanation.textContent = step.explanation;
  conceptChips.innerHTML = step.chips.map((chip) => `<span>${html(chip)}</span>`).join("");
}

function render() {
  renderCode();
  renderGoroutines();
  renderScheduler();
  renderChannels();
  renderExplanation();
  prevBtn.disabled = activeStepIndex === 0;
  stepBtn.disabled = activeStepIndex === getScenario().steps.length - 1;
}

function stopPlayback() {
  if (playTimer) {
    window.clearInterval(playTimer);
    playTimer = null;
  }
  playBtn.textContent = "";
  playBtn.innerHTML = '<span aria-hidden="true">P</span>Play';
}

function nextStep() {
  const scenario = getScenario();
  if (activeStepIndex < scenario.steps.length - 1) {
    activeStepIndex += 1;
    render();
    return true;
  }
  stopPlayback();
  return false;
}

function previousStep() {
  stopPlayback();
  activeStepIndex = Math.max(0, activeStepIndex - 1);
  render();
}

function resetScenario() {
  stopPlayback();
  activeStepIndex = 0;
  render();
}

function startPlayback() {
  if (playTimer) {
    stopPlayback();
    return;
  }

  if (activeStepIndex === getScenario().steps.length - 1) {
    activeStepIndex = 0;
    render();
  }

  playBtn.innerHTML = '<span aria-hidden="true">S</span>Stop';
  playTimer = window.setInterval(() => {
    nextStep();
  }, currentInterval());
}

function currentInterval() {
  const min = Number(speedRange.min);
  const max = Number(speedRange.max);
  const val = Number(speedRange.value);
  return min + max - val;
}

scenarioSelect.addEventListener("change", (event) => {
  stopPlayback();
  activeScenarioIndex = Number(event.target.value);
  activeStepIndex = 0;
  render();
});

resetBtn.addEventListener("click", resetScenario);
prevBtn.addEventListener("click", previousStep);
stepBtn.addEventListener("click", () => {
  stopPlayback();
  nextStep();
});
playBtn.addEventListener("click", startPlayback);
speedRange.addEventListener("input", () => {
  if (playTimer) {
    stopPlayback();
    startPlayback();
  }
});

setupScenarioOptions();
render();
