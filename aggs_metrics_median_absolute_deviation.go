package aggretastic

import "github.com/olivere/elastic/v7"

// MedianAbsoluteDeviationAggregation is a single-value aggregation that
// approximates the median absolute deviation of its search results.
//
// See: https://www.elastic.co/guide/en/elasticsearch/reference/7.2/search-aggregations-metrics-median-absolute-deviation-aggregation.html
type MedianAbsoluteDeviationAggregation struct {
	*tree

	field       string
	script      *elastic.Script
	compression int64
	missing     interface{}

	meta map[string]interface{}
}

func NewMedianAbsoluteDeviationAggregation() *MedianAbsoluteDeviationAggregation {
	a := &MedianAbsoluteDeviationAggregation{
		compression: -1,
	}
	a.tree = nilAggregationTree(a)

	return a
}

func (a *MedianAbsoluteDeviationAggregation) Field(field string) *MedianAbsoluteDeviationAggregation {
	a.field = field
	return a
}

func (a *MedianAbsoluteDeviationAggregation) Script(script *elastic.Script) *MedianAbsoluteDeviationAggregation {
	a.script = script
	return a
}

func (a *MedianAbsoluteDeviationAggregation) Compression(compression int64) *MedianAbsoluteDeviationAggregation {
	a.compression = compression
	return a
}

func (a *MedianAbsoluteDeviationAggregation) Missing(missing interface{}) *MedianAbsoluteDeviationAggregation {
	a.missing = missing
	return a
}

func (a *MedianAbsoluteDeviationAggregation) SubAggregation(name string, subAggregation Aggregation) *MedianAbsoluteDeviationAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *MedianAbsoluteDeviationAggregation) Meta(metaData map[string]interface{}) *MedianAbsoluteDeviationAggregation {
	a.meta = metaData
	return a
}

func (a *MedianAbsoluteDeviationAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//      "review_variability" : { "median_absolute_deviation" : { "field" : "rating", "compression": 100 } }
	//    }
	//	}
	// This method returns only the { "median_absolute_deviation" : { "field" : "ration", "compression": 100 } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["median_absolute_deviation"] = opts

	// ValuesSourceAggregationBuilder
	if a.field != "" {
		opts["field"] = a.field
	}
	if a.script != nil {
		src, err := a.script.Source()
		if err != nil {
			return nil, err
		}
		opts["script"] = src
	}
	if a.compression > 0 {
		opts["compression"] = a.compression
	}
	if a.missing != nil {
		opts["missing"] = a.missing
	}

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}
