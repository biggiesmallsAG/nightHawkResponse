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

type DnsStack struct {
	Host           string
	IpAddress      string
	HostCount      int64
	IpAddressCount int64
}

func StackDnsARequest(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		ret                  string
		sc                   StackingConfig
		conf                 *config.ConfigVars
		err                  error
		client               *elastic.Client
		query                elastic.Query
		agg_ipaddr, agg_host elastic.Aggregation
		sr                   *elastic.SearchResult
		dns                  *elastic.AggregationBucketKeyItems
		stackdata            []DnsStack
		found                bool
	)

	// Read http parameters
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /stacking/dns/a, Error encountered")
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

	agg_ipaddr = elastic.NewTermsAggregation().
		Field("Record.RecordDataList.Ipv4Address.keyword").
		Size(sc.SubAggSize).
		OrderByCount(sc.SortDesc)

	agg_host = elastic.NewTermsAggregation().
		Field("Record.Host.keyword").
		Size(sc.SearchLimit).
		OrderByCount(sc.SortDesc).
		SubAggregation("dns_ipaddr", agg_ipaddr)

	if sc.IgnoreGood {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32network-dns"),
				elastic.NewTermQuery("Record.RecordType.keyword", "A"),
				elastic.NewWildcardQuery("CaseInfo.case_name.keyword", sc.CaseName)).
			MustNot(elastic.NewTermQuery("Record.IsGoodEntry.keyword", "true"),
				elastic.NewWildcardQuery("Record.Host.keyword", "_kerberos._tcp*"),
				elastic.NewWildcardQuery("Record.Host.keyword", "_ldap._tcp*"))
	} else {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32network-dns"),
				elastic.NewTermQuery("Record.RecordType.keyword", "A"),
				elastic.NewWildcardQuery("CaseInfo.case_name.keyword", sc.CaseName)).
			MustNot(elastic.NewWildcardQuery("Record.Host.keyword", "_kerberos._tcp*"),
				elastic.NewWildcardQuery("Record.Host.keyword", "_ldap._tcp*"))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(0).
		Aggregation("dns_host", agg_host).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	dns, found = sr.Aggregations.Terms("dns_host")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, dnsBucket := range dns.Buckets {
		ipaddr, found := dnsBucket.Aggregations.Terms("dns_ipaddr")

		if found {
			for _, ipaddrBucket := range ipaddr.Buckets {
				s := DnsStack{
					Host:           dnsBucket.Key.(string),
					IpAddress:      ipaddrBucket.Key.(string),
					HostCount:      dnsBucket.DocCount,
					IpAddressCount: ipaddrBucket.DocCount,
				}
				stackdata = append(stackdata, s)
			}
		}
	}

	ret = api.HttpSuccessMessage("200", &stackdata, sr.TotalHits())
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] GET /stacking/dns/a HTTP 200, returned DNS Stack"))
	fmt.Fprintln(w, ret)
}
