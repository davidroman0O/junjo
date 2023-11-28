package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidroman0O/junjo"
)

func main() {

	// Create a channel to receive OS signals.
	sigCh := make(chan os.Signal, 1)

	// Notify sigCh when receiving SIGINT or SIGTERM signals.
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	var bootstrap *junjo.Junjo
	var err error

	if bootstrap, err = junjo.New(); err != nil {
		panic(err)
	}

	errCh := bootstrap.Start()

	go func() {
		fmt.Println("waiting signal")
		<-sigCh
		bootstrap.Stop()
	}()

	for err = range errCh {
		log.Println("error:", err)
	}

	fmt.Println("close bootstrap")
	bootstrap.Stop()
}
