package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic"
)

type Aggregation interface {
	elastic.Aggregation

	Select(path ...string) Aggregation
	Inject(agg Aggregation, path ...string) error
	InjectX(agg Aggregation, path ...string) error
	WrapBy(wrapper Aggregation, name string)
	InjectWrapper(wrapper Aggregation, path ...string) error
	GetAllSubs() map[string]Aggregation

	getKey() string
	setKey(key string)

	setParent(parent Aggregation)
	getParent() Aggregation

	setSubAggregation(child Aggregation, childName string)
	getSubAggregation(childName string) Aggregation
	deleteSubAggregation(childName string)
	getSubAggregations() map[string]Aggregation
}

type aggregation struct {
	key             string
	parent          Aggregation
	subAggregations map[string]Aggregation
}

func (a *aggregation) Source() (interface{}, error) {
	return nil, fmt.Errorf("nil aggregation")
}

func (a *aggregation) Select(path ...string) Aggregation {
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

func (a *aggregation) Inject(subAggregation Aggregation, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("no path")
	}

	if len(path) == 1 {
		// injection means setting key, parent, child
		subAggregation.setParent(a)
		subAggregation.setKey(path[0])
		a.setSubAggregation(subAggregation, path[0])
		return nil
	}

	// use select+recursive inject to handle long path

	cursor := a.Select(path[:len(path)-1]...)
	if cursor == nil {
		return fmt.Errorf("path not selectable")
	}

	return cursor.Inject(subAggregation, path[len(path)-1])
}

func (a *aggregation) InjectX(subAggregation Aggregation, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("no path")
	}

	if alreadyInjected := a.Select(path...); alreadyInjected == nil {
		return a.Inject(subAggregation, path...)
	}

	return nil
}

func (a *aggregation) WrapBy(wrapper Aggregation, name string) {
	parentKey := a.parent.getKey()
	a.parent.setSubAggregation(wrapper, name)
	this := a.parent.getSubAggregation(a.key) // get `a` object, but in proper type (not *aggregation, but *SomeAggregation)
	wrapper.setKey(a.key)
	wrapper.setSubAggregation(this, a.key)
	a.parent.deleteSubAggregation(a.key)
	a.parent = wrapper
	a.parent.setKey(parentKey)
}

func (a *aggregation) InjectWrapper(wrapper Aggregation, path ...string) error {
	if len(path) == 0 {
		// do nothing
		return fmt.Errorf("nil path")
	}
	if len(path) == 1 {
		a.WrapBy(wrapper, path[0])
		return nil
	}

	n := len(path) - 1
	parent := a.Select(path[:n]...)
	if parent == nil {
		return fmt.Errorf("not selectable")
	}

	wrapperKey := path[n]
	for childKey, child := range parent.getSubAggregations() {
		wrapper.setSubAggregation(child, childKey)
		parent.deleteSubAggregation(childKey)
	}

	parent.setSubAggregation(wrapper, wrapperKey)
	parent.setKey(wrapperKey)

	return nil
}

func (a *aggregation) GetAllSubs() map[string]Aggregation {
	return a.subAggregations
}

// helper util functions

func (a *aggregation) setParent(parent Aggregation) {
	a.parent = parent
}

func (a *aggregation) getParent() Aggregation {
	return a.parent
}

func (a *aggregation) setSubAggregation(child Aggregation, childName string) {
	a.subAggregations[childName] = child
}

func (a *aggregation) getSubAggregation(childName string) Aggregation {
	return a.subAggregations[childName]
}

func (a *aggregation) getSubAggregations() map[string]Aggregation {
	return a.subAggregations
}

func (a *aggregation) deleteSubAggregation(childName string) {
	delete(a.subAggregations, childName)
}

func (a *aggregation) setKey(key string) {
	a.key = key
}
func (a *aggregation) getKey() string {
	return a.key
}

func nilAggregation() *aggregation {
	return &aggregation{subAggregations: make(map[string]Aggregation)}
}
