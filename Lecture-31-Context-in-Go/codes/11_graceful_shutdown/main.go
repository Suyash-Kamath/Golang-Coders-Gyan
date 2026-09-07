// Package main demonstrates graceful HTTP server shutdown driven by context:
// signal.NotifyContext turns an OS signal into a cancelled context, and
// server.Shutdown(ctx) is itself given its own bounded context so it never
// waits forever for slow clients.
//
// This demo self-terminates after ~200ms with a SIMULATED signal so it can run
// unattended; in real usage you'd just press Ctrl+C and remove the simulation.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	})
	server := &http.Server{Addr: "127.0.0.1:8090", Handler: mux}

	go func() {
		log.Println("listening on :8090")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// --- simulate a Ctrl+C after 200ms so this example is runnable unattended ---
	go func() {
		time.Sleep(200 * time.Millisecond)
		log.Println("(demo) simulating SIGINT now — in real usage this would be an actual Ctrl+C")
		stop() // this is exactly what a real SIGINT/SIGTERM would trigger
	}()

	<-ctx.Done() // blocks here until a real signal (or, in this demo, our simulated one) arrives
	log.Println("shutdown signal received, draining in-flight requests...")

	// server.Shutdown needs its OWN context: it must not wait forever for slow
	// clients, so we give it a hard 5s budget independent of anything else.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server shut down cleanly")
}
