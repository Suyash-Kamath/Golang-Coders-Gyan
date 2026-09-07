// Package main catalogues classic goroutine-leak traps around context and channels:
// busy-looping instead of blocking, blocking sends with no escape hatch, and the
// "first result wins" fan-out pattern.
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Leak trap #1: a busy 'default' branch spins the CPU instead of blocking ===")
	busyLoopBad()
	busyLoopGood()

	fmt.Println("\n=== Leak trap #2: blocking on an unbuffered channel send with no escape hatch ===")
	blockedSendLeak()

	fmt.Println("\n=== Fix A: guard the send with select on ctx.Done() ===")
	blockedSendFixedWithSelect()

	fmt.Println("\n=== Fix B: size the channel so every sender can always complete its send ===")
	firstResultWinsBuffered()
}

// busyLoopBad demonstrates a real anti-pattern: a `default:` branch inside a tight
// loop with no blocking point at all burns a full CPU core just polling ctx.Done().
func busyLoopBad() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	iterations := 0
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			iterations++ // no sleep, no blocking — this spins as fast as the CPU allows
		}
	}
	fmt.Printf("busyLoopBad:  %d spins in 20ms (wasted CPU)\n", iterations)
}

// busyLoopGood does the same job but blocks on a real timer between checks.
func busyLoopGood() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	ticks := 0
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			ticks++ // do real work here, then go back to sleep
		}
	}
	fmt.Printf("busyLoopGood: %d ticks in 20ms (CPU idle between checks)\n", ticks)
}

// blockedSendLeak starts a goroutine that tries to send a single result on an
// UNBUFFERED channel. Nobody ever receives it, and the goroutine has no ctx
// check around the send, so it blocks forever — a permanent goroutine leak.
func blockedSendLeak() {
	results := make(chan string) // unbuffered
	go func() {
		time.Sleep(10 * time.Millisecond)
		results <- "done" // nobody is listening -> blocks forever, goroutine leaks
		fmt.Println("this line never runs")
	}()
	fmt.Println("blockedSendLeak: goroutine launched and abandoned (leaks for the life of the program)")
	time.Sleep(30 * time.Millisecond)
}

// blockedSendFixedWithSelect fixes the leak: even the SEND is wrapped in a select
// against ctx.Done(), so if nobody ever reads, the goroutine still exits.
func blockedSendFixedWithSelect() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	results := make(chan string) // still unbuffered
	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(10 * time.Millisecond)
		select {
		case results <- "done":
			fmt.Println("sent result")
		case <-ctx.Done():
			fmt.Println("blockedSendFixedWithSelect: nobody was listening, exiting cleanly:", ctx.Err())
		}
	}()
	<-done
}

// firstResultWinsBuffered fans out N requests but only ever wants the FIRST reply.
// A buffered channel sized to N guarantees every goroutine can complete its send
// and exit immediately, even though N-1 of them will never be read.
func firstResultWinsBuffered() {
	urls := []string{"a", "bb", "ccc", "dddd", "eeeee"}
	results := make(chan string, len(urls)) // buffered = len(urls): every send succeeds immediately

	for _, u := range urls {
		go func() {
			time.Sleep(time.Duration(len(u)*10) * time.Millisecond)
			results <- fmt.Sprintf("result from %q", u) // never blocks, so never leaks
		}()
	}

	first := <-results // take only the first result
	fmt.Println("first result:", first)
	// the other 4 goroutines still run to completion and send into the buffer,
	// then exit normally — no leak, even though we never read their values.
	time.Sleep(80 * time.Millisecond)
}
