# Context in Go — worked examples

Each numbered folder is a standalone, runnable `package main` covering one part of the
`context` package, in the order you'd normally learn it. Run any of them with:

```sh
cd codes
go run ./01_basics
go run ./02_with_cancel
# ...etc
```

| # | Folder | Covers |
|---|--------|--------|
| 01 | `01_basics` | `context.Background()` vs `TODO()`, the four `Context` methods (`Deadline`, `Done`, `Err`, `Value`) |
| 02 | `02_with_cancel` | `WithCancel`, the context tree, parent/child/grandchild cancellation propagation |
| 03 | `03_with_timeout_deadline` | `WithTimeout` vs `WithDeadline`, the `select{time.After vs ctx.Done()}` race, deadline budgeting |
| 04 | `04_with_value` | Type-safe context keys, typed accessor functions, upward value lookup, what belongs (and doesn't) in `Value` |
| 05 | `05_goroutines_and_leaks` | Busy-loop CPU spin trap, blocked-channel-send leaks, the "first result wins" fan-out pattern |
| 06 | `06_worker_pool` | N workers sharing one context and one jobs channel, single `cancel()` shuts the whole pool down |
| 07 | `07_context_and_waitgroup` | Why `cancel()` (signal) and `wg.Wait()` (confirmation) are different things |
| 08 | `08_http_server` | Handler → Service → Repository → DB, `r.Context()`, `QueryRowContext`, `errors.Is` error handling |
| 09 | `09_http_client_and_fanout` | `http.NewRequestWithContext`, concurrent fan-out with shared cancel, an errgroup-style helper |
| 10 | `10_cancelcause_and_withoutcancel` | `WithCancelCause`/`context.Cause`, `WithTimeoutCause`, `WithoutCancel` for detached best-effort work |
| 11 | `11_graceful_shutdown` | `signal.NotifyContext` + `http.Server.Shutdown(ctx)` |
| 12 | `12_middleware` | Request-ID and auth middleware using `r.WithContext` |
| 13 | `13_anti_patterns` | Annotated bad-vs-good: context in structs, context as a value bag, nil context, swallowing propagation |
| 14 | `14_retries_ratelimiting` | Context-aware retry/backoff loop, a `Wait(ctx)`-style rate limiter |
| 15 | `15_internals_miniature_context` | From-scratch mini `Context` implementation showing how cancellation propagation actually works internally |

## The one-paragraph mental model

`context.Context` answers four questions for whoever is doing work on your behalf:
*when should I stop?* (`Done()`), *why did I stop?* (`Err()`), *do I have a hard
deadline?* (`Deadline()`), and *is there any request-scoped metadata I should know
about?* (`Value()`). Contexts form a **tree**: cancelling a node cancels every
descendant, but never its parent or siblings. `cancel()` only *signals* — it never
confirms anything actually stopped; that's what a `sync.WaitGroup` is for.
Cancellation is entirely **cooperative**: nothing happens unless the function doing
the work explicitly checks `ctx.Done()` (usually inside a `select`).
