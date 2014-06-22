package goraph

// Vertex represents a node in the graph. Users should create
// new Vertex values with NewVertex.
type Vertex int

// DirectedGraph provides a space-efficient directed graph.
type DirectedGraph struct {
	edges      map[Vertex][]Vertex
	nextVertex Vertex
}

// NewDirectedGraph creates and initializes a DirectedGraph.
func NewDirectedGraph() *DirectedGraph {
	g := &DirectedGraph{edges: make(map[Vertex][]Vertex)}
	return g
}

// addVertex adds v to the graph in an idempotent fashion. The return value
// indicates whether or not the vertex was already in the graph; if false,
// the value was not in the graph before it was added.
func (g *DirectedGraph) addVertex(v Vertex) bool {
	_, ok := g.edges[v]
	if !ok {
		g.edges[v] = make([]Vertex, 0)
	}
	return ok
}

// AddEdge connects vertices v1 and v2 in the graph.
func (g *DirectedGraph) AddEdge(v1, v2 Vertex) {
	g.addVertex(v1)
	g.addVertex(v2)
	g.edges[v1] = append(g.edges[v1], v2)
}

// NewVertex creates a new Vertex, adds it to the graph, and returns it.
func (g *DirectedGraph) NewVertex() Vertex {
	v := g.nextVertex
	g.addVertex(v)
	g.nextVertex++
	return v
}

// Vertices returns a slice of the vertices that are in the graph.
func (g *DirectedGraph) Vertices() []Vertex {
	result := make([]Vertex, len(g.edges), len(g.edges))
	i := 0
	for k := range g.edges {
		result[i] = k
		i++
	}
	return result
}

// incomingEdges finds the vertices that connect to v
func (g *DirectedGraph) incomingEdges(v Vertex) []Vertex {
	result := make([]Vertex, 0)
	for w, vlist := range g.edges {
		for _, x := range vlist {
			if v == x {
				result = append(result, w)
				break
			}
		}
	}
	return result
}

// countIncomingEdges is like incomingEdges but only delivers a count.
func (g *DirectedGraph) countIncomingEdges(v Vertex) int {
	result := 0
	for _, vlist := range g.edges {
		for _, x := range vlist {
			if v == x {
				result += 1
				break
			}
		}
	}
	return result
}
