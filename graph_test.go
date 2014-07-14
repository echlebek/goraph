package goraph

import (
	"reflect"
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

func TestAdjacencyList(t *testing.T) {
	g := NewAdjacencyList()
	testGraph(t, g)
}

func TestDirectedGraph(t *testing.T) {
	g := NewDirectedGraph()
	testGraph(t, g)
}

func testGraph(t *testing.T, g Graph) {
	vertices := make(VertexSlice, 100)
	for i := 0; i < 100; i++ {
		vertices[i] = g.NewVertex()
	}

	v1 := vertices[0]
	v2 := vertices[1]
	v3 := vertices[2]
	v4 := vertices[3]
	v5 := vertices[4]
	v6 := vertices[5]

	gverts := VertexSlice(g.Vertices())
	vertices.Sort()
	gverts.Sort()
	if !reflect.DeepEqual(vertices, gverts) {
		t.Errorf("bad Vertices. got %v, want %v", gverts, vertices)
	}

	expectedEdges := EdgeSlice{{v1, v2}, {v1, v3}, {v2, v4}, {v2, v5}, {v3, v6}}
	expectedEdges.Sort()

	for _, e := range expectedEdges {
		g.AddEdge(e.U, e.V)
	}

	edges := EdgeSlice(g.Edges())
	edges.Sort()

	if !reflect.DeepEqual(edges, expectedEdges) {
		t.Errorf("bad graph edges: got %v, want %v", edges, expectedEdges)
	}

	g.RemoveEdge(v3, v6)

	edges = EdgeSlice(g.Edges())
	edges.Sort()
	expectedEdges = EdgeSlice{{v1, v2}, {v1, v3}, {v2, v4}, {v2, v5}}

	if !reflect.DeepEqual(edges, expectedEdges) {
		t.Errorf("bad graph edges: got %v, want %v", edges, expectedEdges)
	}

	g.DeleteVertex(v2)
	edges = EdgeSlice(g.Edges())
	edges.Sort()
	expectedEdges = EdgeSlice{{v1, v3}}

	if !reflect.DeepEqual(edges, expectedEdges) {
		t.Errorf("bad graph edges: got %v, want %v", edges, expectedEdges)
	}
}

func TestPredecessors(t *testing.T) {
	g := NewDirectedGraph()
	vertices := make(VertexSlice, 10)
	for i := range vertices {
		vertices[i] = g.NewVertex()
		if i > 0 {
			g.AddEdge(vertices[i], vertices[0])
		}
	}
	preds := VertexSlice(g.Predecessors(vertices[0]))
	preds.Sort()
	expected := VertexSlice{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(preds, expected) {
		t.Errorf("bad predecessors: got %v, want %v", preds, expected)
	}
}
