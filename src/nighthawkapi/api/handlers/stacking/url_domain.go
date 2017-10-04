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

type UrlDomainStack struct {
	Domain      string
	Host        string
	DomainCount int64
	HostCount   int64
}

func StackUrlDomain(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		ret                      string
		sc                       StackingConfig
		conf                     *config.ConfigVars
		err                      error
		client                   *elastic.Client
		query                    elastic.Query
		agg_hostname, agg_domain elastic.Aggregation
		sr                       *elastic.SearchResult
		url                      *elastic.AggregationBucketKeyItems
		stackdata                []UrlDomainStack
		found                    bool
	)
	// Read http parameters
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/url/domain, Error encountered")
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

	agg_hostname = elastic.NewTermsAggregation().
		Field("Record.UrlHostname.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc)

	agg_domain = elastic.NewTermsAggregation().
		Field("Record.UrlDomain.keyword").
		Size(sc.SearchLimit).
		OrderByCount(sc.SortDesc).
		SubAggregation("url_hostname", agg_hostname)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator", "urlhistory"),
				elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName)).
			MustNot(elastic.NewTermQuery("Record.IsGoodEntry", "true"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator", "urlhistory"),
				elastic.NewWildcardQuery("CaseInfo.case_name", sc.CaseName))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("url_domain", agg_domain).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	url, found = sr.Aggregations.Terms("url_domain")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, domainBucket := range url.Buckets {
		hostname, found := domainBucket.Aggregations.Terms("url_hostname")

		if found {
			for _, hostnameBucket := range hostname.Buckets {
				s := UrlDomainStack{
					Domain:      domainBucket.Key.(string),
					Host:        hostnameBucket.Key.(string),
					DomainCount: domainBucket.DocCount,
					HostCount:   hostnameBucket.DocCount,
				}
				stackdata = append(stackdata, s)
			}
		}
	}

	ret = api.HttpSuccessMessage("200", &stackdata, sr.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /stacking/url/domain HTTP 200, returned URL Domain Stack"))
	fmt.Fprintln(w, ret)
}
