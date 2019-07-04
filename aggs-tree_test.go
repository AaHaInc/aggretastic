package aggretastic_test

import (
	"encoding/json"
	"github.com/aahainc/aggretastic"
	"github.com/olivere/elastic/v7"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AggsTree", func() {

	Context("ExtractLeafs", func() {

		// Tests are really poor for now @todo: tests
		It("should work", func() {
			agg := aggretastic.NewFilterAggregation().Filter(elastic.NewTermQuery("a", "b"))
			agg.Inject(aggretastic.NewChildrenAggregation().Type("child"), "down")
			agg.Inject(aggretastic.NewFilterAggregation().Filter(elastic.NewTermQuery("b", "c")), "filtered")
			agg.Inject(aggretastic.NewChildrenAggregation().Type("foobar"), "filtered", "deeper")
			agg.Inject(aggretastic.NewChildrenAggregation().Type("foobarx"), "filtered", "deeper-x")

			leafPaths := agg.ExtractLeafPaths()
			Expect(leafPaths).To(HaveLen(3))
			Expect(leafPaths[0]).To(Equal([]string{"down"}))
			Expect(leafPaths[1]).To(Equal([]string{"filtered", "deeper"}))
			Expect(leafPaths[2]).To(Equal([]string{"filtered", "deeper-x"}))

			resultPaths, err := agg.InjectSafe(
				aggretastic.NewChildrenAggregation().Type("foobar").SubAggregation(
					"deeper",
					aggretastic.NewChildrenAggregation().Type("foobar").SubAggregation(
						"pre-final",
						aggretastic.NewFilterAggregation().Filter(elastic.NewTermQuery("cc", "ee")).SubAggregation(
							"final", aggretastic.NewSumAggregation().Field("abc"),
						),
					),
				),
				"filtered",
			)

			Expect(err).ShouldNot(HaveOccurred())

			Expect(resultPaths).To(HaveLen(1))
			Expect(resultPaths[0]).To(Equal([]string{"filtered", "deeper", "pre-final", "final"}))

			// how agg was changed
			s, _ := agg.Source()
			j, _ := json.Marshal(s)
			Expect(string(j)).To(Equal(`{"aggregations":{"down":{"children":{"type":"child"}},"filtered":{"aggregations":{"deeper":{"aggregations":{"pre-final":{"aggregations":{"final":{"sum":{"field":"abc"}}},"filter":{"term":{"cc":"ee"}}}},"children":{"type":"foobar"}},"deeper-x":{"children":{"type":"foobarx"}}},"filter":{"term":{"b":"c"}}}},"filter":{"term":{"a":"b"}}}`))
		})
	})
})
