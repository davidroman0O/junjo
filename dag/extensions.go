package dag

// Predecessors returns the immediate predecessors of a given vertex.
// @droman: required those two functions
func (g *AcyclicGraph) ImmediateAncestors(me Vertex) ([]Vertex, error) {
	var predecessors []Vertex
	for _, edge := range g.Edges() {
		if edge.Target() == me {
			predecessors = append(predecessors, edge.Source())
		}
	}
	return predecessors, nil
}

// Descendants returns the immediate descendants of a given vertex.
// @droman: required those two functions
func (g *AcyclicGraph) ImmediateDescendants(me Vertex) ([]Vertex, error) {
	var descendants []Vertex
	for _, edge := range g.Edges() {
		if edge.Source() == me {
			descendants = append(descendants, edge.Target())
		}
	}
	return descendants, nil
}
