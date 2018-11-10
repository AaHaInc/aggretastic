package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic"
)

// Aggregation is a tree-ish version of original elastic.Aggregation
// Besides just attaching subAggregations it can get any of children subAggregations
// and add another subAggregation to it
type Aggregation interface {
	// embedding original elastic.Aggregation interface
	// is used to support call of `.Source()` method from aggregations' code
	elastic.Aggregation

	// GetAllSubs returns the map of this aggregation's subAggregations
	GetAllSubs() map[string]Aggregation

	// Inject sets new subAgg into the map of subAggregations
	Inject(subAgg elastic.Aggregation, path ...string) error

	// InjectX sets new subAgg into the map of subAggregations only if it NOT exists already
	InjectX(subAgg elastic.Aggregation, path ...string) error

	// Select returns any subAgg by it's path
	Select(path ...string) Aggregation

	// Pop returns a subAgg by it's path and remove it from tree
	Pop(path ...string) Aggregation

	// Export returns the same object in original Agg interface
	Export() elastic.Aggregation
}

func IsNilTree(t Aggregation) bool {
	return t == nil || t.Export() == nil
}

type tree struct {
	root            elastic.Aggregation
	subAggregations map[string]elastic.Aggregation
}

func nilAggregationTree(root elastic.Aggregation) *tree {
	return &tree{
		root:            root,
		subAggregations: make(map[string]elastic.Aggregation),
	}
}

func (a *tree) Inject(subAggregation elastic.Aggregation, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("no path")
	}

	if len(path) == 1 {
		a.subAggregations[path[0]] = subAggregation
		return nil
	}

	// deeper inject
	cursor := a.Select(path[:len(path)-1]...)
	if IsNilTree(cursor) {
		return fmt.Errorf("path not selectable")
	}

	return cursor.Inject(subAggregation, path[len(path)-1])
}

func (a *tree) InjectX(subAggregation elastic.Aggregation, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("no path")
	}

	if alreadyInjected := a.Select(path...); IsNilTree(alreadyInjected) {
		return a.Inject(subAggregation, path...)
	}

	return nil
}

func (a *tree) GetAllSubs() map[string]Aggregation {
	result := make(map[string]Aggregation)

	for k, v := range a.subAggregations {
		result[k] = v.(Aggregation)
	}

	return result
}

func (a *tree) Select(path ...string) Aggregation {
	if len(path) == 0 {
		return nil
	}

	subAgg, ok := a.subAggregations[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		return subAgg.(Aggregation)
	}

	return subAgg.(Aggregation).Select(path[1:]...)
}

func (a *tree) Pop(path ...string) Aggregation {
	if len(path) == 0 {
		return nil
	}

	subAgg, ok := a.subAggregations[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		delete(a.subAggregations, path[0])
		return subAgg.(Aggregation)
	}

	return subAgg.(Aggregation).Pop(path[1:]...)
}

func (a *tree) Export() elastic.Aggregation {
	return a.root
}

// Shorthand type for collection of Aggregations
type Aggregations map[string]Aggregation

// Export does export() on the map of aggregations
func (a *Aggregations) Export() map[string]elastic.Aggregation {
	result := make(map[string]elastic.Aggregation)

	if a == nil {
		return result
	}

	for k, v := range *a {
		result[k] = v.Export()
	}

	return result
}

// Select selects an aggregation from the map (going deep forwarding the agg.Select() method)
func (a *Aggregations) Select(path ...string) Aggregation {
	if len(path) == 0 {
		return nil
	}

	base, ok := (*a)[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		return base
	}

	return base.Select(path[1:]...)
}

// Pop pops an aggregation from the map (going deep forwarding the agg.Pop() method)
func (a *Aggregations) Pop(path ...string) Aggregation {
	if len(path) == 0 {
		return nil
	}

	base, ok := (*a)[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		delete(*a, path[0])
		return base
	}

	return base.Pop(path[1:]...)
}

// Inject just puts agg into the map of aggregations
func (a *Aggregations) Inject(subAgg elastic.Aggregation, path ...string) error {
	if a == nil {
		return fmt.Errorf("not injectable")
	}

	if len(path) == 0 {
		return fmt.Errorf("no path")
	}

	name := path[0]

	if len(path) == 1 {
		(*a)[name] = subAgg.(Aggregation)
		return nil
	}

	path = path[1:]
	if _, ok := (*a)[name]; !ok {
		return fmt.Errorf("not injectable")
	}

	return (*a)[name].Inject(subAgg, path...)
}

func (a *Aggregations) InjectX(subAgg elastic.Aggregation, path ...string) error {
	if a == nil {
		return fmt.Errorf("not injectable")
	}

	if len(path) == 0 {
		return fmt.Errorf("no path")
	}

	name := path[0]

	if len(path) == 1 {
		if _, ok := (*a)[name]; !ok {
			(*a)[name] = subAgg.(Aggregation)
		}

		return nil
	}

	path = path[1:]
	if _, ok := (*a)[name]; !ok {
		return fmt.Errorf("not injectable")
	}

	return (*a)[name].InjectX(subAgg, path...)
}
