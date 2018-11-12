package aggretastic

// @TODO
// wip

import (
	"fmt"
)

type finiteAggregation struct {
	*aggregation
}

var errNotInjectable = fmt.Errorf("not injectable")

func newFiniteAggregation() *finiteAggregation {
	return &finiteAggregation{}
}

func IsNotInjectable(agg Aggregation) bool {
	_, ok := agg.(*finiteAggregation)
	return ok
}

func (a *finiteAggregation) Inject(subAggregation Aggregation, path ...string) error {
	return errNotInjectable
}

func (a *finiteAggregation) InjectX(subAggregation Aggregation, path ...string) error {
	return errNotInjectable
}

func (a *finiteAggregation) GetAllSubs() map[string]Aggregation {
	return nil
}

func (a *finiteAggregation) Select(path ...string) Aggregation {
	fmt.Printf("finiteAggregation.Select() is ignored. The aggregation doesn't allow to have subAggregations")
	return nil
}

func (a *finiteAggregation) WrapBy(wrapper Aggregation, name string) {
	return
}
func (a *finiteAggregation) InjectWrapper(wrapper Aggregation, path ...string) error {
	return errNotInjectable
}

func (a *finiteAggregation) Source() (interface{}, error) {
	return nil, nil
}
