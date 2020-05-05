package graph

// Iterator is an iterator interface for iterating over a set of items.
type Iterator interface {
	// Value returns the Item.
	Value() interface{}
	// Next progresses the iterator and return true if there are still items to iterate over.
	Next() bool
	// Size returns the total item count in the iterator
	Size() int
	// Channel returns the items in the iterator as a channel.
	Channel() <-chan interface{}
}
