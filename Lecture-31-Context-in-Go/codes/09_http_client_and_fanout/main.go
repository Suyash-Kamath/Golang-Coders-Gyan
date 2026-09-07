// Package main fans a single page's data needs out across three independent,
// concurrent calls. If any one of them fails, there is no point letting the
// others keep working — so the group cancels itself as soon as the first
// error appears. Shown two ways: manually, and with a tiny errgroup-style helper
// (the real thing lives in golang.org/x/sync/errgroup — not a stdlib package).
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func fetchUser(ctx context.Context) (string, error) {
	select {
	case <-time.After(50 * time.Millisecond):
		return "user:ada", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func fetchOrders(ctx context.Context) (string, error) {
	select {
	case <-time.After(30 * time.Millisecond):
		return "", errors.New("orders service unavailable") // this one fails
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func fetchRecommendations(ctx context.Context) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond): // slow — should get cancelled before finishing
		return "recommendations:...", nil
	case <-ctx.Done():
		fmt.Println("fetchRecommendations: cancelled early because a sibling failed:", ctx.Err())
		return "", ctx.Err()
	}
}

func main() {
	fmt.Println("=== manual fan-out with a shared cancel ===")
	manualFanOut()

	fmt.Println("\n=== same thing using a tiny errgroup-style helper ===")
	groupFanOut()
}

func manualFanOut() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	results := make(map[string]string)
	tasks := map[string]func(context.Context) (string, error){
		"user":            fetchUser,
		"orders":          fetchOrders,
		"recommendations": fetchRecommendations,
	}

	wg.Add(len(tasks))
	for name, task := range tasks {
		go func(name string, task func(context.Context) (string, error)) {
			defer wg.Done()
			val, err := task(ctx)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("%s: %w", name, err)
					cancel() // stop every other in-flight task
				}
				return
			}
			results[name] = val
		}(name, task)
	}
	wg.Wait()

	if firstErr != nil {
		fmt.Println("fan-out failed:", firstErr)
	}
	fmt.Println("partial results:", results)
}

// group is a minimal stand-in for golang.org/x/sync/errgroup.Group, just enough
// to show what errgroup.WithContext(ctx) buys you: the first error cancels the
// shared context so every sibling goroutine can stop early.
type group struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.Mutex
	err    error
}

func withGroupContext(ctx context.Context) (*group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &group{cancel: cancel}, ctx
}

func (g *group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.mu.Lock()
			if g.err == nil {
				g.err = err
				g.cancel() // first error wins and cancels every sibling's context
			}
			g.mu.Unlock()
		}
	}()
}

func (g *group) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

func groupFanOut() {
	g, ctx := withGroupContext(context.Background())

	var mu sync.Mutex
	results := make(map[string]string)

	g.Go(func() error {
		v, err := fetchUser(ctx)
		if err != nil {
			return err
		}
		mu.Lock()
		results["user"] = v
		mu.Unlock()
		return nil
	})
	g.Go(func() error {
		v, err := fetchOrders(ctx)
		if err != nil {
			return fmt.Errorf("orders: %w", err)
		}
		mu.Lock()
		results["orders"] = v
		mu.Unlock()
		return nil
	})
	g.Go(func() error {
		v, err := fetchRecommendations(ctx)
		if err != nil {
			return err
		}
		mu.Lock()
		results["recommendations"] = v
		mu.Unlock()
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println("group failed:", err)
	}
	fmt.Println("partial results:", results)
}
