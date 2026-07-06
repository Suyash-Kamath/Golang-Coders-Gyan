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
	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/http/handlers/student"
	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/storage/sqlite"
)

func main() {

	// fmt.Println("Welcome to Student's API")

	// Load config
	cfg := config.MustLoad()

	// Custom logger is want to set

	// Database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	} // db connect nahi ho raha toh application chalake kya faayda

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// Setup router
	router := http.NewServeMux()

	// this is the convention of restapi , and plural rakhte hai resources ko
	router.HandleFunc("POST /api/students", student.New(storage))
	// student.New() me hame db ko use karna hai , so as a dependency receive karna padega , so that plug in play waali hai
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PUT /api/students/{id}", student.Update(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.Delete(storage))
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
