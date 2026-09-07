// Package main shows the standard HTTP middleware pattern for context: each
// middleware reads r.Context(), derives a new context with WithValue, and calls
// next.ServeHTTP with r.WithContext(ctx) — a NEW *http.Request carrying it.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

type ctxKey struct{ name string }

var requestIDKey = ctxKey{"requestID"}
var userKey = ctxKey{"user"}

// RequestIDMiddleware stamps every incoming request with a unique ID and makes it
// available to every downstream handler via the request's context.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := fmt.Sprintf("req-%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddleware pretends to validate a token and stores the resulting identity.
// Any handler further down the chain can read it back out with no extra plumbing.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := "anonymous"
		if token == "Bearer secret-token" {
			user = "ada"
		}
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := ctx.Value(requestIDKey).(string)
	user, _ := ctx.Value(userKey).(string)
	fmt.Fprintf(w, "request %s handled for user %s\n", id, user)
}

func main() {
	handler := RequestIDMiddleware(AuthMiddleware(http.HandlerFunc(finalHandler)))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	fmt.Println("=== anonymous request ===")
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, _ := http.DefaultClient.Do(req)
	printBody(resp)

	fmt.Println("=== authenticated request ===")
	req2, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req2.Header.Set("Authorization", "Bearer secret-token")
	resp2, _ := http.DefaultClient.Do(req2)
	printBody(resp2)
}

func printBody(resp *http.Response) {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Print(string(body))
	fmt.Println("X-Request-ID header:", resp.Header.Get("X-Request-ID"))
}
