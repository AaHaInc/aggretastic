package aggretastic_test

import (
	"encoding/json"
	"github.com/aahainc/aggretastic"
	"github.com/olivere/elastic/v7"
	. "github.com/onsi/ginkgo"
	"log"
)

var _ = Describe("AggsTree", func() {

	Context("ExtractLeafs", func() {

		It("should work", func() {
			agg := aggretastic.NewFilterAggregation().Filter(elastic.NewTermQuery("a", "b"))
			agg.Inject(aggretastic.NewChildrenAggregation().Type("child"), "down")
			agg.Inject(aggretastic.NewFilterAggregation().Filter(elastic.NewTermQuery("b", "c")), "filtered")
			agg.Inject(aggretastic.NewChildrenAggregation().Type("foobar"), "filtered", "deeper")
			agg.Inject(aggretastic.NewChildrenAggregation().Type("foobarx"), "filtered", "deeper-x")

			log.Println("leafs ", agg.ExtractLeafPaths())

			p, e := agg.InjectSafe(
				aggretastic.NewChildrenAggregation().Type("foobar").SubAggregation(
					"deeper",
					aggretastic.NewChildrenAggregation().Type("foobar").SubAggregation(
						"semi-final",
						aggretastic.NewFilterAggregation().Filter(elastic.NewTermQuery("cc", "ee")).SubAggregation(
							"final", aggretastic.NewSumAggregation().Field("abc"),
						),
					),
				),
				"filtered",
			)
			log.Println("err ", e)
			log.Println("result paths ", p)

			s, _ := agg.Source()
			j, _ := json.Marshal(s)
			log.Println(string(j))
		})
	})
})
