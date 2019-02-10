// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package aggretastic

import (
	"encoding/json"
	"testing"
	"github.com/olivere/elastic"
)

func TestMovFnAggregation(t *testing.T) {
	agg := NewMovFnAggregation(
		"the_sum",
		elastic.NewScript("MovingFunctions.min(values)"),
		10,
	)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"moving_fn":{"buckets_path":"the_sum","script":{"source":"MovingFunctions.min(values)"},"window":10}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
