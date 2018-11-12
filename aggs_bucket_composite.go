package aggretastic

import "github.com/olivere/elastic"

type CompositeAggregation struct {
	*aggregation

	after   map[string]interface{}
	size    *int
	sources []elastic.CompositeAggregationValuesSource
	meta    map[string]interface{}
}

func NewCompositeAggregation() *CompositeAggregation {
	a := &CompositeAggregation{sources: make([]elastic.CompositeAggregationValuesSource, 0)}
	a.aggregation = nilAggregation()

	return a
}

func (a *CompositeAggregation) Size(size int) *CompositeAggregation {
	a.size = &size
	return a
}

func (a *CompositeAggregation) AggregateAfter(after map[string]interface{}) *CompositeAggregation {
	a.after = after
	return a
}

func (a *CompositeAggregation) Sources(sources ...elastic.CompositeAggregationValuesSource) *CompositeAggregation {
	a.sources = append(a.sources, sources...)
	return a
}

func (a *CompositeAggregation) SubAggregation(name string, subAggregation Aggregation) *CompositeAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

func (a *CompositeAggregation) Meta(metaData map[string]interface{}) *CompositeAggregation {
	a.meta = metaData
	return a
}

func (a *CompositeAggregation) Source() (interface{}, error) {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["composite"] = opts

	sources := make([]interface{}, len(a.sources))
	for i, s := range a.sources {
		src, err := s.Source()
		if err != nil {
			return nil, err
		}
		sources[i] = src
	}
	opts["sources"] = sources

	if a.size != nil {
		opts["size"] = *a.size
	}

	if a.after != nil {
		opts["after"] = a.after
	}

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
