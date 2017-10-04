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

type PersistenceStack struct {
	Type         string
	Path         string
	FileMd5      string
	ComputerName []string `json:"ComputerName,omitempty"`
	TypeCount    int64
	PathCount    int64
	FileMd5Count int64
	//ComputerCount int64
}

func StackPersistence(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		ret                                            string
		sc                                             StackingConfig
		conf                                           *config.ConfigVars
		err                                            error
		client                                         *elastic.Client
		query                                          elastic.Query
		agg_computer, agg_filemd5, agg_path, agg_ptype elastic.Aggregation
		sr                                             *elastic.SearchResult
		p                                              *elastic.AggregationBucketKeyItems
		stackdata                                      []PersistenceStack
		found                                          bool
	)

	// Read http parameters
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/persistence, Error encountered")
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

	agg_filemd5 = elastic.NewTermsAggregation().
		Field("Record.File.Md5sum.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc).
		SubAggregation("computer", agg_computer)

	agg_path = elastic.NewTermsAggregation().
		Field("Record.Path.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc).
		SubAggregation("filemd5", agg_filemd5)

	agg_ptype = elastic.NewTermsAggregation().
		Field("Record.PersistenceType.keyword").
		Size(sc.SearchLimit).
		OrderByCount(sc.SortDesc).
		SubAggregation("path", agg_path)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32scripting-persistence"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName), elastic.NewWildcardQuery("Record.PersistenceType.keyword", sc.Type)).
			MustNot(elastic.NewTermQuery("Record.IsGoodService.keyword", "true"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32scripting-persistence"), elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName), elastic.NewWildcardQuery("Record.PersistenceType.keyword", sc.Type))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("ptype", agg_ptype).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	p, found = sr.Aggregations.Terms("ptype")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	PrintQueryObject(query)

	for _, ptypeBucket := range p.Buckets {
		ptype, found := ptypeBucket.Aggregations.Terms("path")

		if found {
			for _, pathBucket := range ptype.Buckets {
				path, found := pathBucket.Aggregations.Terms("filemd5")

				if found {
					for _, filemd5Bucket := range path.Buckets {

						filemd5, found := filemd5Bucket.Aggregations.Terms("computer")
						s := PersistenceStack{
							Type:         ptypeBucket.Key.(string),
							Path:         pathBucket.Key.(string),
							FileMd5:      filemd5Bucket.Key.(string),
							TypeCount:    ptypeBucket.DocCount,
							PathCount:    pathBucket.DocCount,
							FileMd5Count: filemd5Bucket.DocCount,
						}

						if found {
							for _, computerBucket := range filemd5.Buckets {
								s.ComputerName = append(s.ComputerName, computerBucket.Key.(string))
								/*
									if s.ComputerCount <= 0 {
										s.ComputerCount = computerBucket.DocCount
									}
								*/
							}
						}
						if !sc.ShowComputer {
							s.ComputerName = nil
						}
						stackdata = append(stackdata, s)
					}
				}
			}
		}

	}

	ret = api.HttpSuccessMessage("200", &stackdata, sr.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /stacking/persistence/%s HTTP 200, returned Persistence Stack", sc.Type))
	fmt.Fprintln(w, ret)
}
