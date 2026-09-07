// Package main draws the line between context.CancelFunc and sync.WaitGroup:
// cancel() only ever SIGNALS "please stop" — it returns immediately and has no
// idea whether anything actually stopped. wg.Wait() is what CONFIRMS it did.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	const n = 3
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("worker %d: received cancel, finishing current unit of work...\n", id)
					time.Sleep(80 * time.Millisecond) // simulate a bit of cleanup work
					fmt.Printf("worker %d: fully stopped\n", id)
					return
				case <-time.After(500 * time.Millisecond):
					fmt.Printf("worker %d: heartbeat\n", id)
				}
			}
		}(i)
	}

	time.Sleep(300 * time.Millisecond)

	fmt.Println(">>> cancel() called — this returns IMMEDIATELY")
	cancel()
	fmt.Println(">>> cancel() returned, but the workers might still be cleaning up!")

	fmt.Println(">>> now waiting on wg.Wait() for actual confirmation...")
	wg.Wait()
	fmt.Println(">>> wg.Wait() returned — NOW we know every worker truly stopped")
}
