package aggretastic

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func makeInjection(agg, inj Aggregation, x bool, path ...string) {
	var injectErr error
	if x {
		injectErr = agg.InjectX(inj, path...)
	} else {
		injectErr = agg.Inject(inj, path...)
	}
	Expect(injectErr).ShouldNot(HaveOccurred())
}

func compareAggregationToSource(firstAgg Aggregation, secondAggSrc interface{}) bool {
	firstAggSrc, err := firstAgg.Source()
	Expect(err).ShouldNot(HaveOccurred())

	jsonOne, _ := json.Marshal(firstAggSrc)
	jsonTwo, _ := json.Marshal(secondAggSrc)
	return bytes.Compare(jsonOne, jsonTwo) == 0
}

func testInjection(maxAgg, avgAgg Aggregation, maxAggSource, avgAggSource interface{}, x bool) {
	makeInjection(maxAgg, avgAgg, x, "any_path")

	equal := compareAggregationToSource(maxAgg, maxAggSource)
	Expect(equal).To(Equal(false))

	equal = compareAggregationToSource(maxAgg.Select("any_path"), avgAggSource)
	Expect(equal).To(Equal(true))

	maxAggSource, _ = maxAgg.Source()

	avgAgg = NewAvgAggregation()
	avgAggSource, _ = avgAgg.Source()

	makeInjection(maxAgg, avgAgg, x, "any_path", "other_path")

	equal = compareAggregationToSource(maxAgg, maxAggSource)
	Expect(equal).To(Equal(false))

	equal = compareAggregationToSource(maxAgg.Select("any_path", "other_path"), avgAggSource)
	Expect(equal).To(Equal(true))

	avgAgg = NewAvgAggregation()
	avgAggSource, _ = avgAgg.Source()

	makeInjection(avgAgg, NewAvgAggregation(), x, "other_path")
	avgAggSource, _ = avgAgg.Source()

	equal = compareAggregationToSource(maxAgg.Select("any_path"), avgAggSource)
	Expect(equal).To(Equal(true))
}

var _ = Describe("Injectable", func() {
	var maxAgg *MaxAggregation
	var maxAggSource interface{}

	var avgAgg *AvgAggregation
	var avgAggSource interface{}

	BeforeEach(func() {
		maxAgg = NewMaxAggregation()
		maxAggSource, _ = maxAgg.Source()

		avgAgg = NewAvgAggregation()
		avgAggSource, _ = avgAgg.Source()
	})
	Describe(" Injects", func() {

		It(" should inject aggregations", func() {
			testInjection(maxAgg, avgAgg, maxAggSource, avgAggSource, false)
		})

		It(" should injectX aggregations", func() {
			testInjection(maxAgg, avgAgg, maxAggSource, avgAggSource, true)

			bucketAggCopy, _ := maxAgg.Source()
			avgAgg := NewAvgAggregation()

			makeInjection(avgAgg, NewAvgAggregation(), false, "alter_path")
			makeInjection(maxAgg, avgAgg, true, "any_path")
			equal := compareAggregationToSource(maxAgg, bucketAggCopy)
			Expect(equal).To(Equal(true))

			makeInjection(maxAgg, avgAgg, true, "any_path", "other_path")
			equal = compareAggregationToSource(maxAgg, bucketAggCopy)
			Expect(equal).To(Equal(true))

		})
	})

	Describe(" Getters", func() {
		It(" should select aggregations", func() {
			makeInjection(maxAgg, avgAgg, false, "any_path")
			sel := maxAgg.Select("any_path")
			Expect(sel).NotTo(BeNil())
			Expect(sel.Source()).To(Equal(avgAggSource))
		})

		It(" should pop aggregations", func() {
			makeInjection(maxAgg, avgAgg, false, "any_path")

			maxAggSource, _ = maxAgg.Source()
			sel := maxAgg.Pop("any_path")
			Expect(sel).NotTo(BeNil())
			Expect(sel.Source()).To(Equal(avgAggSource))

			maxAggCopySrc, sourceErr := maxAgg.Source()
			Expect(sourceErr).ShouldNot(HaveOccurred())
			Expect(maxAggSource).NotTo(Equal(maxAggCopySrc))
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
