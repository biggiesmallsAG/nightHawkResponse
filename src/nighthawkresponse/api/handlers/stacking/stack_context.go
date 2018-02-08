package stacking

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/audit"
	"nighthawkresponse/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

func (sc *StackingConfig) GetRecordField() (recordField string) {
	switch sc.Type {
	case "service":
		// must switch on stack type to get
		recordField = "Record.Path.keyword"
		break
	case "prefetch":
		recordField = "Record.ApplicationFullPath.keyword"
		break
	case "persistence":
		recordField = "Record.Path.keyword"
		break
	case "runkey":
		recordField = "Record.StackPath.keyword"
		break
	case "task":
		recordField = "Record.Name.keyword"
		break
	case "dns/a":
		recordField = "Record.RecordDataList.Ipv4Address.keyword"
		break
	default:
		// Default field needs to be return
		// If not set exception occurs at aggregation
		recordField = "Record.Name.keyword"
	}
	return
}

func GetStackContext(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		sc          StackingConfig
		query       *elastic.BoolQuery
		ret         string
		conf        *config.ConfigVars
		err         error
		client      *elastic.Client
		aggregation elastic.Aggregation
		sr          *elastic.SearchResult
		endpoints   *elastic.AggregationBucketKeyItems
		_a          []audit.Endpoint
		found       bool
	)

	// Load default values to StackingConfig object
	sc.Default()

	if r.ContentLength > 0 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, "[+] POST /stacking/context, Error encountered")
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}

		sc.LoadParams(body)
	}

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	aggregation = elastic.NewTermsAggregation().
		Field("ComputerName.keyword").
		Size(100)

	recordField := sc.GetRecordField()

	query = elastic.NewBoolQuery().
		Must(elastic.NewTermQuery(recordField, sc.ContextItem))

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(100).
		Aggregation("endpoint", aggregation).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	endpoints, found = sr.Aggregations.Terms("endpoint")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, v := range endpoints.Buckets {
		_a = append(_a, audit.Endpoint{
			Endpoints: v.Key.(string),
		})
	}

	ret = api.HttpSuccessMessage("200", &_a, sr.TotalHits())
	api.LogDebug(api.DEBUG, "[+] POST /stacking/context HTTP 200, returned context list.")
	fmt.Fprintln(w, ret)
}

func GetStackContextByEndpoint(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		sc     StackingConfig
		query  *elastic.BoolQuery
		ret    string
		conf   *config.ConfigVars
		err    error
		client *elastic.Client
		sr     *elastic.SearchResult
	)

	sc.Default()

	if r.ContentLength > 0 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, "[+] POST /stacking/context/endpoint, Error encountered")
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}

		sc.LoadParams(body)
	}

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	recordField := sc.GetRecordField()

	query = elastic.NewBoolQuery().
		Must(elastic.NewTermQuery(recordField, sc.ContextItem)).
		Must(elastic.NewTermQuery("ComputerName.keyword", sc.EndpointName))

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(1).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	ret = api.HttpSuccessMessage("200", &sr.Hits.Hits, sr.TotalHits())
	api.LogDebug(api.DEBUG, "[+] POST /stacking/context/endpoint HTTP 200, returned audit data for endpoint")
	fmt.Fprintln(w, ret)
}
