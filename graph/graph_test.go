package graph

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func readTestData(t *testing.T, name string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func TestSubGraph_one_level(t *testing.T) {
	reader := bytes.NewReader(readTestData(t, "simple-graph.json"))
	g, err := NewFromJSON(reader)
	if err != nil {
		t.Fatal(err)
	}

	subg, err := g.SubGraph("node-dog", 1)
	assert.Nil(t, err)

	foo, err := subg.Node("node-foo")
	assert.Equal(
		t,
		Node{
			UID:        "node-foo",
			Label:      "person",
			Properties: map[string][]byte{"name": []byte("foo")},
			inEdges:    map[string]struct{}{},
			outEdges: map[string]struct{}{
				"edge-owns": {},
			},
		},
		foo,
	)

	bar, err := subg.Node("node-bar")
	assert.Equal(
		t,
		Node{
			UID:        "node-bar",
			Label:      "person",
			Properties: map[string][]byte{"name": []byte("bar")},
			inEdges:    map[string]struct{}{},
			outEdges: map[string]struct{}{
				"edge-dislike": {},
			},
		},
		bar,
	)

	dog, err := subg.Node("node-dog")
	assert.Equal(
		t,
		Node{
			UID:        "node-dog",
			Label:      "animal",
			Properties: map[string][]byte{"name": []byte("socks")},
			inEdges: map[string]struct{}{
				"edge-owns":    {},
				"edge-dislike": {},
			},
			outEdges: map[string]struct{}{},
		},
		dog,
	)

	e1, err := subg.Edge("edge-dislike")
	assert.Nil(t, err)
	assert.Equal(t, Edge{UID: "edge-dislike", SourceUID: "node-bar", Label: "dislikes", TargetUID: "node-dog", Properties: map[string][]byte{}}, e1)

	e2, err := subg.Edge("edge-owns")
	assert.Nil(t, err)
	assert.Equal(t, Edge{UID: "edge-owns", SourceUID: "node-foo", Label: "owns", TargetUID: "node-dog", Properties: map[string][]byte{}}, e2)
}

func TestMarshalJSON(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-foo", "person", KV{Key: "name", Value: []byte("foo")})
	n2, _ := g.AddNode("node-bar", "person", KV{Key: "name", Value: []byte("bar")})
	n3, _ := g.AddNode("node-dog", "animal", KV{Key: "name", Value: []byte("socks")})

	g.AddEdge("edge-knows", n1.UID, "knows", n2.UID, KV{Key: "name", Value: []byte("2020")})
	g.AddEdge("edge-owns", n1.UID, "owns", n3.UID)
	g.AddEdge("edge-like", n1.UID, "likes", n2.UID)
	g.AddEdge("edge-dislike", n2.UID, "dislikes", n3.UID)

	dump, err := json.Marshal(g)
	assert.Nil(t, err)

	actual := New()

	// Hack the start times as they are not testable
	now := time.Now()
	g.startTime = now
	actual.startTime = now

	err = json.Unmarshal(dump, &actual)
	assert.Nil(t, err)

	assert.Equal(t, g, actual)
}

func TestUnmarshalJSON(t *testing.T) {
	dump := readTestData(t, "simple-graph.json")
	g := New()

	err := json.Unmarshal(dump, &g)
	assert.Nil(t, err)

	// Check node-foo was imported properly
	foo, err := g.Node("node-foo")
	assert.Nil(t, err)
	assert.Equal(
		t,
		Node{
			UID:        "node-foo",
			Label:      "person",
			Properties: map[string][]byte{"name": []byte("foo")},
			inEdges:    map[string]struct{}{},
			outEdges: map[string]struct{}{
				"edge-like":  {},
				"edge-knows": {},
				"edge-owns":  {},
			},
		},
		foo,
	)

	bar, err := g.Node("node-bar")
	assert.Nil(t, err)
	assert.Equal(
		t,
		Node{
			UID:        "node-bar",
			Label:      "person",
			Properties: map[string][]byte{"name": []byte("bar")},
			inEdges: map[string]struct{}{
				"edge-like":  {},
				"edge-knows": {},
			},
			outEdges: map[string]struct{}{
				"edge-dislike": {},
			},
		},
		bar,
	)

	dog, err := g.Node("node-dog")
	assert.Nil(t, err)
	assert.Equal(
		t,
		Node{
			UID:        "node-dog",
			Label:      "animal",
			Properties: map[string][]byte{"name": []byte("socks")},
			inEdges: map[string]struct{}{
				"edge-dislike": {},
				"edge-owns":    {},
			},
			outEdges: map[string]struct{}{},
		},
		dog,
	)

	e1, err := g.Edge("edge-like")
	assert.Nil(t, err)
	assert.Equal(t, Edge{UID: "edge-like", SourceUID: "node-foo", Label: "likes", TargetUID: "node-bar", Properties: map[string][]byte{}}, e1)

	e2, err := g.Edge("edge-dislike")
	assert.Nil(t, err)
	assert.Equal(t, Edge{UID: "edge-dislike", SourceUID: "node-bar", Label: "dislikes", TargetUID: "node-dog", Properties: map[string][]byte{}}, e2)

	e3, err := g.Edge("edge-knows")
	assert.Nil(t, err)
	assert.Equal(t, Edge{UID: "edge-knows", SourceUID: "node-foo", Label: "knows", TargetUID: "node-bar", Properties: map[string][]byte{"name": []byte("2020")}}, e3)

	e4, err := g.Edge("edge-owns")
	assert.Nil(t, err)
	assert.Equal(t, Edge{UID: "edge-owns", SourceUID: "node-foo", Label: "owns", TargetUID: "node-dog", Properties: map[string][]byte{}}, e4)
}
