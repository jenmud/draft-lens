package cypher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Query       string
	Expected    QueryPlan
	Name        string
	ShouldError bool
}

func TestMatchQueries(t *testing.T) {
	tests := []TestCase{
		TestCase{
			Name:  "SingleMatchNoConditions",
			Query: `MATCH (n) RETURN n`,
			Expected: QueryPlan{
				ReadingClause: []ReadingClause{
					ReadingClause{
						Returns: []string{"n"},
						Matches: []Match{
							Match{
								Nodes: []Node{
									Node{
										Variable: "n",
									},
								},
							},
						},
					},
				},
			},
		},
		TestCase{
			Name: "SingleMatchSingleLabelMultiline",
			Query: `
			MATCH (n:Person)
			RETURN n
			`,
			Expected: QueryPlan{
				ReadingClause: []ReadingClause{
					ReadingClause{
						Returns: []string{"n"},
						Matches: []Match{
							Match{
								Nodes: []Node{
									Node{
										Variable: "n",
										Labels:   []string{"Person"},
									},
								},
							},
						},
					},
				},
			},
		},
		TestCase{
			Name: "MultipleMatchSingleLabel",
			Query: `
			MATCH (n:Person)
			MATCH (m:Animal)
			RETURN n, m
			`,
			Expected: QueryPlan{
				ReadingClause: []ReadingClause{
					ReadingClause{
						Returns: []string{"n", "m"},
						Matches: []Match{
							Match{
								Nodes: []Node{
									Node{
										Variable: "n",
										Labels:   []string{"Person"},
									},
								},
							},
							Match{
								Nodes: []Node{
									Node{
										Variable: "m",
										Labels:   []string{"Animal"},
									},
								},
							},
						},
					},
				},
			},
		},
		TestCase{
			Name: "MultipleMatchReturnsVarsDontMatch",
			Query: `
			MATCH (n:Person)
			MATCH (m:Animal)
			RETURN n, missing
			`,
			Expected:    QueryPlan{},
			ShouldError: true,
		},
		TestCase{
			Name:        "SingleMatchSingleLabelNoReturn",
			Query:       `MATCH (n:Person)`,
			ShouldError: true,
			Expected:    QueryPlan{},
		},
		TestCase{
			Name:        "SingleMatchMultipleLabelNotSupported",
			Query:       `MATCH (n:Person:Mutliple:Not:Supported)`,
			ShouldError: true,
			Expected:    QueryPlan{},
		},
		TestCase{
			Name:  "SingleMatchSingleLabelWithMultipleProperties",
			Query: `MATCH (n:Person {name: "Foo", surname: 'Bar', age: 21, active: true, address: "My address is private"}) RETURN n`,
			Expected: QueryPlan{
				ReadingClause: []ReadingClause{
					ReadingClause{
						Returns: []string{"n"},
						Matches: []Match{
							Match{
								Nodes: []Node{
									Node{
										Variable: "n",
										Labels:   []string{"Person"},
										Properties: map[string][]byte{
											"name":    []byte("Foo"),
											"surname": []byte("Bar"),
											"age":     []byte("21"),
											"active":  []byte("true"),
											"address": []byte("My address is private"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		TestCase{
			Name:        "SingleMatchSingleLabelWithMultipleProperties",
			Query:       `MATCH (n:Person {name: "Foo", name: 'Bar'}) RETURN n`,
			ShouldError: true,
			Expected:    QueryPlan{},
		},
	}

	for _, test := range tests {
		got, err := Parse("", []byte(test.Query))
		if !test.ShouldError {
			assert.Nil(t, err, "%s did not expect an error to be raises: %s (Query: %s)", test.Name, err, test.Query)
			actual := got.(QueryPlan)
			assert.Equal(t, test.Expected, actual, "%s expected %#v but got %#v", test.Name, test.Expected, actual)
		} else {
			assert.NotNil(t, err, "%s Expected query %s to fail", test.Name, test.Query)
		}
	}
}
