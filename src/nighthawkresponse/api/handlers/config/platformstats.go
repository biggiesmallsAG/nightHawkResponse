package config

import (
	"context"
	"fmt"
	"net/http"
	api "nighthawkresponse/api/core"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type PStats struct {
	Memory  *mem.VirtualMemoryStat `json:"memory"`
	Cpu     []cpu.TimesStat        `json:"cpu"`
	Disk    *disk.UsageStat        `json:"disk"`
	Elastic EStats                 `json:"elastic"`
}

type EStats struct {
	Docs      *elastic.ClusterStatsIndicesDocs `json:"docs"`
	Nodes     int                              `json:"nodes"`
	Cases     int                              `json:"cases"`
	Endpoints int                              `json:"endpoints"`
	Cluster   string                           `json:"cluster"`
	Status    string                           `json:"status"`
	Index     string                           `json:"index"`
}

func ReturnTotalCasesAndEndpoints(c *elastic.Client, index string) (cs int, ep int) {
	agg := elastic.NewTermsAggregation().
		Field("CaseInfo.CaseName.keyword").
		Size(5000)

	cases, err := c.Search().
		Index(index).
		Query(elastic.NewMatchAllQuery()).
		Size(0).
		Aggregation("cases", agg).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	a, found := cases.Aggregations.Terms("cases")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	cs = len(a.Buckets)

	agg = elastic.NewTermsAggregation().
		Field("ComputerName.keyword").
		Size(5000)

	endpoints, err := c.Search().
		Index(index).
		Query(elastic.NewMatchAllQuery()).
		Size(0).
		Aggregation("endpoints", agg).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	e, found := endpoints.Aggregations.Terms("endpoints")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found!")
	}

	ep = len(e.Buckets)

	return cs, ep
}

func ReturnPlatformStats(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Times(true)
	d, _ := disk.Usage("/")

	conf, err := ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	var (
		method, ret string
	)

	if conf.Elastic.Elastic_ssl {
		method = api.HTTPS
	} else {
		method = api.HTTP
	}

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s%s:%d", method, conf.Elastic.Elastic_server, conf.Elastic.Elastic_port)))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	e, err := client.ClusterStats().
		Do(context.Background())

	cases, endpoints := ReturnTotalCasesAndEndpoints(client, conf.Elastic.Elastic_index)

	_e := EStats{
		Docs:      e.Indices.Docs,
		Nodes:     e.Nodes.Count.Total,
		Status:    e.Status,
		Cluster:   e.ClusterName,
		Index:     conf.Elastic.Elastic_index,
		Cases:     cases,
		Endpoints: endpoints,
	}

	updated := PStats{
		Memory:  v,
		Cpu:     c,
		Disk:    d,
		Elastic: _e,
	}

	ret = api.HttpSuccessMessage("200", &updated, 0)
	api.LogDebug(api.DEBUG, "[+] GET /platformstats HTTP 200, returned platform statistics.")
	fmt.Fprintln(w, string(ret))
}
