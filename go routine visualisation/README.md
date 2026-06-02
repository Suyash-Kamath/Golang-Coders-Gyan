# Go Routine Visualisation

A Loupe-style visualizer for learning how goroutines move through the Go runtime.

## Running it

Open `index.html` directly in a browser, or serve the folder:

```
python3 -m http.server 8765
# then open http://localhost:8765
```

## Controls

- `Scenario` — switch between the six built-in scenarios.
- `Step` / `Prev` — move one runtime event at a time.
- `Play` — auto-run the selected scenario.
- `Speed` — slow down or speed up playback.
- `Reset` — jump back to step 1.

## Scenarios

1. **Basic goroutine scheduling** — `go`, the run queue, sleep, wakeup, and main-goroutine exit.
2. **Unbuffered channel handoff** — receiver park, sender park, direct hand-off.
3. **Buffered channel and backpressure** — capacity, full-buffer wait, FIFO drain.
4. **sync.WaitGroup coordination** — `Add` / `Done` / `Wait`, counter park-and-release.
5. **select with multiple channels** — parking on every case, first-ready wins, random tie-break.
6. **for loop spawning goroutines** — closure capture vs. passing `i` as a parameter, ordering caveats, why `time.Sleep` is a fragile substitute for coordination.

## What you see in each panel

- **Goroutines** — per-`G` lifecycle timeline. States: `ready`, `running`, `blocked`, `done`.
- **Scheduler** — the `G / M / P` model: `M` is an OS thread, `P` is a processor token, `G` is a goroutine. An `M` must hold a `P` to execute Go code.
- **Channels, timers, waits** — channel buffer slots and sender/receiver wait queues; timer-heap entries; `WaitGroup` counter.
- **What is happening** — plain-English explanation of the current step.

## Caveats

This is a teaching model. The real Go runtime also does work stealing, syscall handling, preemption, network polling, timers, GC cooperation, and more. Step-by-step ordering here is illustrative, not literal.
