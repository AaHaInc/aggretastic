package aggretastic

import "github.com/olivere/elastic"

type NestedAggregation struct {
	*aggregation
}

func NewNestedAggregation() *NestedAggregation {
	return &NestedAggregation{aggregation: nilAggregation(elastic.NewNestedAggregation())}
}

func (a *NestedAggregation) SubAggregation(name string, subAggregation Aggregation) *NestedAggregation {
	a.base.(*elastic.NestedAggregation).SubAggregation(name, subAggregation)
	return a
}

func (a *NestedAggregation) Base() *elastic.NestedAggregation {
	return a.base.(*elastic.NestedAggregation)
}
