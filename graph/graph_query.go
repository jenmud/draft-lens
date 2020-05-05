package graph

import (
	"log"

	"github.com/jenmud/draft/graph/parser/cypher"
)

// Query takes a query string and returns a subgraph containing
// the query results.
func (g *Graph) Query(query string) (*Graph, error) {
	subg := New()

	queryResult, err := cypher.Parse("", []byte(query))
	if err != nil {
		return nil, err
	}

	// search for nodes
	for _, rc := range queryResult.(cypher.QueryPlan).ReadingClause {
		for _, match := range rc.Matches {
			for _, node := range match.Nodes {
				nodes := g.NodesBy(node.Labels, node.Properties)
				for nodes.Next() {
					node := nodes.Value().(Node)
					if _, err := subg.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
						log.Printf("[Query] %s", err)
					}
				}
			}
		}
	}

	return subg, nil
}
