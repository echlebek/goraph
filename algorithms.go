package goraph

import "fmt"

// TopoSort performs a topological sort on g.
// Based on pseudocode from http://en.wikipedia.org/wiki/Topological_sorting
// NB: Because go map keys are iterated in pseudorandom order,
// repeated invocations of TopoSort may differ.
func TopoSort(g *DirectedGraph) ([]Vertex, error) {
	// Shallow-copy the graph and iteratively remove edges from it later.
	newG := &DirectedGraph{make(map[Vertex][]Vertex, len(g.edges)), g.nextVertex}
	for k, v := range g.edges {
		newG.edges[k] = v
	}
	g = newG

	result := make([]Vertex, 0, len(g.edges))
	startVertices := findStartVertices(g)

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
		return nil, fmt.Errorf("topological sort failed: graph is not a DAG: %v", g.edges)
	}

	return result, nil
}

// findStartVertices finds all the vertices with no incoming edges.
func findStartVertices(g *DirectedGraph) []Vertex {
	result := make([]Vertex, 0)
	for candidate := range g.edges {
		if incoming := g.incomingEdges(candidate); len(incoming) == 0 {
			result = append(result, candidate)
		}
	}

	return result
}
