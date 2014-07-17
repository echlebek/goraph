package goraph

import (
	"reflect"
	"testing"
)

func TestEmptyTopoSort(t *testing.T) {
	g := NewDirectedAdjacencyList()
	result, err := TopoSort(g, false)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, []Vertex{}) {
		t.Fatalf("empty topo sort failed, something is seriously wrong.")
	}
}

func TestSimpleTopoSort(t *testing.T) {
	// TODO: construct a more complex toposort test
	graphEdges := map[Vertex][]Vertex{
		0: {1, 2},
		1: {3, 4},
		2: {5, 6},
	}
	g := newDirectedAdjacencyListFromMap(graphEdges)
	if g.nextVertex != 7 {
		t.Errorf("bad nextVertex: got %d, want %d", g.nextVertex, 7)
	}
	result, err := TopoSort(g, true)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Vertex{0, 2, 6, 5, 1, 4, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("bad TopoSort(): got %v, want %v", result, expected)
	}
}

func TestTopoSortCycle(t *testing.T) {
	// Test a graph with a cycle
	graphEdges := map[Vertex][]Vertex{
		0: {0, 1, 2},
		1: {3, 4},
		2: {5, 6},
	}
	g := newDirectedAdjacencyListFromMap(graphEdges)
	if g.nextVertex != 7 {
		t.Errorf("bad nextVertex: got %d, want %d", g.nextVertex, 7)
	}
	_, err := TopoSort(g, true)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestShortestPath(t *testing.T) {
	graphEdges := map[Vertex][]Vertex{
		0: {1},
		1: {2, 3},
		2: {3},
		3: {4, 5},
		4: {5},
		5: {},
	}
	g := &AdjacencyList{graphEdges, Vertex(6)}
	path := ShortestPath(g, Vertex(0), Vertex(5))
	expected := []Vertex{0, 1, 3, 5}
	if !reflect.DeepEqual(path, expected) {
		t.Errorf("bad path: got %v, want %v", path, expected)
	}
}
