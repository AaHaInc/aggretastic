package aggretastic

// ChildrenAggregation is a special single bucket aggregation that enables
// aggregating from buckets on parent document types to buckets on child documents.
// It is available from 1.4.0.Beta1 upwards.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/6.2/search-aggregations-bucket-children-aggregation.html
type ChildrenAggregation struct {
	*tree

	typ  string
	meta map[string]interface{}
}

func NewChildrenAggregation() *ChildrenAggregation {
	a := &ChildrenAggregation{}
	a.tree = nilAggregationTree(a)

	return a
}

func (a *ChildrenAggregation) Type(typ string) *ChildrenAggregation {
	a.typ = typ
	return a
}

func (a *ChildrenAggregation) SubAggregation(name string, subAggregation Aggregation) *ChildrenAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *ChildrenAggregation) Meta(metaData map[string]interface{}) *ChildrenAggregation {
	a.meta = metaData
	return a
}

func (a *ChildrenAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//      "to-answers" : {
	//        "children": {
	//          "type" : "answer"
	//        }
	//      }
	//    }
	//	}
	// This method returns only the { "type" : ... } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["children"] = opts
	opts["type"] = a.typ

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
