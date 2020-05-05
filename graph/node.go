package graph

// NewNode returns a new node instance.
func NewNode(uid, label string, kv ...KV) Node {
	return Node{
		UID:        uid,
		Label:      label,
		Properties: NewProperties(kv...),
		inEdges:    make(map[string]struct{}),
		outEdges:   make(map[string]struct{}),
	}
}

// Node is a node in the graph.
type Node struct {
	UID        string            `json:"uid"`
	Label      string            `json:"label"`
	Properties map[string][]byte `json:"properties"`
	inEdges    map[string]struct{}
	outEdges   map[string]struct{}
}

// InEdges returns all the inbound edges.
// ()-->(n)
func (n Node) InEdges() []string {
	edges := make([]string, len(n.inEdges))
	count := 0
	for edge := range n.inEdges {
		edges[count] = edge
		count++
	}
	return edges
}

// OutEdges returns all the outbound edges.
// (n)-->()
func (n Node) OutEdges() []string {
	edges := make([]string, len(n.outEdges))
	count := 0
	for edge := range n.outEdges {
		edges[count] = edge
		count++
	}
	return edges
}

// Edges returns all the in and outbound edges.
// ()-->(n)-->()
func (n Node) Edges() []string {
	in := n.InEdges()
	out := n.OutEdges()
	edges := []string{}
	edges = append(edges, in...)
	edges = append(edges, out...)
	return edges
}
