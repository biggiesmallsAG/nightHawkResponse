package stacking

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

type TaskStack struct {
	TaskName     string
	TaskPath     string
	ComputerName []string `json:"ComputerName,omitempty"`
	NameCount    int64
	PathCount    int64
	//ComputerCount int64
}

func StackTasks(w http.ResponseWriter, r *http.Request) {
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
		stackdata                   []TaskStack
		found                       bool
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/task, Error encountered")
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
		SubAggregation("tskpath", agg_path)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32tasks"), elastic.NewWildcardQuery("CaseInfo.CaseName.keyword", sc.CaseName)).
			MustNot(elastic.NewTermQuery("Record.IsGoodTask.keyword", "true"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32tasks"), elastic.NewWildcardQuery("CaseInfo.CaseName.keyword", sc.CaseName))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("tskname", agg).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	tsk, found := sr.Aggregations.Terms("tskname")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, tsknameBucket := range tsk.Buckets {
		tskname, found := tsknameBucket.Aggregations.Terms("tskpath")

		if found {
			for _, tskpathBucket := range tskname.Buckets {
				tskpath, found := tskpathBucket.Aggregations.Terms("computer")

				s := TaskStack{
					TaskName:  tsknameBucket.Key.(string),
					TaskPath:  tskpathBucket.Key.(string),
					NameCount: tsknameBucket.DocCount,
					PathCount: tskpathBucket.DocCount,
				}

				if found {
					for _, computerBucket := range tskpath.Buckets {
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
	api.LogDebug(api.DEBUG, "[+] GET /stacking/task HTTP 200, returned Task Stack")
	fmt.Fprintln(w, ret)
}
