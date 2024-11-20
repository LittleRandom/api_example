package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"plainrandom/config"
	"plainrandom/server"
)

func main() {

	// Setting up a signal handler to receive kill signal.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)   // Single buffer that takes os.Signal items
	signal.Notify(c, os.Interrupt) // os.Interrupt is CTRL + C signal, will notify the c channel

	// Start a goroutine to read channel and instantiate cancel functions on os trigger.
	go func() { <-c; cancel() }()

	// Execute program.
	log.Printf("Starting Application...")
	//
	// Creates a new Main object
	conf := config.NewConfig()
	s := server.NewServer(conf)
	//
	// Starts the API server
	s.Start()
	//
	// Line to wait for CTRL-C
	<-ctx.Done()

	// Start shutting down
	log.Printf("Shutting down Server...")
	if err := s.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	log.Printf("main: done. exiting")
}
