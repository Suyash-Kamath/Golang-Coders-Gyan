// Package main shows that context forms a TREE: cancelling a node cancels every
// descendant, but never its ancestors or siblings. Root -> A -> grandchild(C), Root -> B.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// worker simulates work that keeps going until its context says stop.
func worker(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] stopping: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("[%s] working...\n", name)
			time.Sleep(150 * time.Millisecond)
		}
	}
}

func main() {
	fmt.Println("=== Scenario 1: cancelling the ROOT stops the entire tree (A, its grandchild, and B) ===")
	scenario(func(cancelRoot, cancelChildA, cancelGrandchild context.CancelFunc) {
		time.Sleep(300 * time.Millisecond)
		fmt.Println(">>> cancelling ROOT")
		cancelRoot()
	})

	fmt.Println("\n=== Scenario 2: cancelling CHILD A stops A and its grandchild, but NOT sibling B or the root ===")
	scenario(func(cancelRoot, cancelChildA, cancelGrandchild context.CancelFunc) {
		time.Sleep(300 * time.Millisecond)
		fmt.Println(">>> cancelling CHILD A")
		cancelChildA()
		time.Sleep(300 * time.Millisecond)
		fmt.Println(">>> cancelling ROOT to let the demo finish (B was still running until now)")
		cancelRoot()
	})
}

// scenario builds root -> A -> grandchild, and root -> B, starts a worker under
// each leaf, then lets the caller decide what to cancel and when.
func scenario(drive func(cancelRoot, cancelChildA, cancelGrandchild context.CancelFunc)) {
	var wg sync.WaitGroup

	rootCtx, cancelRoot := context.WithCancel(context.Background())  // "root" (parent)
	childACtx, cancelChildA := context.WithCancel(rootCtx)           // "A" — child of root
	grandchildCtx, cancelGrandchild := context.WithCancel(childACtx) // "C" — child of A
	childBCtx, cancelChildB := context.WithCancel(rootCtx)           // "B" — sibling of A
	defer cancelRoot()
	defer cancelChildA()
	defer cancelGrandchild()
	defer cancelChildB()

	wg.Add(3)
	go worker(childACtx, "A", &wg)
	go worker(grandchildCtx, "A.grandchild (C)", &wg)
	go worker(childBCtx, "B", &wg)

	drive(cancelRoot, cancelChildA, cancelGrandchild)

	wg.Wait()
}
