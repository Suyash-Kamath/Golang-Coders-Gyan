// Package main wires context through a realistic layered HTTP service:
// Handler -> Service -> Repository -> DB, all sharing ONE request-scoped context.
// It uses a fake in-memory "DB" that mimics database/sql's *Context methods
// (QueryRowContext) so the example runs with no external dependencies.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type fakeDB struct{}

var errNotFound = errors.New("not found")

// QueryRowContext pretends to run a slow query, but — like a real database/sql
// driver — it honors ctx cancellation instead of blindly running to completion.
func (fakeDB) QueryRowContext(ctx context.Context, userID string) (string, error) {
	select {
	case <-time.After(150 * time.Millisecond): // pretend this is real query latency
		if userID == "" {
			return "", errNotFound
		}
		return "Ada Lovelace", nil
	case <-ctx.Done():
		return "", ctx.Err() // caller cancelled or timed out before the query finished
	}
}

type Repository struct{ db fakeDB }

func (r Repository) GetUserName(ctx context.Context, userID string) (string, error) {
	return r.db.QueryRowContext(ctx, userID)
}

type Service struct{ repo Repository }

func (s Service) GetUserName(ctx context.Context, userID string) (string, error) {
	// A service layer may tighten the budget for its own downstream call, but it
	// can never make it looser than what the caller already gave it.
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	return s.repo.GetUserName(ctx, userID)
}

type Handler struct{ svc Service }

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // cancelled automatically if the client disconnects
	userID := r.URL.Query().Get("id")

	name, err := h.svc.GetUserName(ctx, userID)
	switch {
	case err == nil:
		fmt.Fprintf(w, "user %s: %s\n", userID, name)
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "upstream timed out", http.StatusGatewayTimeout)
	case errors.Is(err, context.Canceled):
		// the client already went away — nobody will read this response, but log it
		log.Println("request cancelled by client for user", userID)
	case errors.Is(err, errNotFound):
		http.Error(w, "user not found", http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func main() {
	h := Handler{svc: Service{repo: Repository{db: fakeDB{}}}}
	server := &http.Server{Addr: "127.0.0.1:8089", Handler: h}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	time.Sleep(50 * time.Millisecond) // let the server come up

	fmt.Println("=== request that completes normally (150ms query fits in the 300ms service budget) ===")
	get("http://127.0.0.1:8089/?id=42")

	fmt.Println("\n=== request with a client-side timeout shorter than the query needs — the server sees ctx.Canceled ===")
	getWithClientTimeout("http://127.0.0.1:8089/?id=42", 50*time.Millisecond)
	time.Sleep(50 * time.Millisecond) // let the server-side log line print before we close it

	_ = server.Close()
}

func get(url string) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("client error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("status:", resp.Status)
}

func getWithClientTimeout(url string, d time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("client error (expected — client gave up before the server could finish):", err)
		return
	}
}
