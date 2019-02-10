package aggretastic

import (
	"errors"
	"github.com/olivere/elastic"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

func Call(function interface{}, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(function)
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

var _ = Describe("Injectable", func() {
	ch1 := make(chan Aggregation, len(aggMap["Injectable"]))
	ch2 := make(chan Aggregation, len(aggMap["Injectable"]))
	ch3 := make(chan Aggregation, len(aggMap["Injectable"]))
	ch4 := make(chan Aggregation, len(aggMap["Injectable"]))
	ch5 := make(chan Aggregation, len(aggMap["Injectable"]))
	ch6 := make(chan Aggregation, len(aggMap["Injectable"]))
	ch7 := make(chan Aggregation, len(aggMap["Injectable"]))
	for name, funct := range aggMap["Injectable"] {
		Describe(name + " Injects", func() {

			ch1 <- getSample(name, funct)
			It(name + " should inject aggregations", func() {
				injectTest(ch1)
			})

			ch2 <- getSample(name, funct)
			It(name + " should injectX aggregations", func() {
				injectXTest(ch2)
			})
		})

		Describe(name + " Getters", func() {
			ch3 <- getSample(name, funct)
			It(name + " should receive subAggregations", func() {
				getAllSubsTest(ch3)
			})

			ch4 <- getSample(name, funct)
			It(name + " should select aggregations", func() {
				selectTest(ch4)
			})

			ch5 <- getSample(name, funct)
			It(name + " should pop aggregations", func() {
				popTest(ch5)
			})
		})

		Describe(name + " Exports", func() {
			ch6 <- getSample(name, funct)
			It(" should export root", func() {
				exportTest(ch6)
			})
			ch7 <- getSample(name, funct)
			It(name + " should export source", func() {
				sourceTest(ch7)
			})

		})
	}
})

func injectTest(ch chan Aggregation) {
	sample := <-ch
	sampleSource, _ := sample.Source()

	injectErr := sample.Inject(NewAvgAggregation(), "any_path")
	Expect(injectErr).ShouldNot(HaveOccurred())

	sampleSource2, sourceErr := sample.Source()
	Expect(sourceErr).ShouldNot(HaveOccurred())

	Expect(sampleSource).ToNot(Equal(sampleSource2))
}

func injectXTest(ch chan Aggregation) {
	sample := <-ch
	sampleSource, _ := sample.Source()
	injectErr := sample.InjectX(NewAvgAggregation(), "any_path")
	Expect(injectErr).ShouldNot(HaveOccurred())

	sampleSource2, sourceErr := sample.Source()
	Expect(sourceErr).ShouldNot(HaveOccurred())
	Expect(sampleSource).ToNot(Equal(sampleSource2))

	injectXErr := sample.InjectX(NewAvgAggregation(), "any_path")
	Expect(injectXErr).ShouldNot(HaveOccurred())

	sampleSource3, sourceErr := sample.Source()
	Expect(sourceErr).ShouldNot(HaveOccurred())
	Expect(sampleSource2).To(Equal(sampleSource3))
}

func getAllSubsTest(ch chan Aggregation) {
	sample := <-ch
	subs := sample.GetAllSubs()
	Expect(subs).NotTo(BeNil())
}

func selectTest(ch chan Aggregation) {
	sample := <-ch
	_ = sample.Inject(NewAvgAggregation(), "any_path")
	sel := sample.Select("any_path")
	Expect(sel).NotTo(BeNil())
}

func popTest(ch chan Aggregation) {
	sample := <-ch
	sampleSource, _ := sample.Source()

	_ = sample.Inject(NewAvgAggregation(), "any_path")
	sampleSource, _ = sample.Source()
	sel := sample.Pop("any_path")
	Expect(sel).NotTo(BeNil())

	sampleSource2, sourceErr := sample.Source()
	Expect(sourceErr).ShouldNot(HaveOccurred())
	Expect(sampleSource).NotTo(Equal(sampleSource2))
}

func exportTest(ch chan Aggregation) {
	sample := <-ch
	root := sample.Export()
	Expect(root).To(Equal(sample))
}

func sourceTest(ch chan Aggregation) {
	sample := <-ch
	src, err := sample.Source()
	Expect(src).NotTo(BeNil())
	Expect(err).ShouldNot(HaveOccurred())
}

func getSample(name string, funct interface{}) Aggregation {
	r, _ := Call(funct)
	sample, _ := r[0].Interface().(Aggregation)
	if name == "NewFilterAggregation" {
		sample.(*FilterAggregation).Filter(elastic.NewTermQuery("user", "olivere"))
	}
	return sample
}
