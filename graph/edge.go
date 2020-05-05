package graph

// NewEdge returns a new edge instance.
func NewEdge(uid, sourceUID, label, targetUID string, kv ...KV) Edge {
	return Edge{
		UID:        uid,
		SourceUID:  sourceUID,
		Label:      label,
		TargetUID:  targetUID,
		Properties: NewProperties(kv...),
	}
}

// Edge is a edge in the graph.
type Edge struct {
	UID        string            `json:"uid"`
	SourceUID  string            `json:"source_uid"`
	Label      string            `json:"label"`
	TargetUID  string            `json:"target_uid"`
	Properties map[string][]byte `json:"properties"`
}
