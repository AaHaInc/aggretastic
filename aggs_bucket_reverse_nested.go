package aggretastic

type ReverseNestedAggregation struct {
	*aggregation

	path string
	meta map[string]interface{}
}

func NewReverseNestedAggregation() *ReverseNestedAggregation {
	return &ReverseNestedAggregation{aggregation: nilAggregation()}
}

func (a *ReverseNestedAggregation) Path(path string) *ReverseNestedAggregation {
	a.path = path
	return a
}

func (a *ReverseNestedAggregation) SubAggregation(name string, subAggregation Aggregation) *ReverseNestedAggregation {
	a.aggregation.setChild(subAggregation, name)
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *ReverseNestedAggregation) Meta(metaData map[string]interface{}) *ReverseNestedAggregation {
	a.meta = metaData
	return a
}

func (a *ReverseNestedAggregation) Source() (interface{}, error) {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["reverse_nested"] = opts

	if a.path != "" {
		opts["path"] = a.path
	}

	// AggregationBuilder (SubAggregations)
	if len(a.children) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.children {
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
