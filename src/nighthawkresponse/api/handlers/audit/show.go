package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/config"
	"strconv"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/gorilla/mux"
)

const tsLayout = "2006-01-02T15:04:05Z"

type ReturnBucket struct {
	Key      string `json:"key"`
	DocCount int64  `json:"doc_count"`
}

type CaseBucket struct {
	KeyAsString string `json:"key_as_string"`
	DocCount    int64  `json:"doc_count"`
}

type AuditFilter struct {
	ColId  string `json:"colId"`
	Filter struct {
		Value      interface{} `json:"filter"`
		FilterType string      `json:"filterType"`
		MatchType  string      `json:"type"`
		StartDate  string      `json:"dateFrom"`
		EndDate    string      `json:"dateTo"`
	} `json:"filterOn"`
}

func GetDocById(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		client      *elastic.Client
		// query elastic.Query
		cases *elastic.GetResult
		conf  *config.ConfigVars
		err   error
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

	cases, err = client.Get().
		Index(conf.Elastic.Elastic_index).
		Type("audit_type").
		Id(vars["doc_id"]).
		Routing(vars["endpoint"]).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
		return
	}

	ret = api.HttpSuccessMessage("200", cases.Source, 1)
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /show/doc/%s/%s HTTP 200, returned document.", vars["doc_id"], vars["endpoint"]))
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

	query = elastic.NewBoolQuery().Must(elastic.NewTermQuery("CaseInfo.CaseName.keyword", vars_case))

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
		Field("CaseInfo.CaseDate").
		Size(500)

	query = elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("CaseInfo.CaseName.keyword", vars_case),
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
		elastic.NewTermQuery("CaseInfo.CaseName.keyword", vars_case),
		elastic.NewTermQuery("ComputerName.keyword", vars_endpoint),
		elastic.NewTermQuery("CaseInfo.CaseDate", vars_case_date))

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

	var (
		from      int
		sortField string
		sortOrder bool
		method    string
		client    *elastic.Client
		at        *elastic.SearchResult
		query     elastic.Query
		filter    AuditFilter
	)

	// Setting default value for sort_order
	// and sort_field
	sortField = "Record.TlnTime" // default field to sort
	sortOrder = true             // sort in ascending order

	uriSortField := r.URL.Query().Get("sort")
	if uriSortField != "" {
		sortField = uriSortField
	}
	uriSortOrder := r.URL.Query().Get("order")

	// Values from RequestURI
	varsCase := vars["case"]
	varsEndpoint := vars["endpoint"]
	varsAuditType := vars["audittype"]
	varsCaseDate := vars["case_date"]

	switch uriSortOrder {
	case "desc":
		sortOrder = true
		break
	case "asc":
		sortOrder = false
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

	var queries []elastic.Query

	// 0xredskull - updated to fix contains and equals filterType
	queries = append(queries, elastic.NewTermQuery("CaseInfo.CaseName.keyword", varsCase))
	queries = append(queries, elastic.NewTermQuery("ComputerName.keyword", varsEndpoint))
	queries = append(queries, elastic.NewTermQuery("AuditType.Generator.keyword", varsAuditType))
	queries = append(queries, elastic.NewTermQuery("CaseInfo.CaseDate", varsCaseDate))

	// Query for GET method. This will be changed by POST Method
	query = elastic.NewBoolQuery().Must(queries...)

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, "[+] POST /show/audit/filter, Error encountered")
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}
		json.Unmarshal(body, &filter)

		switch filter.Filter.Value.(type) {
		case float64:
			filterValue := int(filter.Filter.Value.(float64))
			queries = append(queries, elastic.NewTermQuery(filter.ColId, filterValue))
			query = elastic.NewBoolQuery().Must(queries...)

		default:
			switch filter.ColId {
			case "Record.Index":
				// Handle interger based query
				filterValue, err := strconv.Atoi(filter.Filter.Value.(string))
				if err != nil {
					api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
					return
				}

				queries = append(queries, elastic.NewTermQuery(filter.ColId, filterValue))
				query = elastic.NewBoolQuery().Must(queries...)

			case "Record.GenTime", "Record.WriteTime":
				var filterQuery elastic.Query

				fvStartDate := "1970-01-01T00:00:00Z"
				fvEndDate := "2030-01-01T00:00:00Z"

				switch filter.Filter.MatchType {
				case "equal", "equals":
					fvStartDate = fmt.Sprintf("%vT00:00:00Z", filter.Filter.StartDate)
					fvEndDate = fmt.Sprintf("%vT23:59:59Z", filter.Filter.StartDate)
					filterQuery = elastic.NewRangeQuery(filter.ColId).Gte(fvStartDate).Lte(fvEndDate)
				case "notEqual", "notEquals":
					fvStartDate = fmt.Sprintf("%vT00:00:00Z", filter.Filter.StartDate)
					fvEndDate = fmt.Sprintf("%vT23:59:59Z", filter.Filter.StartDate)
					filterQuery = elastic.NewRangeQuery(filter.ColId).Lt(fvStartDate).Gt(fvEndDate)
				case "greaterThan":
					fvStartDate = fmt.Sprintf("%vT00:00:00Z", filter.Filter.StartDate)
					filterQuery = elastic.NewRangeQuery(filter.ColId).Gt(fvStartDate)
				case "lessThan":
					fvStartDate = fmt.Sprintf("%vT00:00:00Z", filter.Filter.StartDate)
					filterQuery = elastic.NewRangeQuery(filter.ColId).Lt(fvStartDate)
				default:
					if filter.Filter.StartDate != "" {
						fvStartDate = fmt.Sprintf("%vT00:00:00Z", filter.Filter.StartDate)
					}
					if filter.Filter.EndDate != "" {
						fvEndDate = fmt.Sprintf("%vT00:00:00Z", filter.Filter.EndDate)
					}
					filterQuery = elastic.NewRangeQuery(filter.ColId).Gte(fvStartDate).Lte(fvEndDate)
				}
				query = elastic.NewBoolQuery().Must(queries...).Filter(filterQuery)

			default:
				// String Data
				filterValue := fmt.Sprintf("%v", filter.Filter.Value)
				switch filter.Filter.MatchType {
				case "equal", "equals":
					queries = append(queries, elastic.NewTermQuery(filter.ColId, filterValue))
					query = elastic.NewBoolQuery().Must(queries...)
				case "notEqual", "notEquals":
					query = elastic.NewBoolQuery().Must(queries...).MustNot(elastic.NewTermQuery(filter.ColId, filterValue))
				case "notContain", "notContains":
					filterValue = fmt.Sprintf("*%v*", filterValue)
					query = elastic.NewBoolQuery().Must(queries...).MustNot(elastic.NewWildcardQuery(filter.ColId, filterValue))
				default:
					// contains, notContains, greaterThan,....
					filterValue = fmt.Sprintf("*%v*", filterValue)
					queries = append(queries, elastic.NewWildcardQuery(filter.ColId, filterValue))
					query = elastic.NewBoolQuery().Must(queries...)
				}
			}
		}
	}

	// 0xredskull - print QueryDSL
	/*
		qsrc, _ := query.Source()
		qdsl, _ := json.MarshalIndent(qsrc, "", " ")
		fmt.Println(string(qdsl))
	*/

	at, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(query).
		Size(100).
		Sort(sortField, sortOrder).
		From(from).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, errors.New(fmt.Sprintf("%s - %s", r.RequestURI, err.Error())))
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	if at.TotalHits() < 1 {
		api.HttpResponseReturn(w, r, "failed", "No matching data for filter", nil)
		return
	}

	api.HttpResponseReturn(w, r, "success", "Filtered data", at.Hits.Hits)
	return

}
