package junjo

import (
	"github.com/davidroman0O/junjo/data"
	"github.com/davidroman0O/seigyo"
)

/// Second attempt to make something a bit better
/// I just want a stateless server with an api (probably `fiber`) to build the templates
/// I think it's stupid that i want to elaborate on so much features, i don't need much so FUCK IT
/// We will see those interfaces later on, just need a first toy to start somewhere
/// - simple handlers endpoint api
/// - simple web app to edit my dag templates: i don't fucking care about the tech, probably will do a dirty htmx app first and see from there
/// - i don't want to do code edges behaviors or vertices behaviors, i need to have a list of built-in features, i will later for customs
/// - vertices need to have a json schema-ish validation, i just need to have inputs and outputs to those json schema, could have multiple edges in or out, edges can "plug" to properties of the json schema
/// -

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
