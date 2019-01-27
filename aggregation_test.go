package aggretastic

import (
	"encoding/json"
	"github.com/olivere/elastic"
	"log"
	"testing"
)

func TestC1(t *testing.T) {
	sessions := NewFilterAggregation()
	sessions.Filter(elastic.NewTermQuery("sessions", true))

	experiments := NewFilterAggregation()
	experiments.Filter(elastic.NewTermQuery("experiments", true))
	experiments.Inject(sessions, "sessions")

	toExperiments := NewNestedAggregation()
	toExperiments.Path("experiments")
	toExperiments.Inject(experiments, "experiments.name")

	// NESTED
	// 	-> experiments.name:FILTER
	//		-> sessions:FILTER

	toExperiments.Select("experiments.name", "sessions").WrapBy(NewReverseNestedAggregation(), "back")

	//render(selected)
	render(toExperiments)
}

func TestC2(t *testing.T) {
	sessions := NewFilterAggregation()
	sessions.Filter(elastic.NewTermQuery("sessions", true))
	hits := NewFilterAggregation()
	hits.Filter(elastic.NewTermQuery("hits", true))

	experiments := NewFilterAggregation()
	experiments.Filter(elastic.NewTermQuery("experiments", true))
	experiments.Inject(sessions, "sessions")
	experiments.Inject(hits, "hits")

	toExperiments := NewNestedAggregation()
	toExperiments.Path("experiments")
	toExperiments.Inject(experiments, "experiments.name")

	// NESTED
	// 	-> experiments.name:FILTER
	//		-> sessions:FILTER

	toExperiments.InjectWrapper(NewReverseNestedAggregation(), "experiments.name", "back")

	//render(selected)
	render(toExperiments)
}

/*
func TestB(t *testing.T) {
	x := NewChildrenAggregation().Type("something")
	ww := NewFilterAggregation().Filter(elastic.NewTermQuery("a", true))
	render(x)
	render(ww)

	x.WrapBy(ww, "a")
	render(x)
	render(ww)
}

func TestA(t *testing.T) {
	x := NewChildrenAggregation().Type("something")
	ww := NewFilterAggregation().Filter(elastic.NewTermQuery("a", true))
	ww2 := NewFilterAggregation().Filter(elastic.NewTermQuery("b", false))
	_ = ww2
	render(x)
	render(ww)

	x.Inject(ww, "a")
	x.Inject(ww2, "a", "b")
	//x.Inject(ww2, "a", "b")
	render(x)
}
*/

func render(x elastic.Aggregation) {
	s, _ := x.Source()
	j, _ := json.Marshal(s)
	log.Println(string(j))
}
