package stacking

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/audit"
	"nighthawkapi/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

func (sc *StackingConfig) GetRecordField() (record_field string) {
	switch sc.Type {
	case "service":
		// must switch on stack type to get
		record_field = "Record.Path.keyword"
		break
	case "prefetch":
		record_field = "Record.ApplicationFullPath.keyword"
		break
	case "persistence":
		record_field = "Record.Path.keyword"
		break
	case "runkey":
		record_field = "Record.StackPath.keyword"
		break
	case "task":
		record_field = "Record.Name.keyword"
		break
	case "dns/a":
		record_field = "Record.RecordDataList.Ipv4Address.keywordgrace "
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] POST /stacking/context, Error encountered")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	sc.LoadParams(body)

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

	record_field := sc.GetRecordField()

	query = elastic.NewBoolQuery().
		Must(elastic.NewTermQuery(record_field, sc.ContextItem))

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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] POST /stacking/context/endpoint, Error encountered")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	sc.LoadParams(body)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	record_field := sc.GetRecordField()

	query = elastic.NewBoolQuery().
		Must(elastic.NewTermQuery(record_field, sc.ContextItem)).
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
