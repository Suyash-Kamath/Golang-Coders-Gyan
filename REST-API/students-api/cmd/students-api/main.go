package main

import (
	"context"
	// "fmt"
	// "fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/config"
)

func main() {

	// fmt.Println("Welcome to Student's API")

	// Load config
	cfg := config.MustLoad()

	// Custom logger is want to set

	// Database setup

	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Student's API"))
	})

	// Setup server

	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	// fmt.Println("Server Started")
	// this below line is blocking call, so if there is an error in starting the server, it will be logged and the program will exit
	// err:=server.ListenAndServe()
	// if err != nil && err != http.ErrServerClosed {
	// 	log.Fatalf("Failed to Start Server: %s", err.Error())
	// }

	// Production me itna simple nahi banate , iske upar aur ek chiz add karni hai which is graceful shutdown, so that when we want to stop the server, it will wait for all the ongoing requests to complete before shutting down the server,

	slog.Info("Server Started at ", slog.String("address", cfg.HTTPServer.Addr))

	// this is blocking call, so if there is an error in starting the server, it will be logged and the program will exit

	// also there is one problem , production me itna simple nahi banate , iske upar aur ek chiz add karni hai which is graceful shutdown, so that when we want to stop the server, it will wait for all the ongoing requests to complete before shutting down the server,
	// first ham alag go routine banalete hai

	done := make(chan os.Signal, 1)
	// signal by default will send os.Interrupt signal when we press ctrl+c in the terminal, so we are listening for that signal and when we receive that signal, we will send a signal to the done channel which will be used to stop the server gracefully
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to Start Server: %s", err.Error())
		}

	}()

	<-done // with this , main goroutine will wait here until it receives a signal from the done channel, which will be sent when we want to stop the server

	// Now we will stop the server gracefully, so that it will wait for all the ongoing requests to complete before shutting down the server, for that we will use the Shutdown method of the http.Server which takes a context as an argument, so we will create a context with timeout of 5 seconds, so that if there are any ongoing requests that are taking more than 5 seconds to complete, then the server will forcefully shut down after 5 seconds
	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // this will cancel the context after 5 seconds, so that the server will not wait for more than 5 seconds to complete the ongoing requests

	// you can make this one liner
	// err:= server.Shutdown(ctx)
	// if err!=nil{
	// 	slog.Error("Failed to Shutdown Server",slog.String("error",err.Error()))
	// }

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to Shutdown Server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}
