// Package main is a miniature, educational re-implementation of the ideas behind
// context.Context — just enough to show HOW cancellation propagation actually
// works under the hood (no Value()/deadline support here — see the real
// standard library for that). The key trick: a parent keeps a registry of its
// children, so ONE parent.cancel() call reaches every descendant without
// spawning a goroutine per context.
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var Canceled = errors.New("canceled")

type MyContext struct {
	mu       sync.Mutex
	done     chan struct{}
	err      error
	children map[*MyContext]struct{}
}

// NewRoot creates a root context, exactly like context.Background().
func NewRoot() *MyContext {
	return &MyContext{
		done:     make(chan struct{}),
		children: make(map[*MyContext]struct{}),
	}
}

// WithCancel derives a child from parent and registers it in the parent's
// children map — the registry that lets cancellation propagate downward.
func WithCancel(parent *MyContext) (*MyContext, func()) {
	child := &MyContext{
		done:     make(chan struct{}),
		children: make(map[*MyContext]struct{}),
	}

	parent.mu.Lock()
	if parent.err != nil {
		// parent is already dead — cancel the child immediately, no need to register it
		parent.mu.Unlock()
		child.cancel(parent.err)
		return child, func() { child.cancel(Canceled) }
	}
	parent.children[child] = struct{}{}
	parent.mu.Unlock()

	return child, func() { child.cancel(Canceled) }
}

func (c *MyContext) cancel(err error) {
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // idempotent: already cancelled, calling again is a no-op
	}
	c.err = err
	close(c.done) // closing (not sending) broadcasts to every goroutine waiting on Done()
	children := c.children
	c.children = nil
	c.mu.Unlock()

	// propagate downward to every registered child, recursively
	for child := range children {
		child.cancel(err)
	}
}

func (c *MyContext) Done() <-chan struct{} { return c.done }

func (c *MyContext) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func main() {
	root := NewRoot()
	child, cancelChild := WithCancel(root)
	grandchild, _ := WithCancel(child)

	go func() {
		<-grandchild.Done()
		fmt.Println("grandchild cancelled:", grandchild.Err())
	}()

	fmt.Println("cancelling the CHILD (not the root)...")
	cancelChild()
	time.Sleep(20 * time.Millisecond)
	fmt.Println("root itself is still alive:", root.Err() == nil)

	fmt.Println("\ncancelling the ROOT now propagates through any remaining descendants")
	rootChild, _ := WithCancel(root)
	rootGrandchild, _ := WithCancel(rootChild)
	go func() {
		<-rootGrandchild.Done()
		fmt.Println("rootGrandchild cancelled:", rootGrandchild.Err())
	}()
	root.cancel(Canceled)
	time.Sleep(20 * time.Millisecond)
}
