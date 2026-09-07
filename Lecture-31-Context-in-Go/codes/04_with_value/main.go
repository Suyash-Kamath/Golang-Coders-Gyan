// Package main shows the safe pattern for context.WithValue: typed unexported
// keys, typed accessor functions, upward-walking lookup, and what should (and
// should not) ever be stashed in a context.
package main

import (
	"context"
	"fmt"
)

// --- Step 1: never use a raw string (or other exported/basic type) as a context key. ---
// Two packages could both use the string "requestID" and silently collide/shadow each other.
// The fix: an unexported struct type that only this package can construct.

type contextKey struct{ name string }

var requestIDKey = contextKey{name: "requestID"}
var userIdentityKey = contextKey{name: "userIdentity"}

// --- Step 2: wrap Get/Set in typed helper functions so callers never touch the key directly. ---

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

type UserIdentity struct {
	UserID string
	Role   string
}

func WithUserIdentity(ctx context.Context, u UserIdentity) context.Context {
	return context.WithValue(ctx, userIdentityKey, u)
}

func UserIdentityFromContext(ctx context.Context) (UserIdentity, bool) {
	u, ok := ctx.Value(userIdentityKey).(UserIdentity)
	return u, ok
}

func main() {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-abc-123")
	ctx = WithUserIdentity(ctx, UserIdentity{UserID: "u42", Role: "admin"})

	handleRequest(ctx)

	fmt.Println("\n=== context.Value lookup walks UP the chain ===")
	// WithValue never mutates its parent — it wraps it. Value() checks itself, then
	// delegates to the parent, and so on, until it finds a match or runs out of chain.
	base := context.WithValue(context.Background(), contextKey{"a"}, "A")
	mid := context.WithValue(base, contextKey{"b"}, "B")
	leaf := context.WithValue(mid, contextKey{"c"}, "C")
	fmt.Println("leaf sees a:      ", leaf.Value(contextKey{"a"})) // found by walking up to `base`
	fmt.Println("base sees c (nil):", base.Value(contextKey{"c"})) // base was created BEFORE c existed

	fmt.Println("\n=== What should and shouldn't live in context.Value ===")
	goodAndBadExamples()
}

func handleRequest(ctx context.Context) {
	if id, ok := RequestIDFromContext(ctx); ok {
		fmt.Println("handling request:", id)
	}
	if u, ok := UserIdentityFromContext(ctx); ok {
		fmt.Printf("authenticated as %s (role=%s)\n", u.UserID, u.Role)
	}
	processOrder(ctx, "order-777")
}

// processOrder is two layers deep and never received the requestID as an explicit
// parameter, yet it can still log/trace it — that's the whole point of WithValue.
func processOrder(ctx context.Context, orderID string) {
	id, _ := RequestIDFromContext(ctx)
	fmt.Printf("[req=%s] processing %s\n", id, orderID)
}

func goodAndBadExamples() {
	fmt.Println("GOOD: request-scoped, cross-cutting metadata")
	fmt.Println("  - request ID / correlation ID")
	fmt.Println("  - trace / span ID (tracing systems)")
	fmt.Println("  - authenticated user identity")
	fmt.Println("  - deadline/cancellation (built in)")

	fmt.Println("BAD: things that should just be normal function arguments")
	fmt.Println("  - *sql.DB or other long-lived service objects")
	fmt.Println("  - application config")
	fmt.Println("  - business parameters like orderID, amount, etc.")
	fmt.Println("  - a *log.Logger you depend on for control flow")
	// e.g. do NOT do: ctx = context.WithValue(ctx, dbKey, db)
	// instead:        func ProcessOrder(ctx context.Context, db *sql.DB, orderID string)
}
