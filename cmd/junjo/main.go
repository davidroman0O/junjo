package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidroman0O/junjo/data"
	"github.com/davidroman0O/seigyo"
)

// What i want is to be able to launch any kind of process easily
// Here my dream list for the scope of this library
// - leverage simple dag
// - able to trigger workflows with multiple set of parameters for multiple nodes
// - able to have sub-graphs
// - able to transmit data on edge
// - programmatically it should allow a node to adapt the graph on the fly, a node is both data an function
// - a node is not business rules, a node got data and can react to it's mutations
// - provide multiple interfaces for receiving and sending data to adapt any kind of network procotol, same for storage
// - a dag mutate, data mutate, actions notify of changes, a graph can be locked, a graph can be reconciliated with another in case of distributed data
// Also what it is not
// - not flow based programming
// - workflow engine, only task engine, the dag is the rule, workers will get data to do some work and report
// - integrating all different kind of network protocols
// - no favorite way of doing things
func main() {

	// Create a channel to receive OS signals.
	sigCh := make(chan os.Signal, 1)

	// Notify sigCh when receiving SIGINT or SIGTERM signals.
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	bootstrap := seigyo.New[*data.GlobalState, interface{}](&data.GlobalState{})

	bootstrap.RegisterProcess(
		"graph",
		seigyo.ProcessConfig[*data.GlobalState, interface{}]{
			Process: &data.GraphNetProcess{
				Pid: "graph",
			},
			ShouldRecover:         true,
			RunMaxRetries:         99,
			InitMaxRetries:        99,
			DeinitMaxRetries:      99,
			MessageSendMaxRetries: 99,
		})

	errCh := bootstrap.Start()

	go func() {
		<-sigCh
		bootstrap.Stop()
	}()

	for err := range errCh {
		log.Println("error:", err)
	}

	bootstrap.Stop()

}
