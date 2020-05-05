package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNode(t *testing.T) {
	g := New()
	expected := NewNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	actual, err := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddNode_Duplicate(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	actual, err := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	assert.NotNil(t, err)
	assert.Equal(t, Node{}, actual)
}

func TestRemoveNode(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	err := g.RemoveNode("abcd-1234")
	assert.Nil(t, err)
	assert.Equal(t, false, g.HasNode("abcd-1234"))
}

func TestRemoveNode_does_not_exist(t *testing.T) {
	g := New()
	err := g.RemoveNode("node-1")
	assert.NotNil(t, err)
}

func TestRemoveNode_with_edges(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1", n1.UID, "knows", n2.UID)

	err := g.RemoveNode("node-1")
	assert.NotNil(t, err)
	assert.Equal(t, true, g.HasNode("node-1"))
	assert.Equal(t, true, g.HasNode("node-2"))
	assert.Equal(t, true, g.HasEdge("edge-1"))
}

func TestRemoveNode_after_edge_removal(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1", n1.UID, "knows", n2.UID)

	err := g.RemoveNode("node-1")
	assert.NotNil(t, err)

	g.RemoveEdge("edge-1")
	err = g.RemoveNode("node-1")
	assert.Nil(t, err)

	assert.Equal(t, false, g.HasNode("node-1"))
	assert.Equal(t, true, g.HasNode("node-2"))
	assert.Equal(t, false, g.HasEdge("edge-1"))
}

func TestHasNode(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	assert.Equal(t, true, g.HasNode("abcd-1234"))
}

func TestHasNode_not_found(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	assert.Equal(t, false, g.HasNode("missing"))
}

func TestNode(t *testing.T) {
	g := New()
	expected, _ := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	actual, err := g.Node("abcd-1234")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestNode_not_found(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	actual, err := g.Node("abcd-1234-missing")
	assert.NotNil(t, err)
	assert.Equal(t, Node{}, actual)
}

func TestlabelReducer(t *testing.T) {
	nodes := make(chan Node, 3)

	n1 := NewNode("node-1", "person")
	n2 := NewNode("node-2", "person")
	n3 := NewNode("node-3", "animal")

	nodes <- n1
	nodes <- n2
	nodes <- n3
	close(nodes)

	out := make(chan Node, 3)

	expected := []Node{n1, n2}
	actual := []Node{}

	labelReducer([]string{"person"}, nodes, out)

	for node := range out {
		actual = append(actual, node)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestlabelReducer__no_labels(t *testing.T) {
	nodes := make(chan Node, 3)

	n1 := NewNode("node-1", "person")
	n2 := NewNode("node-2", "person")
	n3 := NewNode("node-3", "animal")

	nodes <- n1
	nodes <- n2
	nodes <- n3
	close(nodes)

	out := make(chan Node, 3)

	expected := []Node{n1, n2, n3}
	actual := []Node{}

	labelReducer([]string{}, nodes, out)

	for node := range out {
		actual = append(actual, node)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestNodeMapper(t *testing.T) {
	nodes := make(chan interface{}, 3)

	n1 := NewNode("node-1", "person")
	n2 := NewNode("node-2", "person")
	n3 := NewNode("node-3", "animal")

	nodes <- n1
	nodes <- n2
	nodes <- n3
	close(nodes)

	out := make(chan Node, 3)

	expected := []Node{n1, n2, n3}
	actual := []Node{}

	nodeMapper(nodes, out)

	for node := range out {
		actual = append(actual, node)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestPropReducer(t *testing.T) {
	nodes := make(chan Node, 3)

	n1 := NewNode("node-1", "person", KV{Key: "name", Value: []byte("Foo")})
	n2 := NewNode("node-1", "person", KV{Key: "name", Value: []byte("Bar")})
	n3 := NewNode("node-1", "person", KV{Key: "age", Value: []byte("21")}, KV{Key: "name", Value: []byte("Foo")})

	nodes <- n1
	nodes <- n2
	nodes <- n3
	close(nodes)

	out := make(chan Node, 3)

	expected := []Node{n1, n3}
	actual := []Node{}

	propReducer(map[string][]byte{"name": []byte("Foo")}, nodes, out)

	for node := range out {
		actual = append(actual, node)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestPropReducer__multiple_props(t *testing.T) {
	nodes := make(chan Node, 3)

	n1 := NewNode("node-1", "person", KV{Key: "name", Value: []byte("Foo")})
	n2 := NewNode("node-1", "person", KV{Key: "name", Value: []byte("Bar")})
	n3 := NewNode("node-1", "person", KV{Key: "age", Value: []byte("21")}, KV{Key: "name", Value: []byte("Foo")})

	nodes <- n1
	nodes <- n2
	nodes <- n3
	close(nodes)

	out := make(chan Node, 3)

	expected := []Node{n3}
	actual := []Node{}

	propReducer(map[string][]byte{"name": []byte("Foo"), "age": []byte("21")}, nodes, out)

	for node := range out {
		actual = append(actual, node)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestPropReducer__empty_props(t *testing.T) {
	nodes := make(chan Node, 3)

	n1 := NewNode("node-1", "person", KV{Key: "name", Value: []byte("Foo")})
	n2 := NewNode("node-1", "person", KV{Key: "name", Value: []byte("Bar")})
	n3 := NewNode("node-1", "person", KV{Key: "age", Value: []byte("21")}, KV{Key: "name", Value: []byte("Foo")})

	nodes <- n1
	nodes <- n2
	nodes <- n3
	close(nodes)

	out := make(chan Node, 3)

	expected := []Node{n1, n2, n3}
	actual := []Node{}

	propReducer(map[string][]byte{}, nodes, out)

	for node := range out {
		actual = append(actual, node)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestNodesBy__label_filtered(t *testing.T) {
	g := New()
	g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "pet")
	n3, _ := g.AddNode("node-3", "bike")
	g.AddNode("node-4", "person")

	expected := []Node{n2, n3}
	actual := []Node{}

	iter := g.NodesBy([]string{"pet", "bike"}, map[string][]byte{})
	for iter.Next() {
		actual = append(actual, iter.Value().(Node))
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestNodesBy__prop_filtered(t *testing.T) {
	g := New()
	g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "pet", KV{Key: "name", Value: []byte("socks")})
	g.AddNode("node-3", "bike")
	g.AddNode("node-4", "person")

	expected := []Node{n2}
	actual := []Node{}

	iter := g.NodesBy([]string{"pet", "bike"}, map[string][]byte{"name": []byte("socks")})
	for iter.Next() {
		actual = append(actual, iter.Value().(Node))
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestNodesBy__empty_labels_prop_filtered(t *testing.T) {
	g := New()
	g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "pet", KV{Key: "name", Value: []byte("socks")}, KV{Key: "enabled", Value: []byte("true")})
	n3, _ := g.AddNode("node-3", "bike", KV{Key: "enabled", Value: []byte("true")})
	g.AddNode("node-4", "person")

	expected := []Node{n2, n3}
	actual := []Node{}

	iter := g.NodesBy([]string{}, map[string][]byte{"enabled": []byte("true")})
	for iter.Next() {
		actual = append(actual, iter.Value().(Node))
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestNodesBy__emtpy_lables_empty_props(t *testing.T) {
	g := New()
	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "pet", KV{Key: "name", Value: []byte("socks")})
	n3, _ := g.AddNode("node-3", "bike")
	n4, _ := g.AddNode("node-4", "person")

	expected := []Node{n1, n2, n3, n4}
	actual := []Node{}

	iter := g.NodesBy([]string{}, map[string][]byte{})
	for iter.Next() {
		actual = append(actual, iter.Value().(Node))
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestNodes(t *testing.T) {
	g := New()
	expected1, _ := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	expected2, _ := g.AddNode("abcd-4321", "person", KV{Key: "name", Value: []byte("bar")})

	expected := []Node{expected1, expected2}
	actual := []Node{}

	iter := g.Nodes()
	for iter.Next() {
		actual = append(actual, iter.Value().(Node))
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestUpdateNode(t *testing.T) {
	g := New()

	old, err := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: []byte("foo")})
	old.Properties["name"] = []byte("bar")

	updated, err := g.UpdateNode(old)
	node, _ := g.Node(old.UID)

	assert.Nil(t, err)
	assert.Equal(t, updated, node)
}

func TestUpdateNode_missing_node(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	g.RemoveNode(n1.UID)
	n1.Properties["surname"] = []byte("Blah")

	updated, err := g.UpdateNode(n1)
	assert.NotNil(t, err)
	assert.Equal(t, updated, n1)
}

func TestNodeCount(t *testing.T) {
	g := New()

	g.AddNode("node-1", "person")
	g.AddNode("node-2", "person")
	g.AddNode("node-3", "person")

	assert.Equal(t, 3, g.NodeCount())
}

func TestNode_InEdges(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	n3, _ := g.AddNode("node-3", "person")

	g.AddEdge("edge-knows", n1.UID, "knows", n2.UID)
	g.AddEdge("edge-likes", n1.UID, "likes", n2.UID)
	g.AddEdge("edge-dislikes", n3.UID, "dislikes", n1.UID)

	expected := []string{"edge-dislikes"}
	assert.ElementsMatch(t, expected, n1.InEdges())
}

func TestNode_OutEdges(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	n3, _ := g.AddNode("node-3", "person")

	g.AddEdge("edge-knows", n1.UID, "knows", n2.UID)
	g.AddEdge("edge-likes", n1.UID, "likes", n2.UID)
	g.AddEdge("edge-dislikes", n3.UID, "dislikes", n1.UID)

	expected := []string{"edge-knows", "edge-likes"}
	assert.ElementsMatch(t, expected, n1.OutEdges())
}

func TestNode_Edges(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	n3, _ := g.AddNode("node-3", "person")

	g.AddEdge("edge-knows", n1.UID, "knows", n2.UID)
	g.AddEdge("edge-likes", n1.UID, "likes", n2.UID)
	g.AddEdge("edge-dislikes", n3.UID, "dislikes", n1.UID)

	expected := []string{"edge-knows", "edge-likes", "edge-dislikes"}
	assert.ElementsMatch(t, expected, n1.Edges())
}

func TestNode_Edges__remove_edge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	n3, _ := g.AddNode("node-3", "person")

	g.AddEdge("edge-knows", n1.UID, "knows", n2.UID)
	g.AddEdge("edge-likes", n1.UID, "likes", n2.UID)
	g.AddEdge("edge-dislikes", n3.UID, "dislikes", n1.UID)

	g.RemoveEdge("edge-dislikes")

	expected := []string{"edge-knows", "edge-likes"}
	assert.ElementsMatch(t, expected, n1.Edges())
}
