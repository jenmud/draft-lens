package graph

import (
	"bytes"
	"fmt"

	"github.com/jenmud/draft/graph/iterator"
)

// HasNode returns true if the graph has a node with the provided uid.
func (g *Graph) HasNode(uid string) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	_, ok := g.nodes[uid]
	return ok
}

// AddNode adds a new node to the graph.
func (g *Graph) AddNode(uid, label string, kv ...KV) (Node, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.nodes[uid]; ok {
		return Node{}, fmt.Errorf("[AddNode] Node UID %s already exists", uid)
	}

	node := NewNode(uid, label, kv...)
	g.nodes[node.UID] = node
	return node, nil
}

// UpdateNode updates the graph node with the new node.
func (g *Graph) UpdateNode(node Node) (Node, error) {
	if !g.HasNode(node.UID) {
		return node, fmt.Errorf("[UpdateNode] Node does not exists, can not update node %s", node)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	g.nodes[node.UID] = node
	return node, nil
}

// RemoveNode removes the node from the graph.
func (g *Graph) RemoveNode(uid string) error {
	node, err := g.Node(uid)
	if err != nil {
		return fmt.Errorf("[RemoveNode] %s", err)
	}

	edgeCount := len(node.inEdges) + len(node.outEdges)
	if edgeCount > 0 {
		return fmt.Errorf("[RemoveNode] Can not remove node with edges attached (edge count: %d)", edgeCount)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	delete(g.nodes, uid)
	return nil
}

// Node returns the node with the provided uid.
func (g *Graph) Node(uid string) (Node, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	node, ok := g.nodes[uid]
	if !ok {
		return Node{}, fmt.Errorf("[GetNode] No such node with UID %s found", uid)
	}

	return node, nil
}

// Nodes returns a node iterator.
func (g *Graph) Nodes() Iterator {
	g.lock.RLock()
	defer g.lock.RUnlock()

	nodes := make([]interface{}, len(g.nodes))
	count := 0
	for _, node := range g.nodes {
		nodes[count] = node
		count++
	}

	return iterator.New(nodes)
}

// nodeMapper converts interfaces into Nodes
func nodeMapper(in <-chan interface{}, out chan<- Node) {
	for item := range in {
		out <- item.(Node)
	}
	close(out)
}

// labelReducer filters for nodes that have the given labels.
func labelReducer(labels []string, in <-chan Node, out chan<- Node) {
	for node := range in {
		if len(labels) == 0 {
			out <- node
			continue
		}

		for _, label := range labels {
			if node.Label == label {
				out <- node
				continue
			}
		}
	}

	close(out)
}

// propReducer filters for nodes that have the given properties.
func propReducer(props map[string][]byte, in <-chan Node, out chan<- Node) {
	for node := range in {
		if len(props) == 0 {
			out <- node
			continue
		}

		var allMatched = false
		for key, value := range props {
			nvalue, ok := node.Properties[key]
			if !ok {
				allMatched = false
				break
			}

			if !bytes.Equal(value, nvalue) {
				allMatched = false
				break
			}

			allMatched = true
		}

		if allMatched {
			out <- node
		}
	}

	close(out)
}

// NodesBy returns a node iterator with filtered nodes.
// If labels is an empty list, then any label will be used.
// If props is an empty map, no properties will be used for filtering.
func (g *Graph) NodesBy(labels []string, props map[string][]byte) Iterator {
	g.lock.RLock()
	defer g.lock.RUnlock()

	in := make(chan Node, len(g.nodes))
	labelFiltered := make(chan Node, len(g.nodes))
	final := make(chan Node)

	go nodeMapper(g.Nodes().Channel(), in)
	go labelReducer(labels, in, labelFiltered)
	go propReducer(props, labelFiltered, final)

	nodes := []interface{}{}
	for node := range final {
		nodes = append(nodes, node)
	}

	return iterator.New(nodes)
}

// NodeCount returns the total number of nodes in the graph.
func (g *Graph) NodeCount() int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return len(g.nodes)
}
