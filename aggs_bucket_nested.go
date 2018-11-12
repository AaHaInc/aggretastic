package aggretastic

type NestedAggregation struct {
	*aggregation

	path string
	meta map[string]interface{}
}

func NewNestedAggregation() *NestedAggregation {
	return &NestedAggregation{aggregation: nilAggregation()}
}

func (a *NestedAggregation) SubAggregation(name string, subAggregation Aggregation) *NestedAggregation {
	a.setSubAggregation(subAggregation, name)
	return a
}

func (a *NestedAggregation) Meta(metaData map[string]interface{}) *NestedAggregation {
	a.meta = metaData
	return a
}

func (a *NestedAggregation) Path(path string) *NestedAggregation {
	a.path = path
	return a
}

func (a *NestedAggregation) Source() (interface{}, error) {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["nested"] = opts

	opts["path"] = a.path

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
