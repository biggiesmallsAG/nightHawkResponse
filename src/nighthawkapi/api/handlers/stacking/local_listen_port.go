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

type LocalListenPortStack struct {
	LocalPort    float64
	Process      string
	ComputerName []string `json:"ComputerName,omitempty"`
	LPortCount   int64
	ProcessCount int64
	//ComputerCount int64
}

func StackLocalListenPort(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		ret                            string
		sc                             StackingConfig
		conf                           *config.ConfigVars
		err                            error
		client                         *elastic.Client
		query                          elastic.Query
		agg_computer, agg_process, agg elastic.Aggregation
		sr                             *elastic.SearchResult
		ao                             *elastic.AggregationBucketKeyItems
		stackdata                      []LocalListenPortStack
		found                          bool
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/locallistport, Error encountered")
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

	agg_process = elastic.NewTermsAggregation().
		Field("Record.Process.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc).
		SubAggregation("computer", agg_computer)

	agg = elastic.NewTermsAggregation().
		Field("Record.LocalPort").
		Size(sc.SearchLimit).
		OrderByCount(sc.SortDesc).
		SubAggregation("process", agg_process)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator", "w32ports"), elastic.NewTermQuery("Record.State.keyword", "LISTENING"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName)).
			MustNot(elastic.NewTermQuery("Record.IsGoodTask.keyword", "true"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator", "w32ports"), elastic.NewTermQuery("Record.State.keyword", "LISTENING"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("lport", agg).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	ao, found = sr.Aggregations.Terms("lport")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, lportBucket := range ao.Buckets {
		lport, found := lportBucket.Aggregations.Terms("process")

		if found {
			for _, processBucket := range lport.Buckets {
				process, found := processBucket.Aggregations.Terms("computer")

				s := LocalListenPortStack{
					LocalPort:    lportBucket.Key.(float64),
					Process:      processBucket.Key.(string),
					LPortCount:   lportBucket.DocCount,
					ProcessCount: processBucket.DocCount,
				}

				if found {
					for _, computerBucket := range process.Buckets {
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
	api.LogDebug(api.DEBUG, "[+] GET /stacking/locallistport HTTP 200, returned LocalListenPort Stack")
	fmt.Fprintln(w, ret)
}
