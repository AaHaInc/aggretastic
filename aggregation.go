package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic"
	"log"
)

type Aggregation interface {
	elastic.Aggregation

	getKey() string
	setKey(key string)
	setParent(parent Aggregation)

	setChild(child Aggregation, childName string)
	getChild(childName string) Aggregation
	removeChild(childName string)
	getChildren() []string

	WrapBy(wrapper Aggregation, name string) error
	Inject(agg Aggregation, path ...string) error
	Select(path ...string) Aggregation
}

type aggregation struct {
	base elastic.Aggregation

	key      string
	parent   Aggregation
	children map[string]Aggregation
}

func (a *aggregation) setParent(parent Aggregation) {
	a.parent = parent
}

func (a *aggregation) setChild(child Aggregation, childName string) {
	a.children[childName] = child
	a.addSubAggregation(child, childName)
}

func (a *aggregation) getChild(childName string) Aggregation {
	return a.children[childName]
}

func (a *aggregation) getChildren() []string {
	r := make([]string, 0)
	for k := range a.children {
		r = append(r, k)
	}
	return r
}

func (a *aggregation) removeChild(childName string) {
	delete(a.children, childName)
	a.base
}

func (a *aggregation) setKey(key string) {
	a.key = key
}
func (a *aggregation) getKey() string {
	return a.key
}

func (a *aggregation) Select(path ...string) Aggregation {
	if len(path) == 0 {
		return nil
	}

	subAgg, ok := a.children[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		return subAgg
	}

	return subAgg.Select(path[1:]...)
}

func (a *aggregation) Source() (interface{}, error) {
	return a.base.Source()
}

func (a *aggregation) WrapBy(wrapper Aggregation, name string) error {


	// clean the parent
	log.Println("SET CHILD " + name)
	a.parent.setChild(wrapper, name)
	log.Println("REMOVE " + a.key)
	a.parent.removeChild(a.key)

	injected := wrapper.Inject(a, a.key)


	return injected
}

func (a *aggregation) Inject(subAggregation Aggregation, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("no path")
	}

	if len(path) == 1 {
		// injection means setting key, parent, child
		subAggregation.setParent(a)
		subAggregation.setKey(path[0])
		a.setChild(subAggregation, path[0])

		// and setting elastic's base subAggregation
		return a.subAggregation(subAggregation, path[0])
	}

	// use select+recursive inject to handle long path

	cursor := a.Select(path[:len(path)-1]...)
	if cursor == nil {
		return fmt.Errorf("path not selectable")
	}

	return cursor.Inject(subAggregation, path[len(path)-1])
}

func (a *aggregation) addSubAggregation(subAggregation Aggregation, name string) error {
	if childrenAgg, ok := a.base.(*elastic.ChildrenAggregation); ok {
		childrenAgg.SubAggregation(name, subAggregation)
	} else if filterAgg, ok := a.base.(*elastic.FilterAggregation); ok {
		filterAgg.SubAggregation(name, subAggregation)
	} else if nestedAgg, ok := a.base.(*elastic.NestedAggregation); ok {
		nestedAgg.SubAggregation(name, subAggregation)
	} else if reverseNestedAgg, ok := a.base.(*elastic.ReverseNestedAggregation); ok {
		reverseNestedAgg.SubAggregation(name, subAggregation)
	} else {
		return fmt.Errorf("unknown aggregation type")
	}

	return nil
}

func (a *aggregation) delSubAggregation(name string) error {
	if childrenAgg, ok := a.base.(*elastic.ChildrenAggregation); ok {
		childrenAgg.SubAggregation(name, subAggregation)
	} else if filterAgg, ok := a.base.(*elastic.FilterAggregation); ok {
		filterAgg.SubAggregation(name, subAggregation)
	} else if nestedAgg, ok := a.base.(*elastic.NestedAggregation); ok {
		nestedAgg.SubAggregation(name, subAggregation)
	} else if reverseNestedAgg, ok := a.base.(*elastic.ReverseNestedAggregation); ok {
		reverseNestedAgg.SubAggregation(name, subAggregation)
	} else {
		return fmt.Errorf("unknown aggregation type")
	}

	return nil
}

func nilAggregation(base elastic.Aggregation) *aggregation {
	return &aggregation{children: make(map[string]Aggregation), base: base}
}
