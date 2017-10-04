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

type PrefetchStack struct {
	AppFullPath      string
	Path             string
	ComputerName     []string `json:"ComputerName,omitempty"`
	AppFullPathCount int64
	PathCount        int64
	//ComputerCount    int64
}

func StackPrefetch(w http.ResponseWriter, r *http.Request) {
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
		pf                          *elastic.AggregationBucketKeyItems
		stackdata                   []PrefetchStack
		found                       bool
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/prefetch, Error encountered")
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
		Field("Record.ApplicationFullPath.keyword").
		Size(sc.SearchLimit).
		OrderByCount(sc.SortDesc).
		SubAggregation("path", agg_path)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32prefetch"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName)).
			MustNot(elastic.NewTermQuery("Record.IsGoodPrefetch.keyword", "true"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32prefetch"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("afp", agg).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	pf, found = sr.Aggregations.Terms("afp")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, afpBucket := range pf.Buckets {
		afp, found := afpBucket.Aggregations.Terms("path")

		if found {
			for _, pathBucket := range afp.Buckets {
				path, found := pathBucket.Aggregations.Terms("computer")
				s := PrefetchStack{
					AppFullPath:      afpBucket.Key.(string),
					Path:             pathBucket.Key.(string),
					AppFullPathCount: afpBucket.DocCount,
					PathCount:        pathBucket.DocCount,
				}

				if found {
					for _, computerBucket := range path.Buckets {
						s.ComputerName = append(s.ComputerName, computerBucket.Key.(string))
					}
					//s.ComputerCount = int64(len(s.ComputerName))
				}
				if !sc.ShowComputer {
					s.ComputerName = nil
				}
				stackdata = append(stackdata, s)
			}
		}
	}

	ret = api.HttpSuccessMessage("200", &stackdata, sr.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /stacking/prefetch HTTP 200, returned Prefetch Stack")
	fmt.Fprintln(w, ret)
}
