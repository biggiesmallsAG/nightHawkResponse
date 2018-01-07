package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

type GlobalSearchConfig struct {
	CaseName   string `json:"case_name"`   // Filter by CaseName, if required
	Type       string `json:"type"`        // Filter by AuditType, if required
	Endpoint   string `json:"endpoint"`    // Filter by Endpoint, if required
	SearchTerm string `json:"search_term"` // Mandatory field
	SearchSize int    `json:"search_size"`
}

func (config *GlobalSearchConfig) Default() {
	config.CaseName = "*" // Default wildcard match to any case
	config.Type = "*"
	config.Endpoint = "*"   // Default null
	config.SearchTerm = ""  // Default null
	config.SearchSize = 500 // Default only return first 500 entries
}

func (config *GlobalSearchConfig) LoadParams(data []byte) {
	config.Default()

	var tconfig GlobalSearchConfig // Temporary GlobalSearchConfig object
	json.Unmarshal(data, &tconfig)

	if tconfig.CaseName != "" {
		config.CaseName = tconfig.CaseName
	}

	if tconfig.Type != "" {
		config.Type = tconfig.Type
	}

	if tconfig.Endpoint != "" {
		config.Endpoint = tconfig.Endpoint
	}

	if tconfig.SearchTerm != "" {
		config.SearchTerm = tconfig.SearchTerm
	}

	if tconfig.SearchSize > 0 {
		config.SearchSize = tconfig.SearchSize
	}
}

func GetGlobalSearch(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		gsc     GlobalSearchConfig
		gs_data []SearchOutput
		query   elastic.Query
		sr      *elastic.SearchResult
		conf    *config.ConfigVars
		client  *elastic.Client
		ret     string
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] POST /stacking/context/endpoint, Error encountered")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read request body"))
		return
	}

	gsc.LoadParams(body)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	query = elastic.NewBoolQuery().
		Must(elastic.NewMatchQuery("_all", gsc.SearchTerm),
			elastic.NewWildcardQuery("CaseInfo.case_name.keyword", gsc.CaseName),
			elastic.NewWildcardQuery("AuditType.Generator.keyword", gsc.Type),
			elastic.NewWildcardQuery("ComputerName.keyword", gsc.Endpoint))

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Size(gsc.SearchSize).
		Do(context.Background())

	if err != nil {
		api.HttpFailureMessage(fmt.Sprintf("Elasticsearch Error: %s", err.Error()))
		api.LogError(api.DEBUG, err)
		return
	}

	if sr.TotalHits() < 1 {
		api.LogDebug(api.DEBUG, "[+] POST /search HTTP 200, returned no hits.")
		ret = api.HttpSuccessMessage("No hits", make([]*SearchOutput, 0), 0)
		fmt.Fprintln(w, string(ret))
		return
	}

	for _, hit := range sr.Hits.Hits {
		var so SearchOutput
		so.Unmarshal(hit, true)
		gs_data = append(gs_data, so)
	}

	ret = api.HttpSuccessMessage("200", &gs_data, sr.TotalHits())
	api.LogDebug(api.DEBUG, "[+] POST /search HTTP 200, returned search data")
	fmt.Fprintln(w, ret)
}
