package main

import (
	"MicroService-Go/handlers"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	logfile, err := os.OpenFile(".\\microservice-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error: Unable to open logfile\n")
	}
	defer logfile.Close()

	log_wrt := io.MultiWriter(os.Stdout, logfile)
	logger := log.New(log_wrt, "[Production-API] ", log.LstdFlags)

	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)
	logger.Println("Starting Microservice Application")

	serverMux := http.NewServeMux()
	serverMux.Handle("/", helloHandler)
	serverMux.Handle("/goodbye", goodbyeHandler)
	log.Fatal(http.ListenAndServe(":9090", serverMux))
}
