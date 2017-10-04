package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

const Layout = "2006-01-02T15:04:05Z"

// TimelineConfig is a structure for timeline search
// configuration option.
// This structure provides user flexibility to perform
// timeline search
type TimelineConfig struct {
	CaseName     string   `json:"case_name"`
	Type         string   `json:"type"` // Option to timeline serach on audit_generator
	SearchLimit  int      `json:"search_limit"`
	SortDesc     bool     `json:"sort_desc"`
	IgnoreGood   bool     `json:"ignore_good"`
	Endpoint     string   `json:"endpoint"`
	ComputerList []string `json:"endpoint_list"` // If provided list of computers perform timeline search on computer
	StartTime    string   `json:"start_time"`
	EndTime      string   `json:"end_time"`
	TimeDelta    int      `json:"time_delta"` // Filter result by StartTime += TimeDelta
}

// Default sets default values to TimelineConfig object
// Default values are:
// SearchLimit = 500,
// SortDesc = true,
// IgnoreGood = false
func (tc *TimelineConfig) Default() {
	tc.CaseName = "*"     // search all casese
	tc.Type = "*"         // search all audit_generator
	tc.SearchLimit = 500  // show only first 500 entries matching condition in specified sort order
	tc.SortDesc = false   // default sort the result in ascending order by TlnTime
	tc.IgnoreGood = false // show all records in timeline output
	tc.Endpoint = ""      // default: no endpoint
	tc.ComputerList = nil // default search all computers. If computer names provided only timeline the specified
	// computers in tc.ComputerList

	// StartTime and EndTime field can be either "" else has to match the
	// timestamp layout YYYY-MM-DDTHH:MM:SSZ
	tc.StartTime = "1970-01-01T00:00:01Z" // default start date if not provided
	tc.EndTime = "2030-01-01T00:00:01Z"   // default end date for filter
	tc.TimeDelta = 0                      // Only use this field if TimeDelta value is greater than 0.
	// If TimeDelta is provided, it takes precedence over EndTime
}

// Update TimelineConfig with API provided data
func (tc *TimelineConfig) LoadParams(data []byte) {
	tc.Default()

	var ttc TimelineConfig
	json.Unmarshal(data, &ttc)

	if ttc.CaseName != "" {
		tc.CaseName = ttc.CaseName
	}
	if ttc.Type != "" {
		tc.Type = ttc.Type
	}
	if ttc.SearchLimit > 0 {
		tc.SearchLimit = ttc.SearchLimit
	}
	if ttc.SortDesc {
		tc.SortDesc = true
	}
	if ttc.IgnoreGood {
		tc.IgnoreGood = true
	}
	if ttc.ComputerList != nil {
		tc.ComputerList = ttc.ComputerList
	}

	// If Endpoint is speicified append it to ComputerList
	if ttc.Endpoint != "" {
		tc.Endpoint = ttc.Endpoint
		tc.ComputerList = append(tc.ComputerList, tc.Endpoint)
	}

	if ttc.StartTime != "" {
		tc.StartTime = ttc.StartTime
	}
	if ttc.EndTime != "" {
		tc.EndTime = ttc.EndTime
	}
	if ttc.TimeDelta > 0 {
		tc.TimeDelta = ttc.TimeDelta
		meanTime, _ := time.Parse(Layout, ttc.StartTime)
		startTime := meanTime.Add(time.Duration(-1*ttc.TimeDelta) * time.Minute)
		endTime := meanTime.Add(time.Duration(ttc.TimeDelta) * time.Minute)

		//fmt.Printf("MeanTime: %s, StartTime: %s, EndTime: %s\n\n", meanTime.Format(Layout), startTime.Format(Layout), endTime.Format(Layout))

		tc.StartTime = startTime.Format(Layout)
		tc.EndTime = endTime.Format(Layout)
	}
}

// Handler function for TimelineSearch
func GetTimelineSearch(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json, charset=UTF-8")

	var (
		tc      TimelineConfig
		conf    *config.ConfigVars
		client  *elastic.Client
		query   elastic.Query
		sr      *elastic.SearchResult
		tl_data []SearchOutput
		ret     string
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "[+] GET /search/timeline, error encountered")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	tc.LoadParams(body)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	// Build Elasticsearch query based on Endpoint or ComputerList
	if len(tc.ComputerList) > 0 {
		cl := make([]interface{}, len(tc.ComputerList))
		for i, v := range tc.ComputerList {
			cl[i] = v
		}

		// Run this query if endpoint or endpoint_list is provided
		query = elastic.NewBoolQuery().
			Must(elastic.NewWildcardQuery("CaseInfo.case_name", tc.CaseName),
				elastic.NewWildcardQuery("AuditType.Generator", tc.Type),
				elastic.NewTermsQuery("ComputerName.keyword", cl...)).
			Filter(elastic.NewRangeQuery("Record.TlnTime").From(tc.StartTime).To(tc.EndTime))
	} else {
		// If Endpoint or ComputerList is not provided run this Elastic query
		query = elastic.NewBoolQuery().
			Must(elastic.NewWildcardQuery("CaseInfo.case_name", tc.CaseName),
				elastic.NewWildcardQuery("AuditType.Generator", tc.Type)).
			Filter(elastic.NewRangeQuery("Record.TlnTime").From(tc.StartTime).To(tc.EndTime))
	}

	sr, err = client.Search().
		Index(conf.ServerIndex()).
		Query(query).
		Sort("Record.TlnTime", true).
		Size(tc.SearchLimit).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if sr.TotalHits() < 1 {
		ret = api.HttpSuccessMessage("No hits", make([]*SearchOutput, 0), 0)
		api.LogDebug(api.DEBUG, "[+] GET /search/timeline HTTP 200, returned no hits for timeline")
		fmt.Fprintf(w, string(ret))
		return
	}

	for _, hit := range sr.Hits.Hits {
		var so SearchOutput
		so.Unmarshal(hit, false)
		tl_data = append(tl_data, so)
	}

	ret = api.HttpSuccessMessage("200", &tl_data, sr.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /search/timeline HTTP 200, returned timeline")
	fmt.Fprintf(w, string(ret))
}
