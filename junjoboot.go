package junjo

import (
	"github.com/davidroman0O/junjo/data"
	"github.com/davidroman0O/seigyo"
)

type JunjoConfiguration func(j *JunjoConfig) error

type JunjoConfig struct{}

type Junjo struct {
	bootstrap *seigyo.Seigyo[*data.GlobalState, interface{}]
}

func New(cfgs ...JunjoConfiguration) (*Junjo, error) {
	cfg := JunjoConfig{} // TODO: defaults
	for i := 0; i < len(cfgs); i++ {
		if err := cfgs[i](&cfg); err != nil {
			return nil, err
		}
	}
	// TODO: do something with it
	return &Junjo{
		bootstrap: seigyo.New[*data.GlobalState, interface{}](&data.GlobalState{}),
	}, nil
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
