package aggretastic

// MultiTermsAggregation is very similar to the terms aggregation,
// however in most cases it will be slower than the terms aggregation and will consume more memory.
//
// See: https://www.elastic.co/guide/en/elasticsearch//reference/master/search-aggregations-bucket-multi-terms-aggregation.html
type MultiTermsAggregation struct {
	*tree

	terms []*MultiTermsField
	meta  map[string]interface{}

	size                  *int
	shardSize             *int
	minDocCount           *int
	shardMinDocCount      *int
	collectionMode        string
	showTermDocCountError *bool
	order                 []TermsOrder
}

func NewMultiTermsAggregation() *MultiTermsAggregation {
	a := &MultiTermsAggregation{
		terms: make([]*MultiTermsField, 0),
	}
	a.tree = nilAggregationTree(a)

	return a
}

func (a *MultiTermsAggregation) Term(field string, missingArg ...interface{}) *MultiTermsAggregation {
	term := &MultiTermsField{field: field}
	if len(missingArg) > 0 {
		term.missing = missingArg[0]
	}

	a.terms = append(a.terms, term)
	return a
}

func (a *MultiTermsAggregation) SubAggregation(name string, subAggregation Aggregation) *MultiTermsAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *MultiTermsAggregation) Meta(metaData map[string]interface{}) *MultiTermsAggregation {
	a.meta = metaData
	return a
}

func (a *MultiTermsAggregation) Size(size int) *MultiTermsAggregation {
	a.size = &size
	return a
}

func (a *MultiTermsAggregation) ShardSize(shardSize int) *MultiTermsAggregation {
	a.shardSize = &shardSize
	return a
}

func (a *MultiTermsAggregation) MinDocCount(minDocCount int) *MultiTermsAggregation {
	a.minDocCount = &minDocCount
	return a
}

func (a *MultiTermsAggregation) ShardMinDocCount(shardMinDocCount int) *MultiTermsAggregation {
	a.shardMinDocCount = &shardMinDocCount
	return a
}

func (a *MultiTermsAggregation) Order(order string, asc bool) *MultiTermsAggregation {
	a.order = append(a.order, TermsOrder{Field: order, Ascending: asc})
	return a
}

func (a *MultiTermsAggregation) OrderByCount(asc bool) *MultiTermsAggregation {
	// "order" : { "_count" : "asc" }
	a.order = append(a.order, TermsOrder{Field: "_count", Ascending: asc})
	return a
}

func (a *MultiTermsAggregation) OrderByCountAsc() *MultiTermsAggregation {
	return a.OrderByCount(true)
}

func (a *MultiTermsAggregation) OrderByCountDesc() *MultiTermsAggregation {
	return a.OrderByCount(false)
}

func (a *MultiTermsAggregation) OrderByTerm(asc bool) *MultiTermsAggregation {
	// "order" : { "_term" : "asc" }
	a.order = append(a.order, TermsOrder{Field: "_term", Ascending: asc})
	return a
}

func (a *MultiTermsAggregation) OrderByTermAsc() *MultiTermsAggregation {
	return a.OrderByTerm(true)
}

func (a *MultiTermsAggregation) OrderByTermDesc() *MultiTermsAggregation {
	return a.OrderByTerm(false)
}

// OrderByAggregation creates a bucket ordering strategy which sorts buckets
// based on a single-valued calc get.
func (a *MultiTermsAggregation) OrderByAggregation(aggName string, asc bool) *MultiTermsAggregation {
	// {
	//     "aggs" : {
	//         "genders_age" : {
	//             "multi_terms" : {
	//				   "terms" : [{"field" : "gender"}, {"field":"age_group"}],
	//                 "order" : { "avg_height" : "desc" }
	//             },
	//             "aggs" : {
	//                 "avg_height" : { "avg" : { "field" : "height" } }
	//             }
	//         }
	//     }
	// }
	a.order = append(a.order, TermsOrder{Field: aggName, Ascending: asc})
	return a
}

// OrderByAggregationAndMetric creates a bucket ordering strategy which
// sorts buckets based on a multi-valued calc get.
func (a *MultiTermsAggregation) OrderByAggregationAndMetric(aggName, metric string, asc bool) *MultiTermsAggregation {
	// {
	//     "aggs" : {
	//         "genders_age" : {
	//             "multi_terms" : {
	//				   "terms" : [{"field" : "gender"}, {"field":"age_group"}],
	//                 "order" : { "height_stats.avg" : "desc" }
	//             },
	//             "aggs" : {
	//                 "height_stats" : { "stats" : { "field" : "height" } }
	//             }
	//         }
	//     }
	// }
	a.order = append(a.order, TermsOrder{Field: aggName + "." + metric, Ascending: asc})
	return a
}

// CollectionMode can be depth_first or breadth_first as of 1.4.0.
func (a *MultiTermsAggregation) CollectionMode(collectionMode string) *MultiTermsAggregation {
	a.collectionMode = collectionMode
	return a
}

func (a *MultiTermsAggregation) ShowTermDocCountError(showTermDocCountError bool) *MultiTermsAggregation {
	a.showTermDocCountError = &showTermDocCountError
	return a
}

func (a *MultiTermsAggregation) Source() (interface{}, error) {
	// Example:
	//	{
	//    "aggs" : {
	//      "genders_age" : {
	//        "multi_terms" : [{"field" : "gender"}, {"field" : "age_group"}]
	//      }
	//    }
	//	}
	// This method returns only the { "multi_terms" : ... } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["multi_terms"] = opts

	// TermsBuilder
	terms := make([]interface{}, 0)
	for _, term := range a.terms {
		termSrc, err := term.Source()
		if err != nil {
			return nil, err
		}

		terms = append(terms, termSrc)
	}
	opts["terms"] = terms

	// MultiTermsBuilder
	if a.size != nil && *a.size >= 0 {
		opts["size"] = *a.size
	}
	if a.shardSize != nil && *a.shardSize >= 0 {
		opts["shard_size"] = *a.shardSize
	}
	if a.minDocCount != nil && *a.minDocCount >= 0 {
		opts["min_doc_count"] = *a.minDocCount
	}
	if a.shardMinDocCount != nil && *a.shardMinDocCount >= 0 {
		opts["shard_min_doc_count"] = *a.shardMinDocCount
	}
	if a.showTermDocCountError != nil {
		opts["show_term_doc_count_error"] = *a.showTermDocCountError
	}
	if a.collectionMode != "" {
		opts["collect_mode"] = a.collectionMode
	}
	if len(a.order) > 0 {
		var orderSlice []interface{}
		for _, order := range a.order {
			src, err := order.Source()
			if err != nil {
				return nil, err
			}
			orderSlice = append(orderSlice, src)
		}
		opts["order"] = orderSlice
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

type MultiTermsField struct {
	field   string
	missing interface{}
}

// Source returns serializable JSON of the MultiTermsField.
func (f *MultiTermsField) Source() (interface{}, error) {
	source := make(map[string]interface{})
	source["field"] = f.field
	if f.missing != nil {
		source["missing"] = f.missing
	}

	return source, nil
}
