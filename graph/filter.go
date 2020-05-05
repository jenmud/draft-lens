package graph

// FilterType indicates what to apply the filtering on.
type FilterType int

const (
	// LABEL is used for filtering by label.
	LABEL FilterType = iota
	// PROPERTY is used for filtering on properties.
	PROPERTY
)

// ItemType indicates the type of item.
type ItemType int

const (
	// NODE is a node item type.
	NODE ItemType = iota
	// EDGE is a edge item type.
	EDGE
)
