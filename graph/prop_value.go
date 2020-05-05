package graph

// KV is a property key and value pair.
type KV struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

// NewProperties takes one or more key value pairs and returns a property map.
func NewProperties(kv ...KV) map[string][]byte {
	props := make(map[string][]byte)

	for _, pair := range kv {
		props[pair.Key] = pair.Value
	}

	return props
}
