package goraph

import (
	"fmt"
)

// Vertex represents a node in the graph. Users should create
// new Vertex values with NewVertex.
type Vertex uint64

// DirectedGraph provides a space-efficient directed graph.
type DirectedGraph struct {
	edges          map[Vertex][]Vertex
	vertexSerialId uint64
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
	v := Vertex(g.vertexSerialId)
	g.addVertex(v)
	g.vertexSerialId += 1
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

// TopoSort performs a topological sort on g.
// Based on pseudocode from http://en.wikipedia.org/wiki/Topological_sorting
// NB: Because go map keys are iterated in pseudorandom order,
// repeated invocations of TopoSort may differ.
func (g *DirectedGraph) TopoSort() []Vertex {
	// Shallow-copy the graph and iteratively remove edges from it later.
	newG := &DirectedGraph{make(map[Vertex][]Vertex, len(g.edges)), g.vertexSerialId}
	for k, v := range g.edges {
		newG.edges[k] = v
	}
	g = newG

	result := make([]Vertex, 0, len(g.edges))
	startVertices := g.findStartVertices()

	for len(startVertices) > 0 {
		v := startVertices[0]
		startVertices = startVertices[1:]
		result = append(result, v)
		for _, w := range g.edges[v] {
			// w has no incoming edges except for v's
			if incoming := g.countIncomingEdges(w); incoming == 1 {
				startVertices = append(startVertices, w)
			}
		}
		delete(g.edges, v)
	}

	if len(g.edges) != 0 {
		panic(fmt.Sprintf("topological sort failed: graph is not a DAG: %v", g.edges))
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

// findStartVertices finds all the vertices with no incoming edges.
func (g *DirectedGraph) findStartVertices() []Vertex {
	result := make([]Vertex, 0)
	for candidate := range g.edges {
		if incoming := g.incomingEdges(candidate); len(incoming) == 0 {
			result = append(result, candidate)
		}
	}

	return result
}
