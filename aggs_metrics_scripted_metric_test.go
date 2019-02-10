package aggretastic

import (
	"encoding/json"
	"testing"
	"github.com/olivere/elastic"
)

func TestScriptedMetricAggregation(t *testing.T) {
	agg := NewScriptedMetricAggregation().
		InitScript(elastic.NewScript("state.transactions = []")).
		MapScript(elastic.NewScript("state.transactions.add(doc.type.value == 'sale' ? doc.amount.value : -1 * doc.amount.value)")).
		CombineScript(elastic.NewScript("double profit = 0; for (t in state.transactions) { profit += t } return profit")).
		ReduceScript(elastic.NewScript("double profit = 0; for (a in states) { profit += a } return profit"))

	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"scripted_metric":{"combine_script":{"source":"double profit = 0; for (t in state.transactions) { profit += t } return profit"},"init_script":{"source":"state.transactions = []"},"map_script":{"source":"state.transactions.add(doc.type.value == 'sale' ? doc.amount.value : -1 * doc.amount.value)"},"reduce_script":{"source":"double profit = 0; for (a in states) { profit += a } return profit"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestScriptedMetricAggregationWithParams(t *testing.T) {
	agg := NewScriptedMetricAggregation().
		MapScript(elastic.NewScript("r=0;")).
		ReduceScript(elastic.NewScript("return params.a;")).
		Params(map[string]interface{}{"a": 3})

	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"scripted_metric":{"map_script":{"source":"r=0;"},"params":{"a":3},"reduce_script":{"source":"return params.a;"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestScriptedMetricAggregationWithMeta(t *testing.T) {
	agg := NewScriptedMetricAggregation().
		MapScript(elastic.NewScript("r=0;")).
		ReduceScript(elastic.NewScript("return params.a;")).
		Meta(map[string]interface{}{"foo": "bar"})

	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"meta":{"foo":"bar"},"scripted_metric":{"map_script":{"source":"r=0;"},"reduce_script":{"source":"return params.a;"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
