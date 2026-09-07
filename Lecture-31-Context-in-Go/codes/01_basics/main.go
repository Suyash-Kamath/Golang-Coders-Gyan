// Package main demonstrates the absolute basics of context.Context:
// Background() vs TODO(), and the four methods every Context exposes
// (Deadline, Done, Err, Value).
package main

import (
	"context"
	"fmt"
	"time"
)

type demoKey string

func main() {
	fmt.Println("=== 1. context.Background() and context.TODO() ===")
	bg := context.Background()
	todo := context.TODO()
	// Both are empty root contexts: never cancelled, no deadline, no values.
	// Background() = "I am intentionally starting a new context tree" (main, tests, top-level).
	// TODO()       = "I know this function should take a context, but I haven't wired it up yet."
	fmt.Printf("Background: %v\n", bg)
	fmt.Printf("TODO:       %v\n", todo)

	fmt.Println("\n=== 2. The four things every Context can answer ===")
	inspect(bg)

	fmt.Println("\n=== 3. Done() — a channel that gets CLOSED, not sent on ===")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done() // blocks until the channel is closed
		fmt.Println("goroutine: context cancelled, cleaning up")
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond) // let the goroutine print before main exits

	fmt.Println("\n=== 4. Err() tells you WHY it stopped ===")
	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2()
	fmt.Println("Err after manual cancel:   ", ctx2.Err()) // context.Canceled

	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel3()
	time.Sleep(20 * time.Millisecond)
	fmt.Println("Err after timeout elapses: ", ctx3.Err()) // context.DeadlineExceeded

	fmt.Println("\n=== 5. Deadline() ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), time.Second)
	defer cancel4()
	deadline, ok := ctx4.Deadline()
	fmt.Printf("Has deadline: %v, deadline: %s\n", ok, deadline.Format(time.RFC3339))
	_, ok = context.Background().Deadline()
	fmt.Printf("Background has deadline: %v\n", ok)

	fmt.Println("\n=== 6. Value() — do NOT use raw strings as keys (see 04_with_value) ===")
	ctxVal := context.WithValue(context.Background(), demoKey("greeting"), "hello")
	fmt.Println("Value():          ", ctxVal.Value(demoKey("greeting")))
	fmt.Println("Missing key nil?  ", ctxVal.Value(demoKey("missing")) == nil)
}

// inspect prints what the four Context methods report for a fresh context.
func inspect(ctx context.Context) {
	deadline, hasDeadline := ctx.Deadline()
	fmt.Println("Deadline():", deadline, "/ ok =", hasDeadline)
	fmt.Println("Done():    ", ctx.Done()) // nil channel for Background/TODO — never fires
	fmt.Println("Err():     ", ctx.Err())
	fmt.Println("Value():   ", ctx.Value(demoKey("anything")))
}
