package goraph

import (
	"fmt"
)

// Vertex represents a node in the graph. Users should create
// new Vertex values with NewVertex.
type Vertex uint64

// DirectedGraph provides a space-efficient directed graph.
type DirectedGraph struct {
	data           map[Vertex][]Vertex
	vertexSerialId uint64
}

// NewDirectedGraph creates and initializes a DirectedGraph.
func NewDirectedGraph() *DirectedGraph {
	g := &DirectedGraph{data: make(map[Vertex][]Vertex)}
	return g
}

// addVertex adds v to the graph in an idempotent fashion. The return value
// indicates whether or not the vertex was already in the graph; if false,
// the value was not in the graph before it was added.
func (g *DirectedGraph) addVertex(v Vertex) bool {
	_, ok := g.data[v]
	if !ok {
		g.data[v] = make([]Vertex, 0)
	}
	return ok
}

// AddEdge connects vertices v1 and v2 in the graph.
func (g *DirectedGraph) AddEdge(v1, v2 Vertex) {
	g.addVertex(v1)
	g.addVertex(v2)
	g.data[v1] = append(g.data[v1], v2)
}

// NewVertex creates a new Vertex, adds it to the graph, and returns it.
func (g *DirectedGraph) NewVertex() Vertex {
	v := Vertex(g.vertexSerialId)
	g.addVertex(v)
	g.vertexSerialId += 1
	return v
}

// TopoSort performs a topological sort on g.
// Based on pseudocode from http://en.wikipedia.org/wiki/Topological_sorting
// NB: Because go map keys are iterated in pseudorandom order,
// repeated invocations of TopoSort may differ.
func (g *DirectedGraph) TopoSort() []Vertex {
	// Shallow-copy the graph and iteratively remove edges from it later.
	newG := &DirectedGraph{make(map[Vertex][]Vertex, len(g.data)), g.vertexSerialId}
	for k, v := range g.data {
		newG.data[k] = v
	}
	g = newG

	result := make([]Vertex, 0, g.vertexSerialId)
	startVertices := g.findStartVertices()

	for len(startVertices) > 0 {
		v := startVertices[0]
		startVertices = startVertices[1:]
		result = append(result, v)
		for _, w := range g.data[v] {
			// w has no incoming edges
			if incoming := g.incomingEdges(w); len(incoming) == 1 {
				startVertices = append(startVertices, w)
			}
		}
		delete(g.data, v)
	}

	if len(g.data) != 0 {
		panic(fmt.Sprintf("topological sort failed: graph is not a DAG: %v", g.data))
	}

	return result
}

// IncomingEdges finds all the vertices that connect to v
func (g *DirectedGraph) incomingEdges(v Vertex) []Vertex {
	result := make([]Vertex, 0)
	for w, vlist := range g.data {
		for _, x := range vlist {
			if v == x {
				result = append(result, w)
				break
			}
		}
	}
	return result
}

// findStartVertices finds all the vertices with no incoming edges.
func (g *DirectedGraph) findStartVertices() []Vertex {
	result := make([]Vertex, 0)
	for candidate := range g.data {
		if incoming := g.incomingEdges(candidate); len(incoming) == 0 {
			result = append(result, candidate)
		}
	}

	return result
}
