package goraph

import (
	"fmt"
)

// TopoSort returns a slice of the vertices in g in topologically
// sorted order. (https://en.wikipedia.org/wiki/Topological_sorting)
//
// If deterministic is true, then the vertices are sorted by id before
// the algorithm begins, guaranteeing a repeatable result.
func TopoSort(g Graph, deterministic bool) (result []Vertex, err error) {
	verts := VertexSlice(g.Vertices())
	if deterministic {
		verts.Sort()
	}

	rlen := len(verts)
	result = make([]Vertex, 0, rlen)
	marked := make(map[Vertex]bool, rlen) // visited vertices

	var visit func(Vertex)
	visit = func(vtx Vertex) {
		if permanent, ok := marked[vtx]; ok {
			if !permanent {
				err = fmt.Errorf("cannot perform toposort: not a DAG")
			} else {
				return
			}
		}
		// Mark vtx temporarily
		marked[vtx] = false
		for _, v := range g.Neighbours(vtx) {
			visit(v)
		}
		// mark vtx permanently
		marked[vtx] = true
		result = append(result, vtx)
	}

	for _, v := range verts {
		visit(v)
	}

	// The algorithm asks us to prepend to the result, but since we're using a
	// slice here, just reverse it after appending items.
	for i := 0; i < len(result)/2; i++ {
		result[i], result[rlen-i-1] = result[rlen-i-1], result[i]
	}

	return
}

// ShortestPath returns the shortest path between source and target using
// Dijkstra's algorithm.
func ShortestPath(g Graph, source, target Vertex) []Vertex {
	pq := make(priorityQueue, 0)

	// Vertex ancestry mapping
	prev := make(map[Vertex]Vertex)

	// Map of visited nodes to their distance from the source
	visited := make(map[Vertex]int)
	prev[source] = source

	s := &qItem{vertex: source, priority: 0}
	pq.Push(s)

	for pq.Len() > 0 {
		v := pq.Pop().(*qItem)
		if v.vertex == target {
			return walkAncestry(prev, v.vertex)
		}
		for _, w := range g.Neighbours(v.vertex) {
			newPrio := v.priority + 1
			if oldPrio, ok := visited[w]; ok {
				if newPrio > oldPrio {
					continue
				}
			}
			visited[w] = newPrio
			prev[w] = v.vertex
			q := &qItem{vertex: w, priority: v.priority + 1}
			pq.Push(q)
		}
	}

	return nil
}

func walkAncestry(ancestors map[Vertex]Vertex, v Vertex) []Vertex {
	result := []Vertex{v}
	if ancestors[v] == v {
		return result
	}
	result = append(walkAncestry(ancestors, ancestors[v]), result...)
	return result
}
