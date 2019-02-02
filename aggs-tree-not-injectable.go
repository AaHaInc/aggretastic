package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic"
)

type notInjectable struct {
	root elastic.Aggregation
}

func newNotInjectable(root elastic.Aggregation) *notInjectable {
	return &notInjectable{root: root}
}

func IsNotInjectable(agg Aggregation) bool {
	_, ok := agg.(*notInjectable)
	return ok
}

func (a *notInjectable) Inject(subAggregation elastic.Aggregation, path ...string) error {
	return ErrAggIsNotInjectable
}

func (a *notInjectable) InjectX(subAggregation elastic.Aggregation, path ...string) error {
	return ErrAggIsNotInjectable
}

func (a *notInjectable) GetAllSubs() map[string]Aggregation {
	return nil
}

func (a *notInjectable) Select(path ...string) Aggregation {
	// nothing to select because of no subAggregations
	s, _ := a.root.Source()
	fmt.Printf("notInjectable.Select() is ignored. The aggregation doesn't allow to have subAggregations. Root: %s", s)
	return nil
}

func (a *notInjectable) Pop(path ...string) Aggregation {
	// nothing to select because of no subAggregations
	s, _ := a.root.Source()
	fmt.Printf("notInjectable.Select() is ignored. The aggregation doesn't allow to have subAggregations. Root: %s", s)
	return nil
}

func (a *notInjectable) Export() elastic.Aggregation {
	return a.root
}

func (a *notInjectable) Source() (interface{}, error) {
	return a.root.Source()
}
