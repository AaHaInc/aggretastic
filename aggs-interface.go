//go:generate ./sync.sh

// The aggretastic package enables manipulation of
// a olivere/elastic subAggregations tree with high fidelity.
//
//	Injectable: https://godoc.org/github.com/AaHaInc/aggretastic/#Injectable
//	NotInjectable: https://godoc.org/github.com/AaHaInc/aggretastic/#NotInjectable
//	Aggregations: https://godoc.org/github.com/AaHaInc/aggretastic/#Aggregations
//
//
package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic"
)

var (
	ErrNoPath             = fmt.Errorf("no path")
	ErrPathNotSelectable  = fmt.Errorf("path is not selectable")
	ErrAggIsNotInjectable = fmt.Errorf("agg is not injectable")
)

// Aggregation is extended version of original elastic.Aggregation
// Besides just attaching subAggregations it can get any of children subAggregations
// and add another subAggregation to it
type Aggregation interface {
	// embedding original elastic.Aggregation interface
	// is used to support call of `.Source()` method from aggregations' code
	elastic.Aggregation

	// GetAllSubs returns the map of this aggregation's subAggregations
	GetAllSubs() map[string]Aggregation

	// Inject sets new subAgg into the map of subAggregations
	Inject(subAgg Aggregation, path ...string) error

	// InjectX sets new subAgg into the map of subAggregations only if it NOT exists already
	InjectX(subAgg Aggregation, path ...string) error

	// Select returns any subAgg by it's path
	Select(path ...string) Aggregation

	// Pop returns a subAgg by it's path and remove it from Injectable
	Pop(path ...string) Aggregation

	// Export returns the same object in original Agg interface
	Export() elastic.Aggregation
}