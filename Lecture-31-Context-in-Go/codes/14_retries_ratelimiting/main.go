// Package main shows context-aware retries (check ctx before every attempt, and
// back off using select instead of a plain time.Sleep so cancellation interrupts
// the wait too) and a minimal context-aware rate limiter modeled after
// golang.org/x/time/rate.Limiter.Wait(ctx).
//
// Note: not everything is context-aware — e.g. sync.Mutex.Lock() cannot be
// interrupted by a context. Cancellation only works when the function you're
// calling actually checks ctx.Done() itself.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errFlaky = errors.New("transient failure")

func flakyCall(attempt int) error {
	if attempt < 3 {
		return errFlaky
	}
	return nil
}

func callWithRetry(ctx context.Context, maxAttempts int) error {
	backoff := 20 * time.Millisecond
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err // don't even try again if the caller already gave up
		}

		err := flakyCall(attempt)
		if err == nil {
			fmt.Printf("succeeded on attempt %d\n", attempt)
			return nil
		}
		fmt.Printf("attempt %d failed: %v\n", attempt, err)

		if attempt == maxAttempts {
			return fmt.Errorf("giving up after %d attempts: %w", maxAttempts, err)
		}

		select {
		case <-time.After(backoff): // wait for the backoff...
		case <-ctx.Done(): // ...unless the context dies first
			return ctx.Err()
		}
		backoff *= 2
	}
	return nil
}

// tokenBucket is a minimal, context-aware rate limiter. Wait blocks until either
// a token is available or the context is cancelled.
type tokenBucket struct {
	tokens   chan struct{}
	interval time.Duration
}

func newTokenBucket(capacity int, refill time.Duration) *tokenBucket {
	tb := &tokenBucket{tokens: make(chan struct{}, capacity), interval: refill}
	for i := 0; i < capacity; i++ {
		tb.tokens <- struct{}{}
	}
	go tb.refillLoop()
	return tb
}

func (tb *tokenBucket) refillLoop() {
	ticker := time.NewTicker(tb.interval)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case tb.tokens <- struct{}{}:
		default: // bucket already full
		}
	}
}

func (tb *tokenBucket) Wait(ctx context.Context) error {
	select {
	case <-tb.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	fmt.Println("=== retry with context-aware backoff ===")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := callWithRetry(ctx, 5); err != nil {
		fmt.Println("final error:", err)
	}

	fmt.Println("\n=== rate limiting: only 2 tokens available up front, 5 requests want one ===")
	limiter := newTokenBucket(2, 100*time.Millisecond)
	rlCtx, rlCancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer rlCancel()

	for i := 1; i <= 5; i++ {
		start := time.Now()
		if err := limiter.Wait(rlCtx); err != nil {
			fmt.Printf("request %d: rate limiter wait aborted: %v\n", i, err)
			continue
		}
		fmt.Printf("request %d: got a token after %v\n", i, time.Since(start).Round(time.Millisecond))
	}
}
