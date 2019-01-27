package aggretastic

import "github.com/olivere/elastic"

type AdjacencyMatrixAggregation struct {
	*aggregation

	filters map[string]elastic.Query
	meta    map[string]interface{}
}

func NewAdjacencyMatrixAggregation() *AdjacencyMatrixAggregation {
	return &AdjacencyMatrixAggregation{
		filters:     make(map[string]elastic.Query),
		aggregation: nilAggregation(),
	}
}

// Filters adds the filter
func (a *AdjacencyMatrixAggregation) Filters(name string, filter elastic.Query) *AdjacencyMatrixAggregation {
	a.filters[name] = filter
	return a
}

// SubAggregation adds a sub-aggregation to this aggregation.
func (a *AdjacencyMatrixAggregation) SubAggregation(name string, subAggregation Aggregation) *AdjacencyMatrixAggregation {
	a.setSubAggregation(subAggregation, name)
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *AdjacencyMatrixAggregation) Meta(metaData map[string]interface{}) *AdjacencyMatrixAggregation {
	a.meta = metaData
	return a
}

func (a *AdjacencyMatrixAggregation) Source() (interface{}, error) {
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
