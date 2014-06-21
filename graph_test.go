package goraph

import (
	"reflect"
	"testing"
)

func TestEmptyTopoSort(t *testing.T) {
	g := NewDirectedGraph()
	result := g.TopoSort()
	if !reflect.DeepEqual(result, []Vertex{}) {
		t.Fatalf("empty topo sort failed, something is seriously wrong.")
	}
}

// newDirectedGraphFromMap is a convenience function to create
// a DirectedGraph from a map.
func newDirectedGraphFromMap(edges map[Vertex][]Vertex) *DirectedGraph {
	g := NewDirectedGraph()
	var id uint64
	for k, v := range edges {
		if id < uint64(k) {
			id = uint64(k)
		}
		for _, vertex := range v {
			if id < uint64(vertex) {
				id = uint64(vertex)
			}
			g.AddEdge(k, vertex)
		}
	}
	g.vertexSerialId = id + 1
	return g
}

func TestSimpleTopoSort(t *testing.T) {
	// TODO: construct a more complex toposort test
	graphEdges := map[Vertex][]Vertex{
		0: {1, 2},
		1: {3, 4},
		2: {5, 6},
	}
	g := newDirectedGraphFromMap(graphEdges)
	if g.vertexSerialId != 7 {
		t.Errorf("bad vertexSerialId: got %d, want %d", g.vertexSerialId, 7)
	}
	result := g.TopoSort()
	expected := []Vertex{0, 1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("bad TopoSort(): got %v, want %v", result, expected)
	}
}

func TestAddVertices(t *testing.T) {
	g := NewDirectedGraph()

	v1 := g.NewVertex()
	v2 := g.NewVertex()
	v3 := g.NewVertex()
	v4 := g.NewVertex()
	v5 := g.NewVertex()
	v6 := g.NewVertex()

	g.AddEdge(v1, v2)
	g.AddEdge(v1, v3)
	g.AddEdge(v2, v4)
	g.AddEdge(v2, v5)
	g.AddEdge(v3, v6)

	expectedEdges := map[Vertex][]Vertex{
		0: {1, 2},
		1: {3, 4},
		2: {5},
		3: {},
		4: {},
		5: {},
	}
	if !reflect.DeepEqual(g.edges, expectedEdges) {
		t.Errorf("bad graph edges: got %v, want %v", g.edges, expectedEdges)
	}
}

func TestVertices(t *testing.T) {
	g := NewDirectedGraph()

	vertices := make([]Vertex, 100)
	for i := 0; i < 100; i++ {
		vertices[i] = g.NewVertex()
	}

	gverts := g.Vertices()
	if !reflect.DeepEqual(vertices, gverts) {
		t.Errorf("bad Vertices. got %v, want %v", gverts, vertices)
	}
}
