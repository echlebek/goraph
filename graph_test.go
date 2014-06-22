package goraph

import (
	"reflect"
	"sort"
	"testing"
)

// newDirectedGraphFromMap is a convenience function to create
// a DirectedGraph from a map.
func newDirectedGraphFromMap(edges map[Vertex][]Vertex) *DirectedGraph {
	g := NewDirectedGraph()
	var id Vertex
	for k, v := range edges {
		if id < k {
			id = k
		}
		for _, vertex := range v {
			if id < vertex {
				id = vertex
			}
			g.AddEdge(k, vertex)
		}
	}
	g.nextVertex = id + 1
	return g
}
func TestAddVertices(t *testing.T) {
	g := NewDirectedGraph()

	v1 := g.NewVertex()
	v2 := g.NewVertex()
	v3 := g.NewVertex()
	v4 := g.NewVertex()
	v5 := g.NewVertex()
	v6 := g.NewVertex()

	expectedEdges := []Edge{{v1, v2}, {v1, v3}, {v2, v4}, {v2, v5}, {v3, v6}}

	for _, e := range expectedEdges {
		g.AddEdge(e.v1, e.v2)
	}

	edges := g.Edges()

	if !reflect.DeepEqual(edges, expectedEdges) {
		t.Errorf("bad graph edges: got %v, want %v", edges, expectedEdges)
	}
}

type byId []Vertex

func (v byId) Less(i, j int) bool { return v[i] < v[j] }
func (v byId) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byId) Len() int           { return len(v) }

func TestVertices(t *testing.T) {
	g := NewDirectedGraph()

	vertices := make([]Vertex, 100)
	for i := 0; i < 100; i++ {
		vertices[i] = g.NewVertex()
	}

	gverts := g.Vertices()
	sort.Sort(byId(vertices))
	sort.Sort(byId(gverts))
	if !reflect.DeepEqual(vertices, gverts) {
		t.Errorf("bad Vertices. got %v, want %v", gverts, vertices)
	}
}
