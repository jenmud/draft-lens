package iterator

// New returns a new iterator.
func New(items []interface{}) *Iterator {
	return &Iterator{current: -1, data: items}
}

// Iterator is an iterator.
type Iterator struct {
	current int
	data    []interface{}
}

// Value returns the current value.
func (it *Iterator) Value() interface{} {
	if it.Size() == 0 {
		return nil
	}

	if it.current == -1 {
		return it.data[0]
	}

	return it.data[it.current]
}

// Next progresses the iterator returning true if there are still items to iterator over.
func (it *Iterator) Next() bool {
	it.current++

	if it.current >= len(it.data) {
		it.current--
		return false
	}

	return true
}

// Size returns the count of items in the iterator.
func (it *Iterator) Size() int {
	return len(it.data)
}

// Channel returns the items in the iterator as a channel.
func (it *Iterator) Channel() <-chan interface{} {
	out := make(chan interface{}, it.Size())
	for it.Next() {
		out <- it.Value()
	}
	close(out)
	return out
}
