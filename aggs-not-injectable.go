package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic"
)

// Aggregation which can't store subAggregations or inject it
//
type NotInjectable struct {
	root elastic.Aggregation
}

// Creates new NotInjectable Aggregation
//
func newNotInjectable(root elastic.Aggregation) *NotInjectable {
	return &NotInjectable{root: root}
}

// Checks if Aggregation is NotInjectable
//
func IsNotInjectable(agg Aggregation) bool {
	_, ok := agg.(*NotInjectable)
	return ok
}

// Always returns isNotInjectable error
//
func (a *NotInjectable) Inject(subAggregation Aggregation, path ...string) error {
	return ErrAggIsNotInjectable
}

// Always returns isNotInjectable error
//
func (a *NotInjectable) InjectX(subAggregation Aggregation, path ...string) error {
	return ErrAggIsNotInjectable
}

// Always returns nil
//
func (a *NotInjectable) GetAllSubs() map[string]Aggregation {
	return nil
}

// Always returns nil
//
func (a *NotInjectable) Select(path ...string) Aggregation {
	// nothing to select because of no subAggregations
	s, _ := a.root.Source()
	fmt.Printf("NotInjectable.Select() is ignored. The aggregation doesn't allow to have subAggregations. Root: %s", s)
	return nil
}

// Always returns nil
//
func (a *NotInjectable) Pop(path ...string) Aggregation {
	// nothing to select because of no subAggregations
	s, _ := a.root.Source()
	fmt.Printf("NotInjectable.Select() is ignored. The aggregation doesn't allow to have subAggregations. Root: %s", s)
	return nil
}

// Export returns the same object in original Agg interface
//
func (a *NotInjectable) Export() elastic.Aggregation {
	return a.root
}

// Source returns a JSON-serializable aggregation that is a fragment
// of the request sent to Elasticsearch.
//
func (a *NotInjectable) Source() (interface{}, error) {
	return a.root.Source()
}
