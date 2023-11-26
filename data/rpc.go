package data

import (
	"context"
	"log"
)

// `RaftProcess` represent a permanent process that will wait for instructions from client
type RaftProcess struct {
	Pid string
}

func (p *RaftProcess) Init(ctx context.Context, stateGetter func() *GlobalState, stateMutator func(mutateFunc func(*GlobalState) *GlobalState), sender func(pid string, data interface{})) error {
	log.Println("agent initializing")
	return nil
}

func (p *RaftProcess) Run(ctx context.Context, stateGetter func() *GlobalState, stateMutator func(mutateFunc func(*GlobalState) *GlobalState), sender func(pid string, data interface{}), shutdownCh chan struct{}, errCh chan<- error, selfShutdown func()) error {
	log.Println("agent running")

	log.Println("agent shutdown")
	return nil
}

func (p *RaftProcess) Deinit(ctx context.Context, stateGetter func() *GlobalState, stateMutator func(mutateFunc func(*GlobalState) *GlobalState), sender func(pid string, data interface{})) error {
	log.Println("agent deinitialization")
	return nil
}

func (p *RaftProcess) Received(pid string, data interface{}) error {
	switch data.(type) {

	}
	return nil
}
