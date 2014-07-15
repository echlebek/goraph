package bench

import (
	gp "github.com/echlebek/goraph"
	"testing"
)

var (
	dense gp.Graph
	// TODO test sparse graph
)

func init() {
	initDense()
}

func initDense() {
	dense = gp.NewDirectedAdjacencyList()
	vertices := make([]gp.Vertex, 0, 1000)
	for i := 0; i < 1000; i++ {
		vertices = append(vertices, dense.AddVertex())
	}
	// Make (n^2/2)-1 connections
	for i := range vertices {
		for j := i; j < len(vertices); j++ {
			dense.AddEdge(vertices[i], vertices[j])
		}
	}
}

func BenchmarkTopoSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gp.TopoSort(dense, false)
	}
}

func BenchmarkDeterministicTopoSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gp.TopoSort(dense, true)
	}
}
