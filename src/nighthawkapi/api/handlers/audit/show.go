package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"
	"strconv"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/gorilla/mux"
)

type ReturnBucket struct {
	Key      string `json:"key"`
	DocCount int64  `json:"doc_count"`
}

type CaseBucket struct {
	KeyAsString string `json:"key_as_string"`
	DocCount    int64  `json:"doc_count"`
}

type FilterOn struct {
	ColId  string `json:"colId"`
	Filter struct {
		FilterQuery interface{} `json:"filter"`
		FilterType  string      `json:"filterType"`
	} `json:"filterOn"`
}

func GetDocById(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		client      *elastic.Client
		query       elastic.Query
		cases       *elastic.SearchResult
		conf        *config.ConfigVars
		err         error
	)

	vars := mux.Vars(r)

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

	query = elastic.NewTermQuery("_id", vars["doc_id"])

	cases, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(query).
		Size(1).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
		return
	}

	ret = api.HttpSuccessMessage("200", &cases.Hits.Hits[0], cases.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /show/doc/%s HTTP 200, returned document.", vars["doc_id"]))
	fmt.Fprintln(w, ret)
}

func GetEndpointByCase(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		conf        *config.ConfigVars
		err         error
		client      *elastic.Client
		query       elastic.Query
		agg         elastic.Aggregation
		at          *elastic.SearchResult
		a           *elastic.AggregationBucketKeyItems
		i           []ReturnBucket
		found       bool
	)

	vars := mux.Vars(r)
	vars_case := vars["case"]

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

	query = elastic.NewBoolQuery().Must(elastic.NewTermQuery("CaseInfo.case_name", vars_case))

	at, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(query).
		Size(0).
		Aggregation("endpoints", agg).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
		return
	}

	a, found = at.Aggregations.Terms("endpoints")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	for _, v := range a.Buckets {
		i = append(i, ReturnBucket{
			Key:      v.Key.(string),
			DocCount: v.DocCount,
		})
	}

	ret = api.HttpSuccessMessage("200", &i, at.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /show/%s HTTP 200, returned endpoints list for case.", vars_case))
	fmt.Fprintln(w, ret)
}

func GetCasedateByEndpoint(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		conf        *config.ConfigVars
		err         error
		client      *elastic.Client
		query       elastic.Query
		agg         elastic.Aggregation
		at          *elastic.SearchResult
		a           *elastic.AggregationBucketKeyItems
		i           []CaseBucket
		found       bool
	)

	vars := mux.Vars(r)
	vars_case := vars["case"]
	vars_endpoint := vars["endpoint"]

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
		Field("CaseInfo.case_date").
		Size(500)

	query = elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("CaseInfo.case_name.keyword", vars_case),
		elastic.NewTermQuery("ComputerName.keyword", vars_endpoint))

	at, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(query).
		Size(0).
		Aggregation("case_dates", agg).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
	}

	a, found = at.Aggregations.Terms("case_dates")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	for _, v := range a.Buckets {
		i = append(i, CaseBucket{
			KeyAsString: *v.KeyAsString,
			DocCount:    v.DocCount,
		})
	}

	ret = api.HttpSuccessMessage("200", &i, at.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /show/%s/%s HTTP 200, returned case dates for endpoint.", vars_case, vars_endpoint))
	fmt.Fprintln(w, ret)
}

func GetAuditTypeByEndpointAndCase(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		client      *elastic.Client
		query       elastic.Query
		agg         elastic.Aggregation
		at          *elastic.SearchResult
		a           *elastic.AggregationBucketKeyItems
		i           []ReturnBucket
		found       bool
	)

	vars := mux.Vars(r)
	vars_case := vars["case"]
	vars_endpoint := vars["endpoint"]
	vars_case_date := vars["case_date"]

	conf, err := config.ReadConfFile()
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
		Field("AuditType.Generator.keyword").
		Size(16)

	query = elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("CaseInfo.case_name.keyword", vars_case),
		elastic.NewTermQuery("ComputerName.keyword", vars_endpoint),
		elastic.NewTermQuery("CaseInfo.case_date", vars_case_date))

	at, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(query).
		Size(0).
		Aggregation("audits", agg).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
	}

	a, found = at.Aggregations.Terms("audits")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	for _, v := range a.Buckets {
		i = append(i, ReturnBucket{
			Key:      v.Key.(string),
			DocCount: v.DocCount,
		})
	}

	ret = api.HttpSuccessMessage("200", &i, at.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /show/%s/%s/%s HTTP 200, returned audit list for endpoint.", vars_case, vars_endpoint, vars_case_date))
	fmt.Fprintln(w, ret)
}

func GetAuditDataByAuditGenerator(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	order := r.URL.Query().Get("order")

	var (
		from        int
		sort_order  bool
		method, ret string
		client      *elastic.Client
		at          *elastic.SearchResult
		query       elastic.Query
		filter      FilterOn
	)

	vars_case := vars["case"]
	vars_endpoint := vars["endpoint"]
	vars_audittype := vars["audittype"]
	vars_case_date := vars["case_date"]

	switch order {
	case "desc":
		sort_order = true
		break
	case "asc":
		sort_order = false
	}

	if f, err := strconv.Atoi(r.URL.Query().Get("from")); err == nil {
		from = f
	}

	conf, err := config.ReadConfFile()
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

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, "[+] POST /show/audit/filter, Error encountered")
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}
		json.Unmarshal(body, &filter)
		query = elastic.NewBoolQuery().Must(
			elastic.NewTermQuery("CaseInfo.case_name.keyword", vars_case),
			elastic.NewTermQuery("ComputerName.keyword", vars_endpoint),
			elastic.NewTermQuery("AuditType.Generator.keyword", vars_audittype),
			elastic.NewTermQuery("CaseInfo.case_date", vars_case_date),
			elastic.NewWildcardQuery(filter.ColId, filter.Filter.FilterQuery.(string)))
	} else {
		query = elastic.NewBoolQuery().Must(
			elastic.NewTermQuery("CaseInfo.case_name.keyword", vars_case),
			elastic.NewTermQuery("ComputerName.keyword", vars_endpoint),
			elastic.NewTermQuery("AuditType.Generator.keyword", vars_audittype),
			elastic.NewTermQuery("CaseInfo.case_date", vars_case_date))
	}

	at, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(query).
		Size(100).
		Sort(r.URL.Query().Get("sort"), sort_order).
		From(from).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
	}

	ret = api.HttpSuccessMessage("200", &at.Hits.Hits, at.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /show/%s/%s/%s/%s HTTP 200, returned audit data for endpoint/case", vars_case, vars_endpoint, vars_case_date, vars_audittype))
	fmt.Fprintln(w, ret)
}
