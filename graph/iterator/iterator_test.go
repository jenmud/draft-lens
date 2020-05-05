package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	iter := New([]interface{}{1, 2})
	assert.Equal(t, 1, iter.Value())
	assert.Equal(t, 1, iter.Value()) // check that calling again returns the same item
}

func TestValue__no_items(t *testing.T) {
	iter := New([]interface{}{})
	assert.Equal(t, nil, iter.Value())
	assert.Equal(t, nil, iter.Value()) // check that calling again returns the same item
}

func TestNext(t *testing.T) {
	iter := New([]interface{}{1, 2, 3})
	assert.Equal(t, true, iter.Next(), "expected true but got false (Expected value: %v, Actual value: %v)", 1, iter.Value())
	assert.Equal(t, true, iter.Next(), "expected true but got false (Expected value: %v, Actual value: %v)", 2, iter.Value())
	assert.Equal(t, true, iter.Next(), "expected true but got false (Expected value: %v, Actual value: %v)", 3, iter.Value())
	assert.Equal(t, false, iter.Next(), "expected false but got true (Expected value: %v, Actual value: %v)", nil, iter.Value())
}

func TestSize(t *testing.T) {
	iter := New([]interface{}{1, 2, 3})
	assert.Equal(t, 3, iter.Size())
}

func TestChannel(t *testing.T) {
	iter := New([]interface{}{1, 2, 3})

	expected := []interface{}{1, 2, 3}
	actual := []interface{}{}

	for item := range iter.Channel() {
		actual = append(actual, item)
	}

	assert.ElementsMatch(t, expected, actual)
}
