package graph

import (
	"fmt"

	"github.com/jenmud/draft/graph/iterator"
)

// HasEdge returns true if the graph has a edge with the provided uid.
func (g *Graph) HasEdge(uid string) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	_, ok := g.edges[uid]
	return ok
}

// UpdateEdge updates the graph edge with the new edge.
func (g *Graph) UpdateEdge(edge Edge) (Edge, error) {
	if !g.HasEdge(edge.UID) {
		return edge, fmt.Errorf("[UpdateEdge] Edge does not exists, can not update edge %s", edge)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	g.edges[edge.UID] = edge
	return edge, nil
}

// AddEdge adds a new edge to the graph.
func (g *Graph) AddEdge(uid, sourceUID, label, targetUID string, kv ...KV) (Edge, error) {
	if !g.HasNode(sourceUID) {
		return Edge{}, fmt.Errorf("[AddEdge] No such not with UID %s", sourceUID)
	}

	if !g.HasNode(targetUID) {
		return Edge{}, fmt.Errorf("[AddEdge] No such not with UID %s", targetUID)
	}

	source, err := g.Node(sourceUID)
	if err != nil {
		return Edge{}, fmt.Errorf("[AddEdge] %s", err)
	}

	target, err := g.Node(targetUID)
	if err != nil {
		return Edge{}, fmt.Errorf("[AddEdge] %s", err)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.edges[uid]; ok {
		return Edge{}, fmt.Errorf("[AddEdge] Edge UID %s already exists", uid)
	}

	edge := NewEdge(uid, sourceUID, label, targetUID, kv...)
	g.edges[edge.UID] = edge

	// (source)->(target)
	source.outEdges[edge.UID] = struct{}{}
	target.inEdges[edge.UID] = struct{}{}

	return edge, nil
}

// RemoveEdge removes the edge from the graph.
func (g *Graph) RemoveEdge(uid string) error {
	edge, err := g.Edge(uid)
	if err != nil {
		return fmt.Errorf("[RemoveEdge] %s", err)
	}

	// (source)->(target)
	source, err := g.Node(edge.SourceUID)
	if err != nil {
		// this is only here for safty, but we shoud not
		// get into a situation where this error is returned.
		return fmt.Errorf("[RemoveEdge] %s", err)
	}
	delete(source.outEdges, uid)

	target, err := g.Node(edge.TargetUID)
	if err != nil {
		// this is only here for safty, but we shoud not
		// get into a situation where this error is returned.
		return fmt.Errorf("[RemoveEdge] %s", err)
	}
	delete(target.inEdges, uid)

	g.lock.Lock()
	defer g.lock.Unlock()

	delete(g.edges, uid)
	return nil
}

// Edge returns the edge with the provided uid.
func (g *Graph) Edge(uid string) (Edge, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	edge, ok := g.edges[uid]
	if !ok {
		return Edge{}, fmt.Errorf("[GetEdge] No such edge with UID %s found", uid)
	}

	return edge, nil
}

// Edges returns a edge iterator.
func (g *Graph) Edges() Iterator {
	g.lock.RLock()
	defer g.lock.RUnlock()

	edges := make([]interface{}, len(g.edges))
	count := 0
	for _, edge := range g.edges {
		edges[count] = edge
		count++
	}

	return iterator.New(edges)
}

// EdgeCount returns the total number of edges in the graph.
func (g *Graph) EdgeCount() int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return len(g.edges)
}
