package aggretastic

import "github.com/olivere/elastic/v7"

// AdjacencyMatrixAggregation returning a form of adjacency matrix.
// The request provides a collection of named filter expressions,
// similar to the filters aggregation request. Each bucket in the
// response represents a non-empty cell in the matrix of intersecting filters.
//
// For details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/6.2/search-aggregations-bucket-adjacency-matrix-aggregation.html
type AdjacencyMatrixAggregation struct {
	*tree

	filters map[string]elastic.Query
	meta    map[string]interface{}
}

// NewAdjacencyMatrixAggregation initializes a new AdjacencyMatrixAggregation.
func NewAdjacencyMatrixAggregation() *AdjacencyMatrixAggregation {
	a := &AdjacencyMatrixAggregation{filters: make(map[string]elastic.Query)}
	a.tree = nilAggregationTree(a)

	return a
}

// Filters adds the filter
func (a *AdjacencyMatrixAggregation) Filters(name string, filter elastic.Query) *AdjacencyMatrixAggregation {
	a.filters[name] = filter
	return a
}

// SubAggregation adds a sub-aggregation to this aggregation.
func (a *AdjacencyMatrixAggregation) SubAggregation(name string, subAggregation Aggregation) *AdjacencyMatrixAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *AdjacencyMatrixAggregation) Meta(metaData map[string]interface{}) *AdjacencyMatrixAggregation {
	a.meta = metaData
	return a
}

// Source returns the a JSON-serializable interface.
func (a *AdjacencyMatrixAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//  "aggs" : {
	//		"interactions" : {
	//			"adjacency_matrix" : {
	//				"filters" : {
	//					"grpA" : { "terms" : { "accounts" : ["hillary", "sidney"] }},
	//					"grpB" : { "terms" : { "accounts" : ["donald", "mitt"] }},
	//					"grpC" : { "terms" : { "accounts" : ["vladimir", "nigel"] }}
	//				}
	//			}
	//		}
	//	}
	// This method returns only the (outer) { "adjacency_matrix" : {} } part.

	source := make(map[string]interface{})
	adjacencyMatrix := make(map[string]interface{})
	source["adjacency_matrix"] = adjacencyMatrix

	dict := make(map[string]interface{})
	for key, filter := range a.filters {
		src, err := filter.Source()
		if err != nil {
			return nil, err
		}
		dict[key] = src
	}
	adjacencyMatrix["filters"] = dict

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}
