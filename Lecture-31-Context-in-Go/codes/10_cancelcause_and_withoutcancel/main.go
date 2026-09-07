// Package main covers three modern additions to context: WithCancelCause /
// context.Cause (attach a real reason to a cancellation), WithTimeoutCause, and
// WithoutCancel (deliberately detach from a parent's cancellation/deadline while
// keeping its values).
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== WithCancelCause / Cause(ctx): attach a REASON to a cancellation ===")
	cancelCauseDemo()

	fmt.Println("\n=== WithTimeoutCause: same idea, but for automatic timeouts ===")
	timeoutCauseDemo()

	fmt.Println("\n=== WithoutCancel: deliberately detach from the parent's cancellation ===")
	withoutCancelDemo()
}

var errRateLimited = errors.New("rate limited by upstream")

func cancelCauseDemo() {
	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel(errRateLimited) // plain cancel(nil) would just report context.Canceled
	}()

	<-ctx.Done()
	fmt.Println("ctx.Err():       ", ctx.Err())        // always the generic context.Canceled
	fmt.Println("context.Cause(): ", context.Cause(ctx)) // the actual, specific reason
}

func timeoutCauseDemo() {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 30*time.Millisecond, errors.New("SLA budget exhausted"))
	defer cancel()

	<-ctx.Done()
	fmt.Println("ctx.Err():       ", ctx.Err())
	fmt.Println("context.Cause(): ", context.Cause(ctx))
}

// withoutCancelDemo shows the classic use case: a request finishes (or is
// cancelled), but you still want to fire off a best-effort audit log entry that
// should NOT be killed just because the request context died a moment earlier.
func withoutCancelDemo() {
	reqCtx, cancel := context.WithCancel(context.Background())

	// simulate the request being cancelled (e.g. client disconnected) almost immediately
	cancel()

	// detachedCtx inherits VALUES from reqCtx, but ignores its cancellation/deadline.
	detachedCtx := context.WithoutCancel(reqCtx)
	fmt.Println("reqCtx already done?     ", reqCtx.Err() != nil)
	fmt.Println("detachedCtx done?        ", detachedCtx.Err() != nil)

	// give the detached work its own independent bound so it can't run forever
	auditCtx, auditCancel := context.WithTimeout(detachedCtx, 100*time.Millisecond)
	defer auditCancel()

	writeAuditLog(auditCtx)
}

func writeAuditLog(ctx context.Context) {
	select {
	case <-time.After(20 * time.Millisecond): // pretend to write to an audit system
		fmt.Println("audit log written successfully, even though the original request was already cancelled")
	case <-ctx.Done():
		fmt.Println("audit log write itself timed out:", ctx.Err())
	}
}
