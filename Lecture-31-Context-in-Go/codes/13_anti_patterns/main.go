// Package main is a reference of context anti-patterns. Each BAD shape is shown
// as a comment (some compile fine and are purely design mistakes — they just
// bite you later), paired with the GOOD, working alternative right below it.
package main

import (
	"context"
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// 1) BAD: passing ctx into pure, non-blocking, non-IO functions "just in case".
//
//	func Add(ctx context.Context, a, b int) int { return a + b }
//
// GOOD: ctx is for cancellation/deadlines/request-scoped values — if a function
// never does I/O, never blocks, and never calls anything that does, skip it.
func Add(a, b int) int { return a + b }

// ---------------------------------------------------------------------------
// 2) BAD: using context.WithValue to avoid writing real parameters.
//
//	ctx = context.WithValue(ctx, "amount", 100)
//	ctx = context.WithValue(ctx, "currency", "USD")
//	func Charge(ctx context.Context) error { amount := ctx.Value("amount"); ... }
//
// GOOD: business inputs are normal, typed, explicit arguments.
func Charge(ctx context.Context, amount int, currency string) error {
	_ = ctx
	fmt.Printf("charging %d %s\n", amount, currency)
	return nil
}

// ---------------------------------------------------------------------------
// 3) BAD: storing a context inside a struct so you "don't have to pass it around".
//
//	type Service struct {
//	    ctx context.Context // smell: whose request is this? how long does it live?
//	    db  *sql.DB
//	}
//
// GOOD: context is a parameter of an OPERATION, not a property of a long-lived object.
type Service struct {
	// db *sql.DB  (real dependencies live here — never a context)
}

func (s Service) DoWork(ctx context.Context) error {
	select {
	case <-time.After(10 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ---------------------------------------------------------------------------
// 4) BAD: forgetting to call cancel() — the context (and its timer, if any) leaks
//    until its parent is garbage collected or the process exits.
//
//	ctx, _ := context.WithTimeout(context.Background(), time.Minute) // cancel discarded!
//
// GOOD: always defer cancel(), even if you expect the context to expire on its own.
func goodTimeoutUsage() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_ = ctx
}

// ---------------------------------------------------------------------------
// 5) BAD: reaching for context.Background() deep inside a call chain instead of
//    propagating the ctx you were given — this silently breaks cancellation.
//
//	func (r Repo) Get(ctx context.Context, id string) (Row, error) {
//	    return r.db.QueryRowContext(context.Background(), ...) // ctx thrown away!
//	}
//
// GOOD: pass the SAME ctx (or a derived child of it) all the way down.
type Repo struct{}

func (r Repo) Get(ctx context.Context, id string) (string, error) {
	select {
	case <-time.After(5 * time.Millisecond):
		return "row-" + id, nil
	case <-ctx.Done():
		return "", ctx.Err() // now a cancelled caller actually cancels this call too
	}
}

// ---------------------------------------------------------------------------
// 6) BAD: cancelling a context that nothing is listening to, and assuming that
//    alone stopped the work. cancel() only signals — see 07_context_and_waitgroup
//    for why you still need a WaitGroup (or similar) to confirm completion.

// ---------------------------------------------------------------------------
// 7) BAD: passing nil where a context is expected.
//
//	doSomething(nil, "x") // panics the moment anything calls ctx.Done()/Err()/Value()
//
// GOOD: use context.TODO() (not yet wired up) or context.Background() (a real root).
func doSomething(ctx context.Context, x string) {
	if ctx == nil {
		ctx = context.TODO()
	}
	fmt.Println("doing", x, "- err:", ctx.Err())
}

func main() {
	fmt.Println("Add(2, 3) =", Add(2, 3))
	_ = Charge(context.Background(), 100, "USD")

	svc := Service{}
	fmt.Println("DoWork:", svc.DoWork(context.Background()))

	goodTimeoutUsage()
	fmt.Println("goodTimeoutUsage: cancel() was deferred, no leak")

	repo := Repo{}
	row, _ := repo.Get(context.Background(), "42")
	fmt.Println("repo.Get:", row)

	doSomething(context.TODO(), "y")
}
