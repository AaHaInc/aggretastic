package aggretastic

import (
	"github.com/olivere/elastic"
)

type FilterAggregation struct {
	*aggregation

	filter elastic.Query
	meta   map[string]interface{}
}

func NewFilterAggregation() *FilterAggregation {
	return &FilterAggregation{aggregation: nilAggregation()}
}

func (a *FilterAggregation) SubAggregation(name string, subAggregation Aggregation) *FilterAggregation {
	a.setSubAggregation(subAggregation, name)
	return a
}

func (a *FilterAggregation) Meta(metaData map[string]interface{}) *FilterAggregation {
	a.meta = metaData
	return a
}

func (a *FilterAggregation) Filter(filter elastic.Query) *FilterAggregation {
	a.filter = filter
	return a
}

func (a *FilterAggregation) Source() (interface{}, error) {
	src, err := a.filter.Source()
	if err != nil {
		return nil, err
	}
	source := make(map[string]interface{})
	source["filter"] = src

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
