package main

import (
	"MicroService-Go/handlers"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	logfile, err := os.OpenFile(".\\microservice-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error: Unable to open logfile\n")
	}
	defer logfile.Close()

	log_wrt := io.MultiWriter(os.Stdout, logfile)
	logger := log.New(log_wrt, "[Production-API] ", log.LstdFlags)

	goodbyeHandler := handlers.NewGoodbye(logger)
	productHandler := handlers.NewProducts(logger)
	logger.Println("Starting Microservice Application")

	serverMux := http.NewServeMux()
	serverMux.Handle("/goodbye", goodbyeHandler)
	serverMux.Handle("/products", productHandler)
	serverMux.Handle("/products/", productHandler)

	srv := &http.Server{
		Addr:         ":9090",
		Handler:      serverMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(tc)
}
