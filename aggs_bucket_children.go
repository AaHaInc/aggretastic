package aggretastic

type ChildrenAggregation struct {
	*aggregation

	typ  string
	meta map[string]interface{}
}

func NewChildrenAggregation() *ChildrenAggregation {
	return &ChildrenAggregation{aggregation: nilAggregation()}
}

func (a *ChildrenAggregation) Type(typ string) *ChildrenAggregation {
	a.typ = typ
	return a
}

func (a *ChildrenAggregation) SubAggregation(name string, subAggregation Aggregation) *ChildrenAggregation {
	a.aggregation.setChild(subAggregation, name)
	return a
}

func (a *ChildrenAggregation) Meta(metaData map[string]interface{}) *ChildrenAggregation {
	a.meta = metaData
	return a
}

func (a *ChildrenAggregation) Source() (interface{}, error) {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["children"] = opts
	opts["type"] = a.typ

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
