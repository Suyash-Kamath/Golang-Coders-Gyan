// Package main shows a worker pool where N workers share ONE context: a single
// cancel() call shuts every worker down, no matter how many there are.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int)
	var wg sync.WaitGroup

	const numWorkers = 5
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(ctx, i, jobs, &wg)
	}

	// feed jobs in the background
	go func() {
		for j := 1; j <= 12; j++ {
			select {
			case jobs <- j:
			case <-ctx.Done():
				return
			}
			time.Sleep(30 * time.Millisecond)
		}
	}()

	time.Sleep(200 * time.Millisecond)
	fmt.Println(">>> shutting down the whole pool with a single cancel()")
	cancel() // one signal stops every worker, no matter how many there are

	wg.Wait()
	fmt.Println("all workers stopped")
}

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: shutting down (%v)\n", id, ctx.Err())
			return
		case j, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("worker %d: processing job %d\n", id, j)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
