package goraph

import (
	"sort"
)

// The Graph interface is implemented by all graph types.
type Graph interface {
	NewVertex() Vertex
	DeleteVertex(v Vertex)
	AddEdge(v1, v2 Vertex)
	RemoveEdge(v1, v2 Vertex)
	Vertices() []Vertex
	Edges() []Edge
	Neighbours(v Vertex) []Vertex
}

var (
	_ Graph = &DirectedGraph{}
	_ Graph = &AdjacencyList{}
)

// Vertex represents a node in the graph. Users should create
// new Vertex values with NewVertex.
type Vertex int

// Edge represents an edge between two vertices.
// In a directed graph, the edge is from U to V.
type Edge struct{ U, V Vertex }

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

// DeleteVertex permanently removes vertex v.
func (g *AdjacencyList) DeleteVertex(v Vertex) {
	delete(g.edges, v)
	for vtx, vertices := range g.edges {
		for idx, candidate := range vertices {
			if candidate == v {
				g.edges[vtx] = append(vertices[:idx], vertices[idx+1:len(vertices)]...)
			}
		}
	}
}

// AddEdge adds an edge between v1 and v2.
func (g *AdjacencyList) AddEdge(v1, v2 Vertex) {
	if v2 < v1 {
		v1, v2 = v2, v1
	}
	edges := g.edges[v1]
	g.edges[v1] = append(edges, v2)
}

// RemoveEdge removes the edge between v2 and v2. It returns true if the edge
// was removed, and false otherwise.
func (g *AdjacencyList) RemoveEdge(v1, v2 Vertex) {
	if v2 < v1 {
		v1, v2 = v2, v1
	}
	vertices, ok := g.edges[v1]
	if !ok {
		return
	}
	var (
		idx int = -1
		vtx Vertex
	)
	for idx, vtx = range vertices {
		if vtx == v2 {
			break
		}
	}
	if idx >= 0 {
		// Remove the edge
		g.edges[v1] = append(g.edges[v1][:idx], g.edges[v1][idx+1:len(g.edges[v1])]...)
	}
}

// VertexSlice is a convenience for sorting vertices by ID.
type VertexSlice []Vertex

func (p VertexSlice) Len() int           { return len(p) }
func (p VertexSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p VertexSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p VertexSlice) Sort()              { sort.Sort(p) }

// EdgeSlice is a convenience for sorted edges by ID.
type EdgeSlice []Edge

func (p EdgeSlice) Len() int { return len(p) }
func (p EdgeSlice) Less(i, j int) bool {
	if p[i].U == p[j].U {
		return p[i].V < p[j].V
	} else {
		return p[i].U < p[j].U
	}
}
func (p EdgeSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p EdgeSlice) Sort()         { sort.Sort(p) }

// Vertices returns a slice of all vertices.
func (g *AdjacencyList) Vertices() []Vertex {
	vertices := make(VertexSlice, len(g.edges))
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

// Neighbours returns a slice of v's neighbours.
func (g *AdjacencyList) Neighbours(v Vertex) []Vertex {
	return g.edges[v]
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

// RemoveEdge removes the edge between v2 and v2. It returns true if the edge
// was removed, and false otherwise.
func (g *DirectedGraph) RemoveEdge(v1, v2 Vertex) {
	vertices, ok := g.edges[v1]
	if !ok {
		return
	}
	var (
		idx int = -1
		vtx Vertex
	)
	for idx, vtx = range vertices {
		if vtx == v2 {
			break
		}
	}
	if idx >= 0 {
		// Remove the edge
		g.edges[v1] = append(g.edges[v1][:idx], g.edges[v1][idx+1:len(g.edges[v1])]...)
	}
}

// NewVertex creates a new Vertex, adds it to the graph, and returns it.
func (g *DirectedGraph) NewVertex() Vertex {
	v := g.nextVertex
	g.addVertex(v)
	g.nextVertex++
	return v
}

// DeleteVertex permanently removes vertex v.
func (g *DirectedGraph) DeleteVertex(v Vertex) {
	delete(g.edges, v)
	for vtx, vertices := range g.edges {
		for idx, candidate := range vertices {
			if candidate == v {
				g.edges[vtx] = append(vertices[:idx], vertices[idx+1:len(vertices)]...)
			}
		}
	}
}

// Vertices returns a slice of the vertices that are in the graph.
func (g *DirectedGraph) Vertices() []Vertex {
	vertices := make([]Vertex, 0, len(g.edges))
	for k := range g.edges {
		vertices = append(vertices, k)
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

// Neighbours returns a slice of v's neighbours.
func (g *DirectedGraph) Neighbours(v Vertex) []Vertex {
	return g.edges[v]
}

func (g *DirectedGraph) Predecessors(v Vertex) (result []Vertex) {
	for vtx, vertices := range g.edges {
		for _, candidate := range vertices {
			if candidate == v {
				result = append(result, vtx)
			}
		}
	}
	return
}
