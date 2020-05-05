package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")

	expected := NewEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})
	actual, err := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	source, _ := g.Node(n1.UID)
	target, _ := g.Node(n2.UID)

	_, ok := source.outEdges[actual.UID]
	assert.Equal(t, true, ok)

	_, ok = target.inEdges[actual.UID]
	assert.Equal(t, true, ok)
}

func TestAddEdge_missing_source(t *testing.T) {
	g := New()

	n2, _ := g.AddNode("node-2", "person")
	actual, err := g.AddEdge("edge-1234", "nissing", "knows", n2.UID, KV{Key: "since", Value: []byte("school")})

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}

func TestAddEdge_missing_target(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	actual, err := g.AddEdge("edge-1234", n1.UID, "knows", "missing", KV{Key: "since", Value: []byte("school")})

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}

func TestAddEdge_duplicate_uid(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})
	actual, err := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}

func TestHasEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})

	assert.Equal(t, true, g.HasEdge("edge-1234"))
}

func TestHasEdge_missing(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})

	assert.Equal(t, false, g.HasEdge("missing"))
}

func TestRemoveEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	n3, _ := g.AddNode("node-3", "person")

	edge1, _ := g.AddEdge("edge-1", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})
	edge2, _ := g.AddEdge("edge-2", n1.UID, "knows", n3.UID)

	err := g.RemoveEdge("edge-1")
	assert.Nil(t, err)
	assert.Equal(t, false, g.HasEdge("edge-1"))

	source, _ := g.Node(n1.UID)
	target, _ := g.Node(n2.UID)

	// (n1)->(n2)
	_, ok := source.outEdges[edge1.UID]
	assert.Equal(t, false, ok)

	// (n1)->(n3)
	_, ok = source.outEdges[edge2.UID]
	assert.Equal(t, true, ok)

	// (n1)<-(n2)
	_, ok = target.inEdges[edge1.UID]
	assert.Equal(t, false, ok)
}

func TestRemoveEdge_missing_edge(t *testing.T) {
	g := New()

	err := g.RemoveEdge("edge-1")
	assert.NotNil(t, err)
}

func TestEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	expected, _ := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})
	actual, err := g.Edge("edge-1234")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEdge_no_such_edge(t *testing.T) {
	g := New()

	actual, err := g.Edge("edge-missing")

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}

func TestEdges(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")

	e1, _ := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID)
	e2, _ := g.AddEdge("edge-2345", n1.UID, "knows", n2.UID)
	e3, _ := g.AddEdge("edge-3456", n1.UID, "knows", n2.UID)

	expected := []Edge{e1, e2, e3}
	actual := []Edge{}

	iter := g.Edges()
	for iter.Next() {
		edge := iter.Value().(Edge)
		actual = append(actual, edge)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestUpdateEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")

	old, _ := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})
	old.Properties["since"] = []byte("2020")

	updated, err := g.UpdateEdge(old)

	assert.Nil(t, err)
	assert.Equal(t, updated, old)

	source, _ := g.Node(n1.UID)
	target, _ := g.Node(n2.UID)

	_, ok := source.outEdges[updated.UID]
	assert.Equal(t, true, ok)

	_, ok = target.inEdges[updated.UID]
	assert.Equal(t, true, ok)
}

func TestUpdateEdge_missing_edge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")

	old, _ := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: []byte("school")})
	g.RemoveEdge(old.UID)

	old.Properties["since"] = []byte("2020")

	updated, err := g.UpdateEdge(old)
	assert.NotNil(t, err)
	assert.Equal(t, updated, old)
}

func TestEdgeCount(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")

	g.AddEdge("edge-1", n1.UID, "knows", n2.UID)
	g.AddEdge("edge-2", n1.UID, "knows", n2.UID)
	g.AddEdge("edge-3", n1.UID, "knows", n2.UID)

	assert.Equal(t, 3, g.EdgeCount())
}
