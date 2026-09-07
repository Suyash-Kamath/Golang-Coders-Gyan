// Package main compares WithTimeout vs WithDeadline, the classic
// select{ time.After vs ctx.Done() } race pattern, and deadline budgeting:
// a child context can never end up with a MORE generous deadline than its parent.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== WithTimeout: cancel automatically after a duration ===")
	timeoutDemo()

	fmt.Println("\n=== WithDeadline: cancel automatically at a specific point in time ===")
	deadlineDemo()

	fmt.Println("\n=== Deadline budgeting: a child can only ever be as generous as its parent ===")
	budgetDemo()
}

// slowOperation pretends to do work that takes `work` time, but honors cancellation.
func slowOperation(ctx context.Context, work time.Duration) (string, error) {
	select {
	case <-time.After(work):
		return "operation finished", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func timeoutDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel() // always defer cancel, even though the timer will fire on its own eventually

	result, err := slowOperation(ctx, 2*time.Second) // work takes way longer than the timeout
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("slowOperation timed out as expected:", err)
		} else {
			fmt.Println("unexpected error:", err)
		}
		return
	}
	fmt.Println("result:", result)
}

func deadlineDemo() {
	deadline := time.Now().Add(200 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// WithTimeout(ctx, d) is literally implemented as WithDeadline(ctx, time.Now().Add(d)).
	result, err := slowOperation(ctx, 2*time.Second)
	if err != nil {
		fmt.Println("slowOperation hit the deadline:", err)
		return
	}
	fmt.Println("result:", result)
}

func budgetDemo() {
	fmt.Println("-- parent has 100ms, child asks for 500ms -> child is capped at ~100ms --")
	parent, cancelParent := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelParent()
	child, cancelChild := context.WithTimeout(parent, 500*time.Millisecond)
	defer cancelChild()
	waitAndReport(child)

	fmt.Println("-- parent has 500ms, child asks for only 100ms -> child really does get just 100ms --")
	parent2, cancelParent2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelParent2()
	child2, cancelChild2 := context.WithTimeout(parent2, 100*time.Millisecond)
	defer cancelChild2()
	waitAndReport(child2)
}

func waitAndReport(ctx context.Context) {
	start := time.Now()
	<-ctx.Done()
	fmt.Printf("   stopped after %v: %v\n", time.Since(start).Round(time.Millisecond), ctx.Err())
}
