package dot

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/echlebek/goraph"
)

const (
	digraphSep = "->"
	graphSep   = "--"
)

type Dot struct {
	graph goraph.Graph

	// Name is the name of the graphviz graph.
	Name string

	// GraphAttrs are attributes that apply to the graph as a whole.
	GraphAttrs map[string]interface{}

	// EdgeAttrs are attributes that apply to specific edges in the graph.
	EdgeAttrs map[goraph.Edge]map[string]interface{}

	// VertexAttrs are attributes that apply to specific vertices in the graph.
	VertexAttrs map[goraph.Vertex]map[string]interface{}

	// EdgeGlobalAttrs are attributes that apply to every edge in the graph.
	EdgeGlobalAttrs map[string]interface{}

	// VertexGlobalAttrs are attributes that apply to every vertex in the graph.
	VertexGlobalAttrs map[string]interface{}
}

// NewDot creates a new Dot from g.
func NewDot(g goraph.Graph) Dot {
	result := Dot{
		g,
		"",
		make(map[string]interface{}),
		make(map[goraph.Edge]map[string]interface{}),
		make(map[goraph.Vertex]map[string]interface{}),
		make(map[string]interface{}),
		make(map[string]interface{}),
	}
	for _, v := range g.Vertices() {
		result.VertexAttrs[v] = make(map[string]interface{})
	}
	for _, e := range g.Edges() {
		result.EdgeAttrs[e] = make(map[string]interface{})
	}
	return result
}

func writeAttrs(w *bytes.Buffer, name string, tabs int, attrs map[string]interface{}) {
	ws := ""
	for i := 0; i < tabs; i++ {
		ws += "\t"
	}
	fmt.Fprint(w, ws)
	fmt.Fprintf(w, "%s [ ", name)
	keys := make(sort.StringSlice, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	list := []string{}
	for _, k := range keys {
		list = append(list, fmt.Sprintf("%s=%v", k, attrs[k]))
	}
	fmt.Fprintln(w, strings.Join(list, ", "), "];")
}

// WriteDot writes dot to w. It returns the number of bytes written and
// any error that occurred.
func WriteDot(w io.Writer, dot Dot) (int64, error) {
	var graphT, nodeSep string
	g := dot.graph
	_, ok := g.(*goraph.DirectedAdjacencyList)
	if ok {
		graphT = "digraph"
		nodeSep = digraphSep
	} else {
		graphT = "graph"
		nodeSep = graphSep
	}
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "%s {\n", graphT)

	// Write the global attrs if they exist.
	if len(dot.GraphAttrs) > 0 {
		writeAttrs(buf, "graph", 1, dot.GraphAttrs)
	}
	if len(dot.EdgeGlobalAttrs) > 0 {
		writeAttrs(buf, "edge", 1, dot.EdgeGlobalAttrs)
	}
	if len(dot.VertexGlobalAttrs) > 0 {
		writeAttrs(buf, "node", 1, dot.VertexGlobalAttrs)
	}

	// Write all the vertices
	vertices := goraph.VertexSlice(g.Vertices())
	vertices.Sort()
	for _, v := range vertices {
		attrs, ok := dot.VertexAttrs[v]
		if !ok || len(attrs) == 0 {
			continue
		}
		writeAttrs(buf, fmt.Sprintf("%d", v), 1, attrs)
	}

	// Write all the edges
	edges := goraph.EdgeSlice(g.Edges())
	edges.Sort()
	for _, e := range edges {
		fmt.Fprintf(buf, "\t%d %s %d", e.U, nodeSep, e.V)
		if attrs, ok := dot.EdgeAttrs[e]; ok && len(attrs) > 0 {
			writeAttrs(buf, "", 0, attrs)
		} else {
			fmt.Fprint(buf, ";\n")
		}
	}
	fmt.Fprint(buf, "}\n")
	return io.Copy(w, buf)
}
