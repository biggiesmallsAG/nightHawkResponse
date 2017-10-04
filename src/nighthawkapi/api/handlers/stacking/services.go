package stacking

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

type ServiceStack struct {
	ServiceName  string
	ServicePath  string
	ComputerName []string `json:"ComputerName,omitempty"`
	NameCount    int64
	PathCount    int64
}

func StackServices(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	var (
		ret                         string
		sc                          StackingConfig
		conf                        *config.ConfigVars
		err                         error
		client                      *elastic.Client
		query                       elastic.Query
		agg_computer, agg_path, agg elastic.Aggregation
		sr                          *elastic.SearchResult
		svc                         *elastic.AggregationBucketKeyItems
		stackdata                   []ServiceStack
		found                       bool
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/services, Error encountered")
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

	agg_computer = elastic.NewTermsAggregation().
		Field("ComputerName.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc)

	agg_path = elastic.NewTermsAggregation().
		Field("Record.Path.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc).
		SubAggregation("computer", agg_computer)

	agg = elastic.NewTermsAggregation().
		Field("Record.Name.keyword").
		Size(sc.SearchLimit).
		OrderByCount(sc.SortDesc).
		SubAggregation("svcpath", agg_path)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32services"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName)).
			MustNot(elastic.NewTermQuery("Record.IsGoodService.keyword", "true"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32services"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("svcname", agg).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
	}

	// Parsing SearchResult
	svc, found = sr.Aggregations.Terms("svcname")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, svcnameBucket := range svc.Buckets {
		svcname, found := svcnameBucket.Aggregations.Terms("svcpath")

		if found {
			for _, svcpathBucket := range svcname.Buckets {
				svcpath, found := svcpathBucket.Aggregations.Terms("computer")
				s := ServiceStack{
					ServiceName: svcnameBucket.Key.(string),
					ServicePath: svcpathBucket.Key.(string),
					NameCount:   svcnameBucket.DocCount,
					PathCount:   svcpathBucket.DocCount,
				}

				if found {
					for _, computerBucket := range svcpath.Buckets {
						s.ComputerName = append(s.ComputerName, computerBucket.Key.(string))
					}
				}
				if !sc.ShowComputer {
					s.ComputerName = nil
				}
				stackdata = append(stackdata, s)
			}
		}
	}

	ret = api.HttpSuccessMessage("200", &stackdata, sr.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /stacking/services HTTP 200, returned Service Stack")
	fmt.Fprintln(w, ret)
}
