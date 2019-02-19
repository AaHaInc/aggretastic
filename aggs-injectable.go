package aggretastic

import (
	"github.com/olivere/elastic"
)

// Injectable represents elastic.Aggregation which can
// get any of children subAggregations
// and add another subAggregation to it
type Injectable struct {
	root            elastic.Aggregation
	subAggregations map[string]Aggregation
}

// Creates new Injectable Aggregation
func newInjectable(root elastic.Aggregation) *Injectable {
	return &Injectable{
		root:            root,
		subAggregations: make(map[string]Aggregation),
	}
}


// Checks if Aggregation is Injectable
func IsInjectable(agg Aggregation) bool {
	_, ok := agg.(*Injectable)
	return ok
}

// Checks if Aggregation is nil
func IsNilTree(t Aggregation) bool {
	return t == nil || t.Export() == nil
}

// Inject sets new subAgg into the map of subAggregations
func (a *Injectable) Inject(subAggregation Aggregation, path ...string) error {
	if len(path) == 0 {
		return ErrNoPath
	}

	if len(path) == 1 {
		a.subAggregations[path[0]] = subAggregation
		return nil
	}

	// deeper inject
	cursor := a.Select(path[:len(path)-1]...)
	if IsNilTree(cursor) {
		return ErrPathNotSelectable
	}

	return cursor.Inject(subAggregation, path[len(path)-1])
}

// InjectX sets new subAgg into the map of subAggregations only if it NOT exists already
func (a *Injectable) InjectX(subAggregation Aggregation, path ...string) error {
	if len(path) == 0 {
		return ErrNoPath
	}

	if alreadyInjected := a.Select(path...); IsNilTree(alreadyInjected) {
		return a.Inject(subAggregation, path...)
	}

	return nil
}

// GetAllSubs returns the map of this aggregation's subAggregations
func (a *Injectable) GetAllSubs() map[string]Aggregation {
	return a.subAggregations
}

// Select returns any subAgg by it's path
func (a *Injectable) Select(path ...string) Aggregation {
	if len(path) == 0 {
		return nil
	}

	subAgg, ok := a.subAggregations[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		return subAgg
	}

	return subAgg.Select(path[1:]...)
}

// Pop returns a subAgg by it's path and remove it from Injectable
func (a *Injectable) Pop(path ...string) Aggregation {
	if len(path) == 0 {
		return nil
	}

	subAgg, ok := a.subAggregations[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		delete(a.subAggregations, path[0])
		return subAgg
	}

	return subAgg.Pop(path[1:]...)
}

// Export returns the same object in original Agg interface
func (a *Injectable) Export() elastic.Aggregation {
	return a.root
}

// Source returns a JSON-serializable aggregation that is a fragment
// of the request sent to Elasticsearch.
func (a *Injectable) Source() (interface{}, error) {
	return a.root.Source()
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
func (a *Aggregations) Inject(subAgg Aggregation, path ...string) error {
	if a == nil {
		return ErrAggIsNotInjectable
	}

	if len(path) == 0 {
		return ErrNoPath
	}

	name := path[0]

	if len(path) == 1 {
		(*a)[name] = subAgg
		return nil
	}

	path = path[1:]
	if _, ok := (*a)[name]; !ok {
		return ErrAggIsNotInjectable
	}

	return (*a)[name].Inject(subAgg, path...)
}

// Inject just puts agg into the map of aggregations only if it NOT exists already
func (a *Aggregations) InjectX(subAgg Aggregation, path ...string) error {
	if a == nil {
		return ErrAggIsNotInjectable
	}

	if len(path) == 0 {
		return ErrNoPath
	}

	name := path[0]

	if len(path) == 1 {
		if _, ok := (*a)[name]; !ok {
			(*a)[name] = subAgg
		}

		return nil
	}

	path = path[1:]
	if _, ok := (*a)[name]; !ok {
		return ErrAggIsNotInjectable
	}

	return (*a)[name].InjectX(subAgg, path...)
}