package aggretastic

import (
	"fmt"
	"github.com/olivere/elastic"
)

// declare an elastic script for every math operation
var (
	addScript      = elastic.NewScript("params.a + params.b")
	subtractScript = elastic.NewScript("params.a - params.b")
	multiplyScript = elastic.NewScript("params.a * params.b")
	divideScript   = elastic.NewScript("params.a / params.b")
	percentScript  = elastic.NewScript("(params.a / params.b) * 100")
	valScript      = elastic.NewScript("params.a")
)

// BucketsPath consists bucket's paths
type BucketsPath map[string]string

// BucketScriptAddAggregation performs math plus operation
func BucketScriptAddAggregation(a, b string) *BucketScriptAggregation {
	return newBucketScriptAggregation(BucketsPath{
		"a": a,
		"b": b,
	}, addScript)
}

// BucketScriptSubtractAggregation performs math minus operation
func BucketScriptSubtractAggregation(a, b string) *BucketScriptAggregation {
	return newBucketScriptAggregation(BucketsPath{
		"a": a,
		"b": b,
	}, subtractScript)
}

// BucketScriptMultiplyAggregation performs math multiply operation
func BucketScriptMultiplyAggregation(a, b string) *BucketScriptAggregation {
	return newBucketScriptAggregation(BucketsPath{
		"a": a,
		"b": b,
	}, multiplyScript)
}

// BucketScriptDivideAggregation performs math division operation
func BucketScriptDivideAggregation(a, b string) *BucketScriptAggregation {
	return newBucketScriptAggregation(BucketsPath{
		"a": a,
		"b": b,
	}, divideScript)
}

// BucketScriptPercentAggregation performs math percent operation
func BucketScriptPercentAggregation(a, b string) *BucketScriptAggregation {
	return newBucketScriptAggregation(BucketsPath{
		"a": a,
		"b": b,
	}, percentScript)
}

// BucketScriptValAggregation performs simple equal operation (just return the value)
func BucketScriptValAggregation(a string) *BucketScriptAggregation {
	return newBucketScriptAggregation(BucketsPath{
		"a": a,
	}, valScript)
}

// BucketScriptNumDivideAggregation performs math divide (between bucket value and number) operation
func BucketScriptNumDivideAggregation(a string, num float64) *BucketScriptAggregation {
	return newBucketScriptAggregation(BucketsPath{
		"a": a,
	}, elastic.NewScript("params.a / "+fmt.Sprintf("%f", num)))
}

// newBucketScriptAggregation is a private function, constructor of elastic.BucketScriptAggregation
func newBucketScriptAggregation(bucketPaths BucketsPath, script *elastic.Script) *BucketScriptAggregation {
	bsa := NewBucketScriptAggregation()

	for k, v := range bucketPaths {
		bsa = bsa.AddBucketsPath(k, v)
	}

	bsa = bsa.Script(script)

	return bsa
}

// RewriteBucketScriptPath allows to rewrite the bucket script path map
func RewriteBucketScriptPath(agg *BucketScriptAggregation, rewriter func(string) (string, bool)) {
	for k, v := range agg.bucketsPathsMap {
		rewrited, deleted := rewriter(v)
		if deleted {
			delete(agg.bucketsPathsMap, k)
		} else {
			agg.bucketsPathsMap[k] = rewrited
		}
	}
}

func IsBucketScriptAggregation(agg Aggregation) (ok bool) {
	_, ok = agg.(*BucketScriptAggregation)
	return
}

func IsBucketSortAggregation(agg Aggregation) (ok bool) {
	_, ok = agg.(*BucketSortAggregation)
	return
}
