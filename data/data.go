package data

import (
	"context"
	"log"
	"time"
)

type GlobalState struct {
}

// `GraphNetProcess` represent a permanent process that will wait for instructions from client
type GraphNetProcess struct {
	Pid string
}

func (p *GraphNetProcess) Init(ctx context.Context, stateGetter func() *GlobalState, stateMutator func(mutateFunc func(*GlobalState) *GlobalState), sender func(pid string, data interface{})) error {
	log.Println("graphnet initializing")
	return nil
}

func (p *GraphNetProcess) Run(ctx context.Context, stateGetter func() *GlobalState, stateMutator func(mutateFunc func(*GlobalState) *GlobalState), sender func(pid string, data interface{}), shutdownCh chan struct{}, errCh chan<- error, selfShutdown func()) error {
	log.Println("graphnet running")

	t := time.NewTimer(time.Second * 2)

	select {
	case <-t.C:
		log.Println(p.Pid, "reached timeout")
	}

	log.Println("graphnet shutdown")
	return nil
}

func (p *GraphNetProcess) Deinit(ctx context.Context, stateGetter func() *GlobalState, stateMutator func(mutateFunc func(*GlobalState) *GlobalState), sender func(pid string, data interface{})) error {
	log.Println("graphnet deinitialization")
	return nil
}

func (p *GraphNetProcess) Received(pid string, data interface{}) error {
	switch data.(type) {

	}
	return nil
}
