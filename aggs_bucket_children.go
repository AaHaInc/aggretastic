package aggretastic

import "github.com/olivere/elastic"

type ChildrenAggregation struct {
	*aggregation
}

func NewChildrenAggregation() *ChildrenAggregation {
	return &ChildrenAggregation{
		aggregation: nilAggregation(elastic.NewChildrenAggregation()),
	}
}

func (a *ChildrenAggregation) SubAggregation(name string, subAggregation Aggregation) *ChildrenAggregation {
	a.aggregation.base.(*elastic.ChildrenAggregation).SubAggregation(name, subAggregation)
	return a
}

func (a *ChildrenAggregation) Base() *elastic.ChildrenAggregation {
	return a.base.(*elastic.ChildrenAggregation)
}
