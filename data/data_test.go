package data

import (
	"testing"

	"github.com/davidroman0O/junjo/dag"
)

type NodeID string

type GraphNet struct {
	nodeTypes map[NodeID]Node
	graph     dag.Graph
}

func NewNet() GraphNet {
	return GraphNet{
		nodeTypes: map[NodeID]Node{},
	}
}

func (g GraphNet) AddNode(n Node) {

}

func (g GraphNet) Listen(port string) chan error {
	signal := make(chan error)

	go func() {
		for {

		}
	}()
	return signal
}

type Node interface {
	dag.Vertex
	Init()   // memory init
	Open()   // receive
	Run()    // execute
	Close()  // post
	Deinit() // memory deinit
}

type InitNode struct{}

func (i InitNode) Init()   {}
func (i InitNode) Open()   {}
func (i InitNode) Run()    {}
func (i InitNode) Close()  {}
func (i InitNode) Deinit() {}

func TestBasicGraph(t *testing.T) {

	gnet := NewNet()
	gnet.AddNode(InitNode{})

	s := gnet.Listen(":3000")

	<-s

}
