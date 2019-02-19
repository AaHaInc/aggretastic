package aggretastic

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func makeInjection_not(agg, inj Aggregation, x bool, path ...string) {
	var injectErr error
	if x {
		injectErr = agg.InjectX(inj, path...)
	} else {
		injectErr = agg.Inject(inj, path...)
	}
	Expect(injectErr).Should(HaveOccurred())
}

func testNotInjection(maxAgg, avgAgg Aggregation, maxAggSource, avgAggSource interface{}, x bool) {
	makeInjection_not(maxAgg, avgAgg, x, "any_path")
	equal := compareAggregationToSource(maxAgg, maxAggSource)
	Expect(equal).To(Equal(true))

	maxAggSource, _ = maxAgg.Source()
	makeInjection_not(maxAgg, avgAgg, x, "any_path", "other_path")
	equal = compareAggregationToSource(maxAgg, maxAggSource)
	Expect(equal).To(Equal(true))
}

var _ = Describe("Injectable", func() {
	var maxAgg *AvgBucketAggregation
	var maxAggSource interface{}

	var avgAgg *AvgAggregation
	var avgAggSource interface{}

	BeforeEach(func() {
		maxAgg = NewAvgBucketAggregation()
		maxAggSource, _ = maxAgg.Source()

		avgAgg = NewAvgAggregation()
		avgAggSource, _ = avgAgg.Source()
	})
	Describe(" Injects", func() {

		It(" should not inject aggregations", func() {
			testNotInjection(maxAgg, avgAgg, maxAggSource, avgAggSource, false)
		})

		It(" should not injectX aggregations", func() {
			testNotInjection(maxAgg, avgAgg, maxAggSource, avgAggSource, true)
		})
	})

	Describe(" Getters", func() {
		It(" should not select aggregations", func() {
			sel := maxAgg.Select("any_path")
			Expect(sel).To(BeNil())
		})

		It(" should not pop aggregations", func() {
			refSrc, _ := maxAgg.Source()
			pop := maxAgg.Pop("any_path")
			Expect(pop).To(BeNil())
			Expect(maxAgg.Source()).To(Equal(refSrc))
		})
	})

	Describe(" Exports", func() {
		It(" should export root", func() {
			root := maxAgg.Export()
			Expect(root).To(Equal(maxAgg))
		})
		It(" should export source", func() {
			src, err := maxAgg.Source()
			Expect(src).NotTo(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

})
