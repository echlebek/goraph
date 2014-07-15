package goraph

import (
	"sort"
)

// The Graph interface is implemented by all graph types.
type Graph interface {
	AddVertex() Vertex
	RemoveVertex(v Vertex)
	AddEdge(v1, v2 Vertex)
	RemoveEdge(v1, v2 Vertex)
	Vertices() []Vertex
	Edges() []Edge
	Neighbours(v Vertex) []Vertex
	IsDirected() bool
}

var (
	_ Graph = &AdjacencyList{}
)

// Vertex represents a node in the graph. Users should create
// new Vertex values with AddVertex.
type Vertex int

// Edge represents an edge between two vertices.
// In a directed graph, the edge is from U to V.
type Edge struct{ U, V Vertex }

// AdjacencyList implements an undirected graph using an adjacency list.
type AdjacencyList struct {
	edges      map[Vertex][]Vertex
	nextVertex Vertex
	directed   bool
}

// NewAdjacencyList creates an empty graph.
func NewAdjacencyList() *AdjacencyList {
	return &AdjacencyList{edges: make(map[Vertex][]Vertex)}
}

// NewDirectedAdjacencyList creates an empty digraph.
func NewDirectedAdjacencyList() *AdjacencyList {
	return &AdjacencyList{edges: make(map[Vertex][]Vertex), directed: true}
}

// IsDirected returns true if the graph is a directed graph.
func (g *AdjacencyList) IsDirected() bool {
	return g.directed
}

// AddVertex adds a new vertex.
func (g *AdjacencyList) AddVertex() Vertex {
	v := g.nextVertex
	if _, ok := g.edges[v]; !ok {
		g.edges[v] = make([]Vertex, 0)
	}
	g.nextVertex++
	return v
}

// RemoveVertex permanently removes vertex v.
func (g *AdjacencyList) RemoveVertex(v Vertex) {
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
	if v2 < v1 && !g.directed {
		v1, v2 = v2, v1
	}
	edges := g.edges[v1]
	g.edges[v1] = append(edges, v2)
}

// RemoveEdge removes the edge between v2 and v2.
func (g *AdjacencyList) RemoveEdge(v1, v2 Vertex) {
	if v2 < v1 && !g.directed {
		v1, v2 = v2, v1
	}
	vertices := g.edges[v1]
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
		g.edges[v1] = append(vertices[:idx], vertices[idx+1:len(vertices)]...)
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

// Predecessors returns the vertices that consider v a successor.
func (g *AdjacencyList) Predecessors(v Vertex) (result []Vertex) {
	if !g.directed {
		return
	}
	for vtx, vertices := range g.edges {
		for _, candidate := range vertices {
			if candidate == v {
				result = append(result, vtx)
			}
		}
	}
	return
}
