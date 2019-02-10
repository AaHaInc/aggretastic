# Aggretastic

The `aggretastic` package enables manipulation of a `olivere/elastic` subAggregations tree with high fidelity. 


## About `olivere/elastic` 

 `olivere/elastic` (further in text it is named elastic) is an Elasticsearch client for the Go programming language. 
This package is some kind of ORM for Elasticsearch queries - using  Golang structs and its methods, we can simply compile raw JSON query that can be sent to Elasticsearch.

**Example:**

    #import olivere/elastic
    
    #build new query
    query := elastic.NewBoolQuery().
      Must(elastic.NewTermQuery("action_type.keyword", "some-action")).
      Filter(elastic.NewBoolQuery().
            MustNot(elastic.NewTermQuery("disabled", true))
      )
    
    #such JSON will be returned (that can be sent to elastic)
    source, err := query.Source()
    
    # source contents (JSON) 
    > {"bool":{"filter":{"bool":{"must_not":{"term":{"disabled":true}}}},"must":{"term":{"action_type.keyword":"some-action"}}}}
    



## Where does `elastic`  break?

ElasticSearch aggregations support adding subAggregations inside aggregations. In `elastic` that is done via inserting a new aggregation into a private map, that can’t be accessed within the package. There is no methods of walking, searching, or modifying the `subAggregation` tree.
In real cases it’s very useful to modify the subAggregation tree - on the fly…


## Here comes Aggretastic

To achieve breaking the limits explaining above we created Aggretastic. It’s constructs similar to `elastic` aggregations, but with support of helper methods: **Inject**, **InjectX**, **Select** and **Pop**.


## Installation

    $ go get -u github.com/AaHaInc/aggretastic


## Update project files to `olivere/elastic` upstream

     $ go generate



## Quick Start 


Add this import line to the file you're working in:

    import "github.com/AaHaInc/aggretastic"

Let’s create sample aggregation

    var result aggretastic.Aggregation = aggretastic.NewDateHistogramAggregation()
    >>{"date_histogram":{"interval":""}}

And now let’s inject some new aggregations into it

    nested := aggretastic.NewNestedAggregation()
    
    result.Inject(nested, "exp")
    >> nil
    >>result == {"aggregations":{"exp":{"nested":{"path":""}}},"date_histogram":{"interval":""}}
    
    result.Inject(nested, "exp", "lit")
    >> nil
    >>result == {"aggregations":{"exp":{"aggregations":{"lit":{"nested":{"path":""}}},"nested":{"path":""}}},"date_histogram":{"interval":""}}
    
    result.InjectX(nested.Path("anypath"), "exp", "lit")
    >>nil
    >>nothing changes here

You have to define all path entries before injection:

    result.Inject(aggretastic.NewNestedAggregation(), "one", "two", "three")
    >> path is not selectable
    >>{"aggregations":{"exp":{"aggregations":{"lit":{"nested":{"path":""}}},"nested":{"path":""}}},"date_histogram":{"interval":""}} 

Select will  return you aggregation in path or nil

    result.Select("exp")
    >> {"aggregations":{"lit":{"nested":{"path":""}}},"nested":{"path":""}}
    >> returns aggregation from exp path

Pop will do the same as select. Also it will delete aggregation from path

    result.Pop("exp", "lit")
    >> {"nested":{"path":""}}
    >> result now equals to {"aggregations":{"exp":{"nested":{"path":""}}},"date_histogram":{"interval":""}}


## Types

All Aggregations implements single interface. It contains several functions for subAggregations on-the-fly modifications. Aggregation types is next:
**Injectable:**
This is a single injectable Aggregation. All operations affects `Aggregation.subAggregations` field
**Aggregations:**
This is a  set of Aggregations. 
**Not Injectable:**
This is a single aggregation which are not contain any subAggregations. All method returns error or nil

