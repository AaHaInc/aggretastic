// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package aggretastic

import (
	"encoding/json"
	"testing"
	"github.com/olivere/elastic"
)

func TestBucketSelectorAggregation(t *testing.T) {
	agg := NewBucketSelectorAggregation().
		AddBucketsPath("totalSales", "total_sales").
		Script(elastic.NewScript("totalSales >= 1000"))
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"bucket_selector":{"buckets_path":{"totalSales":"total_sales"},"script":{"source":"totalSales \u003e= 1000"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
