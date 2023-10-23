package types

import (
	"fmt"

	"github.com/davidroman0O/junjo/dag"
)

type WorkUnitDag struct {
	graph                 *dag.AcyclicGraph
	storageImplementation StorageInterface
}

func NewWorkUnitDag(storageImplementation StorageInterface) *WorkUnitDag {
	return &WorkUnitDag{
		graph:                 &dag.AcyclicGraph{},
		storageImplementation: storageImplementation,
	}
}

func (d *WorkUnitDag) Graph() *dag.AcyclicGraph {
	return d.graph
}

type TaskUnitFactory func() dag.Vertex

// Will create a TaskUnitFactory to help you declare connections instead of creating them manually on each vertex
func (d *WorkUnitDag) AddTaskDefinition(taskDescription *TaskDefinition) TaskUnitFactory {
	return func() dag.Vertex {
		uuid, _ := d.storageImplementation.NewUUID()
		node := NewNodeTaskUnit(
			WithNodeWithTaskKey(TaskUnitID(uuid)),
			WithNodeWithTaskStatus(NoneStatus), // make sure that all Node have at least a status
			WithNodeWithTaskDefinition(taskDescription),
		)
		return d.graph.Add(node)
	}
}

// Create a new Vertex based on a TaskUnit, you will have to manually connect each vertexes
// If the node has no owner, it won't be visible by any
func (d *WorkUnitDag) AddTaskUnit(cfgs ...NodeTaskUnitConfig) dag.Vertex {
	cfgs = append(cfgs, WithNodeWithTaskStatus(NoneStatus)) // make sure that all Node have at least a status
	node := NewNodeTaskUnit(cfgs...)
	return d.graph.Add(node)
}

// `Connect` will MANUALLY connect one vertex to another vertex at a time
func (d *WorkUnitDag) Connect(from dag.Vertex, to dag.Vertex) {
	d.graph.Connect(dag.BasicEdge(from, to))
}

// `ConnectDef` will help CREATING vertexes based on a `TaskUnitConnector` which is based on a description
func (d *WorkUnitDag) ConnectDef(from TaskUnitFactory, to TaskUnitFactory) (dag.Vertex, dag.Vertex) {
	source := from()
	target := to()
	d.graph.Connect(dag.BasicEdge(source, target))
	return source, target
}

func (d *WorkUnitDag) AssociateDef(from dag.Vertex, to TaskUnitFactory) dag.Vertex {
	target := to()
	d.graph.Connect(dag.BasicEdge(from, target))
	return target
}

func (d *WorkUnitDag) MConnect(target dag.Vertex, sources ...dag.Vertex) dag.Vertex {
	for i := 0; i < len(sources); i++ {
		d.graph.Connect(dag.BasicEdge(sources[i], target))
	}
	return target
}

func (d *WorkUnitDag) MConnectDef(def TaskUnitFactory, sources ...dag.Vertex) dag.Vertex {
	target := def()
	for i := 0; i < len(sources); i++ {
		d.graph.Connect(dag.BasicEdge(sources[i], target))
	}
	return target
}

func (d *WorkUnitDag) Print() {
	fmt.Println(d.graph.Graph.StringWithNodeTypes())
}

func (d *WorkUnitDag) ToTaskUnits() ([]*TaskUnit, error) {
	var taskUnits []*TaskUnit = []*TaskUnit{}

	vertices := d.graph.Vertices()
	for _, vertex := range vertices {
		nodeTaskUnit, ok := vertex.(*NodeTaskUnit) // Assuming vertex type is *NodeTaskUnit
		if !ok {
			return nil, fmt.Errorf("vertex is not a NodeTaskUnit")
		}

		taskUnit := NewTaskUnit(
			nodeTaskUnit.Unit.Key,
			WithTaskUnitDefinition(nodeTaskUnit.Definition),
			WithTaskUnitStatus(nodeTaskUnit.Unit.Status),
			WithTaskUnitData(nodeTaskUnit.Unit.Data),
		)

		// Collect dependencies
		immediateAncestors, err := d.graph.ImmediateAncestors(vertex)
		if err != nil {
			return nil, err // Handle error
		}
		for _, ancestorVertex := range immediateAncestors {
			ancestor, ok := ancestorVertex.(*NodeTaskUnit)
			if !ok {
				return nil, fmt.Errorf("ancestor vertex is not a NodeTaskUnit")
			}
			taskUnit.DependsOnIDs = append(taskUnit.DependsOnIDs, TaskUnitID(ancestor.Unit.Key))
		}

		taskUnits = append(taskUnits, taskUnit)
	}

	// We need to re-sync the taskunits for their ids and pointers
	// It will create nested pointers lists, which would eventually fuck up your `pp` logs, don't care :)
	for idx := 0; idx < len(taskUnits); idx++ {
		for ididx := 0; ididx < len(taskUnits[idx].DependsOnIDs); ididx++ {
			for jdx := 0; jdx < len(taskUnits); jdx++ {
				if taskUnits[idx].DependsOnIDs[ididx] == taskUnits[jdx].Key {
					// handshake bros
					// it will stay only at runtime
					taskUnits[idx].DependsOn = append(taskUnits[idx].DependsOn, taskUnits[jdx])
				}
			}
		}
	}

	return taskUnits, nil
}

func (d *WorkUnitDag) AvailableNodeUnit() []NodeTaskUnit {
	var availableUnits []NodeTaskUnit
	for _, vertex := range d.graph.Vertices() {
		node, ok := vertex.(*NodeTaskUnit)
		if !ok || node.Unit.Status != NoneStatus {
			continue
		}
		immediateAncestors, err := d.graph.ImmediateAncestors(vertex)
		if err != nil {
			// todo: handle error
			continue
		}
		allAncestorsSuccess := true
		for _, ancestorVertex := range immediateAncestors {
			ancestor, ok := ancestorVertex.(*NodeTaskUnit)
			if !ok || ancestor.Unit.Status != SuccessStatus {
				allAncestorsSuccess = false
				break
			}
		}
		if allAncestorsSuccess {
			availableUnits = append(availableUnits, *node)
		}
	}
	return availableUnits
}

func (d *WorkUnitDag) CanChangeStatus(vertex dag.Vertex) (bool, error) {
	// Cast vertex to *NodeUnit to get the associated NodeUnit
	targetNode, ok := vertex.(*NodeTaskUnit)
	if !ok {
		return false, fmt.Errorf("vertex is not a NodeUnit")
	}

	// Check the current status of the NodeUnit
	switch targetNode.Unit.Status {
	case NoneStatus, QueuedStatus, ProgressStatus, PauseStatus:
		// Allowed statuses, proceed to next step
	default:
		return false, nil
	}

	// Check if all predecessors are in SuccessStatus
	predecessors, err := d.graph.ImmediateAncestors(vertex)
	if err != nil {
		return false, err
	}
	for _, predecessorVertex := range predecessors {
		predecessor, ok := predecessorVertex.(*NodeTaskUnit)
		if !ok || predecessor.Unit.Status != SuccessStatus {
			return false, nil // One of the predecessors is not in SuccessStatus, status cannot be changed
		}
	}

	return true, nil // All checks passed, status can be changed
}

func (d *WorkUnitDag) AvailableNodeUnitWithOwner(ownerID OwnerID) []NodeTaskUnit {
	var availableUnits []NodeTaskUnit
	for _, vertex := range d.graph.Vertices() {
		node, ok := vertex.(*NodeTaskUnit)
		if !ok || node.Definition.Key == "" {
			continue
		}
		if !ok || node.Unit.Status != NoneStatus || node.Definition.OwnerID != ownerID {
			continue
		}
		immediateAncestors, err := d.graph.ImmediateAncestors(vertex)
		if err != nil {
			// handle error (optional based on your error handling strategy)
			continue
		}
		allAncestorsSuccess := true
		for _, ancestorVertex := range immediateAncestors {
			ancestor, ok := ancestorVertex.(*NodeTaskUnit)
			if !ok || (ancestor.Unit.Status != NoneStatus) {
				allAncestorsSuccess = false
				break
			}
		}
		if allAncestorsSuccess {
			availableUnits = append(availableUnits, *node)
		}
	}
	return availableUnits
}

func (d *WorkUnitDag) CanChangeStatusWithOwner(vertex dag.Vertex, ownerID OwnerID) (bool, error) {
	// Cast vertex to *NodeUnit to get the associated NodeUnit
	targetNode, ok := vertex.(*NodeTaskUnit)
	if !ok || targetNode.Definition.Key == "" {
		return false, fmt.Errorf("vertex is does not have description ownership")
	}
	if !ok || targetNode.Definition.OwnerID != ownerID {
		return false, nil
	}

	// Check the current status of the NodeUnit
	switch targetNode.Unit.Status {
	case NoneStatus, QueuedStatus, ProgressStatus, PauseStatus:
		// Allowed statuses, proceed to next step
	default:
		return false, nil
	}

	// Check if all immediate ancestors are in SuccessStatus
	immediateAncestors, err := d.graph.ImmediateAncestors(vertex)
	if err != nil {
		return false, err
	}
	for _, ancestorVertex := range immediateAncestors {
		ancestor, ok := ancestorVertex.(*NodeTaskUnit)
		if !ok || ancestor.Unit.Status != SuccessStatus {
			return false, nil // One of the ancestors is not in SuccessStatus, status cannot be changed
		}
	}

	return true, nil // All checks passed, status can be changed
}

// If the nodes has no owner, they won't be visible by any
func CreateDagFromTaskUnits(si StorageInterface, taskUnits []TaskUnit, definitions []TaskDefinition) (*WorkUnitDag, error) {
	// Create a new Dag
	wdag := NewWorkUnitDag(si)

	// Create a map to hold the NodeUnit instances and to allow for quick lookup by TaskUnitID
	nodeUnitMap := make(map[TaskUnitID]dag.Vertex)

	var def *TaskDefinition
	// Step 1: Add all nodes to the Dag and to the nodeUnitMap
	for idxTask := 0; idxTask < len(taskUnits); idxTask++ {

		for idxDef := 0; idxDef < len(definitions); idxDef++ {
			if definitions[idxDef].Key == taskUnits[idxTask].TaskDefinitionID {
				def = &definitions[idxDef]
			}
		}

		nodeUnitCfgs := []NodeTaskUnitConfig{
			WithNodeWithTaskKey(taskUnits[idxTask].Key),
			WithNodeWithTaskDefinition(def),
			WithNodeWithTaskStatus(taskUnits[idxTask].Status),
			WithNodeWithTaskData(taskUnits[idxTask].Data),
		}

		nodeUnit := wdag.AddTaskUnit(nodeUnitCfgs...)
		nodeUnitMap[taskUnits[idxTask].Key] = nodeUnit
	}

	// Step 2: Connect the nodes in the Dag based on the DependsOnIDs field
	for idxTaskUnit := 0; idxTaskUnit < len(taskUnits); idxTaskUnit++ {
		fromNode, exists := nodeUnitMap[taskUnits[idxTaskUnit].Key]
		if !exists {
			return nil, fmt.Errorf("node doesnt exists for TaskUnitID: %s", taskUnits[idxTaskUnit].Key)
		}
		for _, dependsOnID := range taskUnits[idxTaskUnit].DependsOnIDs {
			toNode, exists := nodeUnitMap[dependsOnID]
			if !exists {
				return nil, fmt.Errorf("node not found for DependsOnID: %s", dependsOnID)
			}
			wdag.Connect(fromNode, toNode)
		}
	}

	vertices := wdag.graph.Vertices()

	for idxVertice := 0; idxVertice < len(vertices); idxVertice++ {
		_, ok := vertices[idxVertice].(*NodeTaskUnit)
		if !ok {
			fmt.Println("not ok for some reason")
		}
	}

	return wdag, nil
}
