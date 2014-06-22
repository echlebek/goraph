// A basic priority queue, nearly identical to the example from the
// container/heap docs.

package goraph

import (
	"container/heap"
)

type qItem struct {
	vertex   Vertex
	priority int
	index    int
}

type priorityQueue []*qItem

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*qItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and vertex of an qItem in the queue.
func (pq *priorityQueue) updateQItem(item *qItem, v Vertex, priority int) {
	item.vertex = v
	item.priority = priority
	heap.Fix(pq, item.index)
}
