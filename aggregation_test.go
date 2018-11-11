package aggretastic

import (
	"encoding/json"
	"github.com/olivere/elastic"
	"log"
	"testing"
)

func TestC(t *testing.T) {
	sessions := NewFilterAggregation()
	sessions.Base().Filter(elastic.NewTermQuery("sessions", true))

	experiments := NewFilterAggregation()
	experiments.Base().Filter(elastic.NewTermQuery("experiments", true))
	experiments.Inject(sessions, "sessions")

	toExperiments := NewNestedAggregation()
	toExperiments.Base().Path("experiments")
	toExperiments.Inject(experiments, "experiments.name")


	log.Printf("1 %#v", toExperiments)
	log.Printf("2 %#v", toExperiments.aggregation)
	log.Printf("3 %#v", toExperiments.aggregation.children["experiments.name"])
	log.Printf("4 %#v", toExperiments.aggregation.children["experiments.name"].getKey())
	log.Printf("4 %#v", toExperiments.aggregation.children["experiments.name"].getChild("sessions").getKey())
	log.Println("")
	//log.Printf("EXPERIMENTS %#v", toExperiments.aggregation.children["experiments"].getKey())

	// NESTED
	// 	-> experiments.name:FILTER
	//		-> sessions:FILTER


	log.Println("===========================")
	render(toExperiments)
	log.Println("---------------------------")
	toExperiments.Select("experiments.name", "sessions").WrapBy(NewReverseNestedAggregation(), "back")

	log.Printf("1 %#v", toExperiments)
	log.Printf("2 %#v", toExperiments.aggregation)
	log.Printf("3 %#v", toExperiments.aggregation.children["experiments.name"])
	log.Printf("4 %#v", toExperiments.aggregation.children["experiments.name"].getKey())
	//log.Printf("4 %#v", toExperiments.aggregation.children["experiments.name"].getChild("back").getKey())
	log.Printf("5 %#v", toExperiments.aggregation.children["experiments.name"].getChildren())
	log.Println("")

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
