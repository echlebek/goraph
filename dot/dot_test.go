package dot

import (
	"bytes"
	"github.com/echlebek/goraph"
	"testing"
)

func TestDirectedDotWriter(t *testing.T) {
	graph := goraph.NewDirectedGraph()
	verts := make([]goraph.Vertex, 0, 10)
	for i := 0; i < 10; i++ {
		verts = append(verts, graph.AddVertex())
	}
	graph.AddEdge(verts[4], verts[3])
	graph.AddEdge(verts[4], verts[5])
	graph.AddEdge(verts[3], verts[1])
	graph.AddEdge(verts[1], verts[0])
	graph.AddEdge(verts[1], verts[2])
	graph.AddEdge(verts[5], verts[6])
	graph.AddEdge(verts[6], verts[9])
	graph.AddEdge(verts[6], verts[7])
	graph.AddEdge(verts[9], verts[8])

	buf := new(bytes.Buffer)

	dot := NewDot(graph)

	dot.VertexAttrs[verts[0]]["label"] = "Happy"
	dot.VertexAttrs[verts[1]]["label"] = "Sleepy"
	dot.VertexAttrs[verts[1]]["shape"] = "egg"

	dot.GraphAttrs["splines"] = false

	dot.EdgeAttrs[goraph.Edge{verts[5], verts[6]}]["arrowhead"] = "diamond"

	WriteDot(buf, dot)

	expected := `digraph {
	graph [ splines=false, ];
	0 [ label=Happy, ];
	1 [ label=Sleepy, shape=egg, ];
	1 -> 0;
	1 -> 2;
	3 -> 1;
	4 -> 3;
	4 -> 5;
	5 -> 6 [ arrowhead=diamond, ];
	6 -> 7;
	6 -> 9;
	9 -> 8;
}
`

	if buf.String() != expected {
		t.Errorf("bad dot output: got %s, want %s", buf.String(), expected)
	}

}
