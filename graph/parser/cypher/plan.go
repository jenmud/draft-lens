package cypher

// toIfaceSlice is a helper function that takes a interface
// returns a slice of interfaces.
func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

// KV is a key/value pair
type KV struct {
	Key   string
	Value []byte
}

// ReadingClause is a immutable read/query.
type ReadingClause struct {
	Matches []Match
	Returns []string
}

// Match is the match query.
type Match struct {
	Nodes []Node
}

// Node is a node used for a query.
type Node struct {
	Variable   string
	UID        string
	Labels     []string
	Properties map[string][]byte
}

// QueryPlan is a query plan for applying a query.
type QueryPlan struct {
	ReadingClause []ReadingClause
}
