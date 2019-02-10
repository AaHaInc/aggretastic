package aggretastic

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotInjectable", func() {
	ch1 := make(chan Aggregation, len(aggMap["NotInjectable"]))
	ch2 := make(chan Aggregation, len(aggMap["NotInjectable"]))
	ch3 := make(chan Aggregation, len(aggMap["NotInjectable"]))
	ch4 := make(chan Aggregation, len(aggMap["NotInjectable"]))
	ch5 := make(chan Aggregation, len(aggMap["NotInjectable"]))
	ch6 := make(chan Aggregation, len(aggMap["NotInjectable"]))
	ch7 := make(chan Aggregation, len(aggMap["NotInjectable"]))
	for name, funct := range aggMap["NotInjectable"] {
		if name == "NewMovFnAggregation" {
			continue
		}
		Describe(name+" Injects", func() {

			ch1 <- getNotInjectable(name, funct)
			It(name+" shouldn't inject aggregations", func() {
				shouldNotInjectTest(ch1)
			})

			ch2 <- getNotInjectable(name, funct)
			It(name+" shouldn't injectX aggregations", func() {
				shouldNotInjectXTest(ch2)
			})
		})

		Describe(name+" Getters", func() {
			ch3 <- getNotInjectable(name, funct)
			It(name+" shouldn't receive subAggregations", func() {
				shouldNotGetAllSubsTest(ch3)
			})

			ch4 <- getNotInjectable(name, funct)
			It(name+" shouldn't select aggregations", func() {
				shouldNotSelectTest(ch4)
			})

			ch5 <- getNotInjectable(name, funct)
			It(name+" shouldn't pop aggregations", func() {
				shouldNotPopTest(ch5)
			})
		})

		Describe(name+" Exports", func() {
			ch6 <- getNotInjectable(name, funct)
			It(" shouldn't export root", func() {
				notInjectableExportTest(ch6)
			})
			ch7 <- getNotInjectable(name, funct)
			It(name+" shouldn't export source", func() {
				notInjectableSourceTest(ch7)
			})

		})
	}
})

func shouldNotInjectTest(ch chan Aggregation) {
	sample := <-ch
	sampleSource, _ := sample.Source()

	injectErr := sample.Inject(NewAvgAggregation(), "any_path")
	Expect(injectErr).Should(HaveOccurred())

	sampleSource2, sourceErr := sample.Source()
	Expect(sourceErr).ShouldNot(HaveOccurred())
	Expect(sampleSource).To(Equal(sampleSource2))
}

func shouldNotInjectXTest(ch chan Aggregation) {
	sample := <-ch
	sampleSource, _ := sample.Source()
	injectErr := sample.InjectX(NewAvgAggregation(), "any_path")
	Expect(injectErr).Should(HaveOccurred())

	sampleSource2, sourceErr := sample.Source()
	Expect(sourceErr).ShouldNot(HaveOccurred())
	Expect(sampleSource).To(Equal(sampleSource2))
}

func shouldNotGetAllSubsTest(ch chan Aggregation) {
	sample := <-ch
	subs := sample.GetAllSubs()
	Expect(subs).To(BeNil())
}

func shouldNotSelectTest(ch chan Aggregation) {
	sample := <-ch
	_ = sample.Inject(NewAvgAggregation(), "any_path")
	sel := sample.Select("any_path")
	Expect(sel).To(BeNil())
}

func shouldNotPopTest(ch chan Aggregation) {
	sample := <-ch
	sampleSource, _ := sample.Source()

	sel := sample.Pop("any_path")
	Expect(sel).To(BeNil())

	sampleSource2, sourceErr := sample.Source()
	Expect(sourceErr).ShouldNot(HaveOccurred())
	Expect(sampleSource).To(Equal(sampleSource2))
}

func notInjectableExportTest(ch chan Aggregation) {
	sample := <-ch
	root := sample.Export()
	Expect(root).To(Equal(sample))
}

func notInjectableSourceTest(ch chan Aggregation) {
	sample := <-ch
	src, err := sample.Source()
	Expect(src).NotTo(BeNil())
	Expect(err).ShouldNot(HaveOccurred())
}

func getNotInjectable(name string, funct interface{}) Aggregation {
	r, _ := Call(funct)
	sample, _ := r[0].Interface().(Aggregation)

	return sample
}
