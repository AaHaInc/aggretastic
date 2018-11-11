package aggretastic

import "github.com/olivere/elastic"

type ReverseNestedAggregation struct {
	*aggregation
}

func NewReverseNestedAggregation() *ReverseNestedAggregation {
	return &ReverseNestedAggregation{aggregation: nilAggregation(elastic.NewReverseNestedAggregation())}
}

func (a *ReverseNestedAggregation) SubAggregation(name string, subAggregation Aggregation) *ReverseNestedAggregation {
	a.base.(*elastic.ReverseNestedAggregation).SubAggregation(name, subAggregation)
	return a
}

func (a *ReverseNestedAggregation) Base() *elastic.ReverseNestedAggregation {
	return a.base.(*elastic.ReverseNestedAggregation)
}
