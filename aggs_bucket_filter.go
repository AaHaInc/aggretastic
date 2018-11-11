package aggretastic

import "github.com/olivere/elastic"

type FilterAggregation struct {
	*aggregation
}

func NewFilterAggregation() *FilterAggregation {
	return &FilterAggregation{aggregation: nilAggregation(elastic.NewFilterAggregation())}
}

func (a *FilterAggregation) SubAggregation(name string, subAggregation Aggregation) *FilterAggregation {
	a.aggregation.base.(*elastic.FilterAggregation).SubAggregation(name, subAggregation)
	return a
}

func (a *FilterAggregation) Base() *elastic.FilterAggregation {
	return a.base.(*elastic.FilterAggregation)
}
