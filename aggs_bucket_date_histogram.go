package aggretastic

import "github.com/olivere/elastic"

type DateHistogramAggregation struct {
	*aggregation

	field   string
	script  *elastic.Script
	missing interface{}
	meta    map[string]interface{}

	interval          string
	order             string
	orderAsc          bool
	minDocCount       *int64
	extendedBoundsMin interface{}
	extendedBoundsMax interface{}
	timeZone          string
	format            string
	offset            string
}

func NewDateHistogramAggregation() *DateHistogramAggregation {
	a := &DateHistogramAggregation{}
	a.aggregation = nilAggregation()

	return a
}

func (a *DateHistogramAggregation) Field(field string) *DateHistogramAggregation {
	a.field = field
	return a
}

func (a *DateHistogramAggregation) Script(script *elastic.Script) *DateHistogramAggregation {
	a.script = script
	return a
}

func (a *DateHistogramAggregation) Missing(missing interface{}) *DateHistogramAggregation {
	a.missing = missing
	return a
}

func (a *DateHistogramAggregation) SubAggregation(name string, subAggregation Aggregation) *DateHistogramAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

func (a *DateHistogramAggregation) Meta(metaData map[string]interface{}) *DateHistogramAggregation {
	a.meta = metaData
	return a
}

func (a *DateHistogramAggregation) Interval(interval string) *DateHistogramAggregation {
	a.interval = interval
	return a
}

func (a *DateHistogramAggregation) Order(order string, asc bool) *DateHistogramAggregation {
	a.order = order
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) OrderByCount(asc bool) *DateHistogramAggregation {
	// "order" : { "_count" : "asc" }
	a.order = "_count"
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) OrderByCountAsc() *DateHistogramAggregation {
	return a.OrderByCount(true)
}

func (a *DateHistogramAggregation) OrderByCountDesc() *DateHistogramAggregation {
	return a.OrderByCount(false)
}

func (a *DateHistogramAggregation) OrderByKey(asc bool) *DateHistogramAggregation {
	// "order" : { "_key" : "asc" }
	a.order = "_key"
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) OrderByKeyAsc() *DateHistogramAggregation {
	return a.OrderByKey(true)
}

func (a *DateHistogramAggregation) OrderByKeyDesc() *DateHistogramAggregation {
	return a.OrderByKey(false)
}

func (a *DateHistogramAggregation) OrderByAggregation(aggName string, asc bool) *DateHistogramAggregation {
	a.order = aggName
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) OrderByAggregationAndMetric(aggName, metric string, asc bool) *DateHistogramAggregation {
	a.order = aggName + "." + metric
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) MinDocCount(minDocCount int64) *DateHistogramAggregation {
	a.minDocCount = &minDocCount
	return a
}

func (a *DateHistogramAggregation) TimeZone(timeZone string) *DateHistogramAggregation {
	a.timeZone = timeZone
	return a
}

func (a *DateHistogramAggregation) Format(format string) *DateHistogramAggregation {
	a.format = format
	return a
}

func (a *DateHistogramAggregation) Offset(offset string) *DateHistogramAggregation {
	a.offset = offset
	return a
}

func (a *DateHistogramAggregation) ExtendedBounds(min, max interface{}) *DateHistogramAggregation {
	a.extendedBoundsMin = min
	a.extendedBoundsMax = max
	return a
}

func (a *DateHistogramAggregation) ExtendedBoundsMin(min interface{}) *DateHistogramAggregation {
	a.extendedBoundsMin = min
	return a
}

func (a *DateHistogramAggregation) ExtendedBoundsMax(max interface{}) *DateHistogramAggregation {
	a.extendedBoundsMax = max
	return a
}

func (a *DateHistogramAggregation) Source() (interface{}, error) {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["date_histogram"] = opts

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
	if a.missing != nil {
		opts["missing"] = a.missing
	}

	opts["interval"] = a.interval
	if a.minDocCount != nil {
		opts["min_doc_count"] = *a.minDocCount
	}
	if a.order != "" {
		o := make(map[string]interface{})
		if a.orderAsc {
			o[a.order] = "asc"
		} else {
			o[a.order] = "desc"
		}
		opts["order"] = o
	}
	if a.timeZone != "" {
		opts["time_zone"] = a.timeZone
	}
	if a.offset != "" {
		opts["offset"] = a.offset
	}
	if a.format != "" {
		opts["format"] = a.format
	}
	if a.extendedBoundsMin != nil || a.extendedBoundsMax != nil {
		bounds := make(map[string]interface{})
		if a.extendedBoundsMin != nil {
			bounds["min"] = a.extendedBoundsMin
		}
		if a.extendedBoundsMax != nil {
			bounds["max"] = a.extendedBoundsMax
		}
		opts["extended_bounds"] = bounds
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
