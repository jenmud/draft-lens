package graph

// convertPropertiesToKV converts a Property map to a array of key values.
func convertPropertiesToKV(props map[string][]byte) []KV {
	kvs := make([]KV, len(props))

	count := 0
	for k, v := range props {
		kvs[count] = KV{Key: k, Value: v}
	}

	return kvs
}
