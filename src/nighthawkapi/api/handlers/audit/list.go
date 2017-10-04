package audit

import (
	"context"
	"fmt"
	"net/http"
	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

type Case struct {
	CaseName string `json:"case_name"`
}

type Endpoint struct {
	Endpoints string `json:"endpoint_name,omitempty"`
}

type AuditType struct {
	AuditType string `json:"audit_name"`
}

func GetCaseList(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		conf        *config.ConfigVars
		err         error
		client      *elastic.Client
		agg         elastic.Aggregation
		cases       *elastic.SearchResult
		a           *elastic.AggregationBucketKeyItems
		_c          Case
		found       bool
	)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if conf.Elastic.Elastic_ssl {
		method = api.HTTPS
	} else {
		method = api.HTTP
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s%s:%d", method, conf.Elastic.Elastic_server, conf.Elastic.Elastic_port)))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}
	agg = elastic.NewTermsAggregation().
		Field("CaseInfo.case_name.keyword").
		Size(5000)

	cases, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(elastic.NewMatchAllQuery()).
		Size(0).
		Aggregation("cases", agg).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
		return
	}

	a, found = cases.Aggregations.Terms("cases")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	_r := make([]interface{}, len(a.Buckets))

	for i, v := range a.Buckets {
		_c.CaseName = v.Key.(string)
		_r[i] = _c
	}

	ret = api.HttpSuccessMessage("200", &_r, cases.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /list/cases HTTP 200, returned case list.")
	fmt.Fprintln(w, ret)
}

func GetEndpointList(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		conf        *config.ConfigVars
		err         error
		client      *elastic.Client
		agg         elastic.Aggregation
		ep          *elastic.SearchResult
		a           *elastic.AggregationBucketKeyItems
		_e          Endpoint
		found       bool
	)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if conf.Elastic.Elastic_ssl {
		method = api.HTTPS
	} else {
		method = api.HTTP
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s%s:%d", method, conf.Elastic.Elastic_server, conf.Elastic.Elastic_port)))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}
	agg = elastic.NewTermsAggregation().
		Field("ComputerName.keyword").
		Size(50000)

	ep, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(elastic.NewMatchAllQuery()).
		Size(0).
		Aggregation("endpoints", agg).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
	}

	a, found = ep.Aggregations.Terms("endpoints")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	_r := make([]interface{}, len(a.Buckets))

	for i, v := range a.Buckets {
		_e.Endpoints = v.Key.(string)
		_r[i] = _e
	}

	ret = api.HttpSuccessMessage("200", &_r, ep.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /list/endpoints HTTP 200, returned endpoint list.")
	fmt.Fprintln(w, ret)
}

func GetAuditTypeList(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	var (
		method, ret string
		conf        *config.ConfigVars
		err         error
		client      *elastic.Client
		agg         elastic.Aggregation
		at          *elastic.SearchResult
		a           *elastic.AggregationBucketKeyItems
		_a          AuditType
		found       bool
	)
	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if conf.Elastic.Elastic_ssl {
		method = api.HTTPS
	} else {
		method = api.HTTP
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s%s:%d", method, conf.Elastic.Elastic_server, conf.Elastic.Elastic_port)))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}
	agg = elastic.NewTermsAggregation().
		Field("AuditType.Generator").
		Size(16)

	at, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(elastic.NewMatchAllQuery()).
		Size(16).
		Aggregation("endpoints", agg).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	a, found = at.Aggregations.Terms("endpoints")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	_r := make([]interface{}, len(a.Buckets))

	for i, v := range a.Buckets {
		_a.AuditType = v.Key.(string)
		_r[i] = _a
	}

	ret = api.HttpSuccessMessage("200", &_a, at.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /list/audittypes HTTP 200, returned audittype list.")
	fmt.Fprintln(w, ret)
}
