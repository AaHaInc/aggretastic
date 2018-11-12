package aggretastic

import "github.com/olivere/elastic"

// PercentilesAggregation is a multi-value metrics aggregation
// that calculates one or more percentiles over numeric values
// extracted from the aggregated documents. These values can
// be extracted either from specific numeric fields in the documents,
// or be generated by a provided script.
//
// See: https://www.elastic.co/guide/en/elasticsearch/reference/6.2/search-aggregations-metrics-percentile-aggregation.html
type PercentilesAggregation struct {
	*aggregation

	field       string
	script      *elastic.Script
	format      string
	meta        map[string]interface{}
	percentiles []float64
	compression *float64
	estimator   string
}

func NewPercentilesAggregation() *PercentilesAggregation {
	a := &PercentilesAggregation{percentiles: make([]float64, 0)}
	a.aggregation = nilAggregation()

	return a
}

func (a *PercentilesAggregation) Field(field string) *PercentilesAggregation {
	a.field = field
	return a
}

func (a *PercentilesAggregation) Script(script *elastic.Script) *PercentilesAggregation {
	a.script = script
	return a
}

func (a *PercentilesAggregation) Format(format string) *PercentilesAggregation {
	a.format = format
	return a
}

func (a *PercentilesAggregation) SubAggregation(name string, subAggregation Aggregation) *PercentilesAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *PercentilesAggregation) Meta(metaData map[string]interface{}) *PercentilesAggregation {
	a.meta = metaData
	return a
}

func (a *PercentilesAggregation) Percentiles(percentiles ...float64) *PercentilesAggregation {
	a.percentiles = append(a.percentiles, percentiles...)
	return a
}

func (a *PercentilesAggregation) Compression(compression float64) *PercentilesAggregation {
	a.compression = &compression
	return a
}

func (a *PercentilesAggregation) Estimator(estimator string) *PercentilesAggregation {
	a.estimator = estimator
	return a
}

func (a *PercentilesAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//      "load_time_outlier" : {
	//           "percentiles" : {
	//               "field" : "load_time"
	//           }
	//       }
	//    }
	//	}
	// This method returns only the
	//   { "percentiles" : { "field" : "load_time" } }
	// part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["percentiles"] = opts

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
	if a.format != "" {
		opts["format"] = a.format
	}
	if len(a.percentiles) > 0 {
		opts["percents"] = a.percentiles
	}
	if a.compression != nil {
		opts["compression"] = *a.compression
	}
	if a.estimator != "" {
		opts["estimator"] = a.estimator
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
