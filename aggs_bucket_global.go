package aggretastic

// GlobalAggregation defines a single bucket of all the documents within
// the search execution context. This context is defined by the indices
// and the document types you’re searching on, but is not influenced
// by the search query itself.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/6.2/search-aggregations-bucket-global-aggregation.html
type GlobalAggregation struct {
	*tree

	meta map[string]interface{}
}

func NewGlobalAggregation() *GlobalAggregation {
	a := &GlobalAggregation{}
	a.tree = nilAggregationTree(a)

	return a
}

func (a *GlobalAggregation) SubAggregation(name string, subAggregation Aggregation) *GlobalAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *GlobalAggregation) Meta(metaData map[string]interface{}) *GlobalAggregation {
	a.meta = metaData
	return a
}

func (a *GlobalAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//         "all_products" : {
	//             "global" : {},
	//             "aggs" : {
	//                 "avg_price" : { "avg" : { "field" : "price" } }
	//             }
	//         }
	//    }
	//	}
	// This method returns only the { "global" : {} } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["global"] = opts

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
