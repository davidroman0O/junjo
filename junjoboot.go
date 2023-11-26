package junjo

import (
	"github.com/davidroman0O/junjo/data"
	"github.com/davidroman0O/seigyo"
)

type Junjo struct {
	bootstrap *seigyo.Seigyo[*data.GlobalState, interface{}]
}

func New() Junjo {
	return Junjo{
		bootstrap: seigyo.New[*data.GlobalState, interface{}](&data.GlobalState{}),
	}
}

func (j *Junjo) Start() <-chan error {
	// add all processes
	j.bootstrap.RegisterProcess(
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

	return j.bootstrap.Start()
}

func (j *Junjo) Stop() {
	if j.bootstrap != nil {
		j.bootstrap.Stop()
	}
}
