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

type RunKeyStack struct {
	RegKey      string
	RegVal      string
	Md5sum      string
	KeyCount    int64
	Md5sumCount int64
}

func StackRunKey(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		ret                           string
		sc                            StackingConfig
		conf                          *config.ConfigVars
		err                           error
		client                        *elastic.Client
		query                         elastic.Query
		agg_filemd5, agg_val, agg_key elastic.Aggregation
		sr                            *elastic.SearchResult
		ao                            *elastic.AggregationBucketKeyItems
		stackdata                     []RunKeyStack
		found                         bool
	)
	// Read http parameters
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/runkey, Error encountered")
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

	agg_filemd5 = elastic.NewTermsAggregation().
		Field("Record.Md5sum.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc)

	agg_val = elastic.NewTermsAggregation().
		Field("Record.Registry.Text.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc).
		SubAggregation("md5sum", agg_filemd5)

	agg_key = elastic.NewTermsAggregation().
		Field("Record.StackPath.keyword").
		Size(sc.SearchLimit).
		OrderByCount(sc.SortDesc).
		SubAggregation("reg_val", agg_val)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32scripting-persistence"),
				elastic.NewTermQuery("Record.PersistenceType.keyword", "Registry"),
				elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName),
				elastic.NewWildcardQuery("Record.RegPath.keyword", "*Run*")).
			MustNot(elastic.NewTermQuery("Record.IsGoodService.keyword", "true"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32scripting-persistence"),
				elastic.NewTermQuery("Record.PersistenceType.keyword", "Registry"),
				elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName),
				elastic.NewWildcardQuery("Record.RegPath.keyword", "*Run*"))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("reg_key", agg_key).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	ao, found = sr.Aggregations.Terms("reg_key")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, keyBucket := range ao.Buckets {
		val, found := keyBucket.Aggregations.Terms("reg_val")

		if found {
			for _, valBucket := range val.Buckets {
				md5sum, found := valBucket.Aggregations.Terms("md5sum")

				if found {
					for _, md5sumBucket := range md5sum.Buckets {

						s := RunKeyStack{
							RegKey:      keyBucket.Key.(string),
							RegVal:      valBucket.Key.(string),
							Md5sum:      md5sumBucket.Key.(string),
							KeyCount:    keyBucket.DocCount,
							Md5sumCount: md5sumBucket.DocCount,
						}
						stackdata = append(stackdata, s)
					}
				}
			}
		}

	}

	ret = api.HttpSuccessMessage("200", &stackdata, sr.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /stacking/runkey HTTP 200, returned RunKey Stack"))
	fmt.Fprintln(w, ret)
}
