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

func TestAdjacencyList(t *testing.T) {
	g := NewAdjacencyList()
	testGraph(t, g)
}

func TestDirectedGraph(t *testing.T) {
	g := NewDirectedGraph()
	testGraph(t, g)
}

func testGraph(t *testing.T, g Graph) {
	vertices := make([]Vertex, 100)
	for i := 0; i < 100; i++ {
		vertices[i] = g.NewVertex()
	}

	v1 := vertices[0]
	v2 := vertices[1]
	v3 := vertices[2]
	v4 := vertices[3]
	v5 := vertices[4]
	v6 := vertices[5]

	gverts := g.Vertices()
	sort.Sort(verticesById(vertices))
	sort.Sort(verticesById(gverts))
	if !reflect.DeepEqual(vertices, gverts) {
		t.Errorf("bad Vertices. got %v, want %v", gverts, vertices)
	}

	expectedEdges := []Edge{{v1, v2}, {v1, v3}, {v2, v4}, {v2, v5}, {v3, v6}}
	sort.Sort(edgesByV(expectedEdges))

	for _, e := range expectedEdges {
		g.AddEdge(e.v1, e.v2)
	}

	edges := g.Edges()
	sort.Sort(edgesByV(edges))

	if !reflect.DeepEqual(edges, expectedEdges) {
		t.Errorf("bad graph edges: got %v, want %v", edges, expectedEdges)
	}
}

type verticesById []Vertex

func (v verticesById) Less(i, j int) bool { return v[i] < v[j] }
func (v verticesById) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v verticesById) Len() int           { return len(v) }

type edgesByV []Edge

func (e edgesByV) Less(i, j int) bool {
	ei := e[i]
	ej := e[j]
	return ei.v1 < ej.v1 && ei.v2 < ej.v2
}

func (e edgesByV) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e edgesByV) Len() int      { return len(e) }
