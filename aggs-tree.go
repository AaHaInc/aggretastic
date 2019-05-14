package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

var (
	ErrNoPath             = fmt.Errorf("no path")
	ErrPathNotSelectable  = fmt.Errorf("path is not selectable")
	ErrAggIsNotInjectable = fmt.Errorf("agg is not injectable")
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
	Inject(subAgg Aggregation, path ...string) (resultPaths [][]string, err error)

	// InjectX sets new subAgg into the map of subAggregations only if it NOT exists already
	InjectX(subAgg Aggregation, path ...string) (resultPaths [][]string, err error)

	// InjectSafe sets new subAgg into the map of subAggregations in the SAFE mode
	InjectSafe(subAgg Aggregation, path ...string) (resultPaths [][]string, err error)

	// Select returns any subAgg by it's path
	Select(path ...string) Aggregation

	// Pop returns a subAgg by it's path and remove it from tree
	Pop(path ...string) Aggregation

	// Export returns the same object in original Agg interface
	Export() elastic.Aggregation

	// ExtractLeafPaths returns paths the leafs
	ExtractLeafPaths() [][]string
}

func IsNilTree(t Aggregation) bool {
	return t == nil || t.Export() == nil
}

type tree struct {
	root            elastic.Aggregation
	subAggregations map[string]Aggregation
}

func nilAggregationTree(root elastic.Aggregation) *tree {
	return &tree{
		root:            root,
		subAggregations: make(map[string]Aggregation),
	}
}

func (a *tree) ExtractLeafPaths() (leafs [][]string) {
	leafs = make([][]string, 0)

	if len(a.subAggregations) == 0 {
		leafs = append(leafs, []string{})
		return
	}

	for leafName, leafAgg := range a.subAggregations {
		for _, leaf := range leafAgg.ExtractLeafPaths() {
			leafs = append(leafs, append([]string{leafName}, leaf...))
		}
	}

	return leafs
}

func (a *tree) Inject(subAggregation Aggregation, path ...string) (resultPaths [][]string, err error) {
	resultPaths = make([][]string, 0)

	if len(path) == 0 {
		err = ErrNoPath
		return
	}

	if len(path) == 1 {
		a.subAggregations[path[0]] = subAggregation
		for _, leaf := range subAggregation.ExtractLeafPaths() {
			resultPaths = append(resultPaths, append(path, leaf...))
		}
		return
	}

	// deeper inject
	cursor := a.Select(path[:len(path)-1]...)
	if IsNilTree(cursor) {
		err = ErrPathNotSelectable
		return
	}

	if _, err = cursor.Inject(subAggregation, path[len(path)-1]); err != nil {
		return
	}

	resultPaths = getIntersectedPaths(a.subAggregations, subAggregation, path)
	return
}

func (a *tree) InjectX(subAggregation Aggregation, path ...string) (resultPaths [][]string, err error) {
	resultPaths = make([][]string, 0)

	if len(path) == 0 {
		err = ErrNoPath
		return
	}

	if alreadyInjected := a.Select(path...); IsNilTree(alreadyInjected) {
		resultPaths, err = a.Inject(subAggregation, path...)
		if err != nil {
			return
		}
	}

	return
}

func (a *tree) InjectSafe(subAggregation Aggregation, path ...string) (resultPaths [][]string, err error) {
	resultPaths = make([][]string, 0)

	if len(path) == 0 {
		err = ErrNoPath
		return
	}

	// extracting the sub tree
	subTree := a.Select(path...)

	if IsNilTree(subTree) {
		return a.Inject(subAggregation, path...)
	}

	for k, subAggDeep := range subAggregation.GetAllSubs() {
		_, injectErr := subTree.InjectSafe(subAggDeep, k)
		if injectErr != nil {
			err = injectErr
			return
		}
	}

	resultPaths = getIntersectedPaths(a.subAggregations, subAggregation, path)

	return
}

func (a *tree) GetAllSubs() map[string]Aggregation {
	return a.subAggregations
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
		return subAgg
	}

	return subAgg.Select(path[1:]...)
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
		return subAgg
	}

	return subAgg.Pop(path[1:]...)
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
func (a *Aggregations) Inject(subAgg Aggregation, path ...string) (resultPaths [][]string, err error) {
	resultPaths = make([][]string, 0)

	if a == nil {
		err = ErrAggIsNotInjectable
		return
	}

	if len(path) == 0 {
		err = ErrNoPath
		return
	}

	name := path[0]

	if len(path) == 1 {
		(*a)[name] = subAgg
		for _, leaf := range subAgg.ExtractLeafPaths() {
			resultPaths = append(resultPaths, append(path, leaf...))
		}

		return
	}

	path = path[1:]
	if _, ok := (*a)[name]; !ok {
		err = ErrAggIsNotInjectable
		return
	}

	resultPaths, err = (*a)[name].Inject(subAgg, path...)
	if err == nil {
		for i := range resultPaths {
			resultPaths[i] = append([]string{name}, resultPaths[i]...)
		}
	}
	return
}

func (a *Aggregations) InjectX(subAgg Aggregation, path ...string) (resultPaths [][]string, err error) {
	resultPaths = make([][]string, 0)

	if a == nil {
		err = ErrAggIsNotInjectable
		return
	}

	if len(path) == 0 {
		err = ErrNoPath
		return
	}

	name := path[0]

	if len(path) == 1 {
		if _, ok := (*a)[name]; !ok {
			(*a)[name] = subAgg
		}
		// @tody return path

		return
	}

	path = path[1:]
	if _, ok := (*a)[name]; !ok {
		err = ErrAggIsNotInjectable
		return
	}

	return (*a)[name].InjectX(subAgg, path...)
}

func (a *Aggregations) InjectSafe(subAgg Aggregation, path ...string) (resultPaths [][]string, err error) {
	resultPaths = make([][]string, 0)

	if len(path) == 0 {
		err = ErrNoPath
		return
	}

	name := path[0]

	if len(path) == 1 {
		if _, ok := (*a)[name]; !ok {
			(*a)[name] = subAgg
			resultPaths = subAgg.ExtractLeafPaths()
			for i := range resultPaths {
				// prepend name
				resultPaths[i] = append([]string{name}, resultPaths[i]...)
			}
			return
		}

		// @todo
		log.Println("warning! maybe unexpected behaviour. Edge case, need handling")
		return
	}

	path = path[1:]
	if _, ok := (*a)[name]; !ok {
		err = ErrAggIsNotInjectable
		return
	}

	if resultPaths, err = (*a)[name].InjectSafe(subAgg, path...); err == nil {
		// prepend name to result paths
		for i := range resultPaths {
			resultPaths[i] = append([]string{name}, resultPaths[i]...)
		}
	}

	return
}

//
// helpers
//

// getIntersectedPaths returns leafs' paths that intersects of agg and its subAgg
func getIntersectedPaths(rootAggs map[string]Aggregation, injectingAgg Aggregation, injectPath []string) [][]string {
	paths := make([][]string, 0)

	rootLeafs := make([][]string, 0)
	if len(rootAggs) == 0 {
		rootLeafs = append(rootLeafs, []string{})
	} else {
		for leafName, leafAgg := range rootAggs {
			for _, leaf := range leafAgg.ExtractLeafPaths() {
				rootLeafs = append(rootLeafs, append([]string{leafName}, leaf...))
			}
		}
	}

	givenLeafPaths := injectingAgg.ExtractLeafPaths()
	for i := range givenLeafPaths {
		givenLeafPaths[i] = append(injectPath, givenLeafPaths[i]...)
	}

	for _, resultLeaf := range rootLeafs {
		for _, givenLeaf := range givenLeafPaths {
			if pathIsLeafOf(givenLeaf, resultLeaf) {
				paths = append(paths, resultLeaf)
			}
		}
	}

	return paths
}

// pathIsLeafOf checks if childPath is a finite leaf of parentPath
func pathIsLeafOf(childPath, parentPath []string) bool {
	if len(parentPath) < len(childPath) {
		return false
	}

	if len(parentPath) > len(childPath) {
		parentPath = parentPath[len(parentPath)-len(childPath):]
	}

	var diff bool
	for i := range childPath {
		if childPath[i] != parentPath[i] {
			diff = true
		}
	}
	return !diff
}
