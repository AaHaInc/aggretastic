package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic/v7"
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

func (a *notInjectable) ExtractLeafPaths() (leafs [][]string) {
	return make([][]string, 0)
}

func (a *notInjectable) Inject(subAggregation Aggregation, path ...string) (resultPaths [][]string, err error) {
	err = ErrAggIsNotInjectable
	return
}

func (a *notInjectable) InjectX(subAggregation Aggregation, path ...string) (resultPaths [][]string, err error) {
	err = ErrAggIsNotInjectable
	return
}

func (a *notInjectable) InjectSafe(subAggregation Aggregation, path ...string) (resultPaths [][]string, err error) {
	err = ErrAggIsNotInjectable
	return
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
