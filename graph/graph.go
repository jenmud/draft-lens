package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"sync"
	"time"
)

const (
	// Unlimited is used when returning a unlimited level subgraph.
	Unlimited = 0
)

// New returns a new empty graph.
func New() *Graph {
	return &Graph{
		startTime: time.Now().UTC(),
		nodes:     make(map[string]Node),
		edges:     make(map[string]Edge),
	}
}

// NewFromJSON takes a JSON formatted output and returns a new Graph.
func NewFromJSON(r io.Reader) (*Graph, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	graph := New()
	err = json.Unmarshal(data, graph)
	return graph, err
}

// Graph is a graph store.
type Graph struct {
	lock      sync.RWMutex
	startTime time.Time
	nodes     map[string]Node
	edges     map[string]Edge
}

// Stats returns some stats on the current graph instance.
func (g *Graph) Stats() Stat {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s := Stat{
		StartTime:     g.startTime,
		NumCPU:        runtime.NumCPU(),
		NumGoroutings: runtime.NumGoroutine(),
		NodeCount:     g.NodeCount(),
		EdgeCount:     g.EdgeCount(),
	}

	runtime.ReadMemStats(&s.MemStats)

	return s
}

// SubGraph takes a starting node UID and returns a new subgraph
// with `n` levels deep. Levels `0` or `Unlimited` will return the subgraph
// with no level limit.
func (g *Graph) SubGraph(uid string, levels int) (*Graph, error) {
	//TODO: add in the levels
	subg := New()

	node, err := g.Node(uid)
	if err != nil {
		return subg, fmt.Errorf("[SubGraph] %s", err)
	}

	var addClosure func(edge Edge) error

	addClosure = func(edge Edge) error {
		source, err := g.Node(edge.SourceUID)
		if err != nil {
			return fmt.Errorf("[SubGraph] %s", err)
		}

		target, err := g.Node(edge.TargetUID)
		if err != nil {
			return fmt.Errorf("[SubGraph] %s", err)
		}

		// Add the node and continue if they are already in the graph.
		subg.AddNode(source.UID, source.Label, convertPropertiesToKV(source.Properties)...)
		subg.AddNode(target.UID, target.Label, convertPropertiesToKV(target.Properties)...)

		if _, err := subg.AddEdge(edge.UID, source.UID, edge.Label, target.UID, convertPropertiesToKV(edge.Properties)...); err != nil {
			return fmt.Errorf("[SubGraph] %s", err)
		}

		return nil
	}

	for edgeUID := range node.inEdges {
		edge, err := g.Edge(edgeUID)
		if err != nil {
			return subg, fmt.Errorf("[SubGraph] %s", err)
		}

		if err := addClosure(edge); err != nil {
			return subg, err
		}
	}

	for edgeUID := range node.outEdges {
		edge, err := g.Edge(edgeUID)
		if err != nil {
			return subg, fmt.Errorf("[SubGraph] %s", err)
		}

		if err := addClosure(edge); err != nil {
			return subg, err
		}
	}

	return subg, nil
}

// See graph_node.go for all the node related methods
// See graph_edges.go for all the edge related methods

// MarshalJSON marchals the graph into a JSON format.
func (g *Graph) MarshalJSON() ([]byte, error) {
	type G struct {
		Nodes []Node `json:"nodes"`
		Edges []Edge `json:"edges"`
	}

	nodes := g.Nodes()
	edges := g.Edges()

	graph := G{
		Nodes: make([]Node, nodes.Size()),
		Edges: make([]Edge, edges.Size()),
	}

	ncount := 0
	for nodes.Next() {
		node := nodes.Value().(Node)
		graph.Nodes[ncount] = node
		ncount++
	}

	ecount := 0
	for edges.Next() {
		edge := edges.Value().(Edge)
		graph.Edges[ecount] = edge
		ecount++
	}

	return json.Marshal(graph)
}

// UnmarshalJSON unmarshals JSON data into the graph.
func (g *Graph) UnmarshalJSON(b []byte) error {
	type G struct {
		Nodes []Node `json:"nodes"`
		Edges []Edge `json:"edges"`
	}

	graph := G{}
	if err := json.Unmarshal(b, &graph); err != nil {
		return err
	}

	for _, node := range graph.Nodes {
		if _, err := g.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
			return err
		}
	}

	for _, edge := range graph.Edges {
		if _, err := g.AddEdge(edge.UID, edge.SourceUID, edge.Label, edge.TargetUID, convertPropertiesToKV(edge.Properties)...); err != nil {
			return err
		}
	}

	return nil
}
