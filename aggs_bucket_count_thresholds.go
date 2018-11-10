package aggretastic

// BucketCountThresholds is used in e.g. terms and significant text aggregations.
type BucketCountThresholds struct {
	MinDocCount      *int64
	ShardMinDocCount *int64
	RequiredSize     *int
	ShardSize        *int
}
