package goraph

// The Graph interface is implemented by all graph types.
type Graph interface {
	NewVertex() Vertex
	AddEdge(v1, v2 Vertex)
	Vertices() []Vertex
	Edges() []Edge
}

var (
	_ Graph = &DirectedGraph{}
	_ Graph = &AdjacencyList{}
)

// Vertex represents a node in the graph. Users should create
// new Vertex values with NewVertex.
type Vertex int

// Edge represents an edge between two vertices.
// In a directed graph the edge is from v1 to v2.
type Edge struct{ v1, v2 Vertex }

// AdjacencyList implements an undirected graph using an adjacency list.
type AdjacencyList struct {
	edges      map[Vertex][]Vertex
	nextVertex Vertex
}

// NewAdjacencyList creates an empty graph.
func NewAdjacencyList() *AdjacencyList {
	return &AdjacencyList{edges: make(map[Vertex][]Vertex)}
}

// NewVertex adds a new vertex.
func (g *AdjacencyList) NewVertex() Vertex {
	v := g.nextVertex
	g.edges[v] = make([]Vertex, 0)
	g.nextVertex++
	return v
}

// AddEdge adds an edge between v1 and v2.2.
func (g *AdjacencyList) AddEdge(v1, v2 Vertex) {
	if v2 < v1 {
		v1, v2 = v2, v1
	}
	edges := g.edges[v1]
	g.edges[v1] = append(edges, v2)
}

// Vertices returns a slice of all vertices.
func (g *AdjacencyList) Vertices() []Vertex {
	vertices := make([]Vertex, len(g.edges))
	var i int
	for k := range g.edges {
		vertices[i] = k
		i++
	}
	return vertices
}

// Edges returns a slice of all edges.
func (g *AdjacencyList) Edges() []Edge {
	var edges []Edge
	for k, neighbors := range g.edges {
		for _, n := range neighbors {
			edges = append(edges, Edge{k, n})
		}
	}
	return edges
}

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
	vertices := make([]Vertex, len(g.edges))
	var i int
	for k := range g.edges {
		vertices[i] = k
		i++
	}
	return vertices
}

// Edges returns all the outgoing edges of the graph.
func (g *DirectedGraph) Edges() []Edge {
	var edges []Edge
	for k, neighbors := range g.edges {
		for _, n := range neighbors {
			edges = append(edges, Edge{k, n})
		}
	}
	return edges
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
