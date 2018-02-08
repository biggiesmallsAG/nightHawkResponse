package stacking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/config"
	gsapi "nighthawkresponse/api/handlers/search"

	"github.com/gorilla/mux"
	elastic "gopkg.in/olivere/elastic.v5"
)

type DiffConfig struct {
	Endpoint     string `json:"endpoint"`
	SearchLimit  int    `json:"search-limit"`
	SubAggSize   int    `json:"sub_agg_size"`
	dRunKey      bool   `json:"diff_runkey"`
	dPersistence bool   `json:"diff_persistence"`
	dTasks       bool   `json:"diff_tasks"`
	dServices    bool   `json:"diff_services"`
	dLocalPorts  bool   `json:"diff_local_port"` // Ports on listening mode only
	dShimCache   bool   `json:"diff_shimcache"`
}

func (config *DiffConfig) Default() {
	config.Endpoint = ""
	config.SearchLimit = 100
	config.SubAggSize = 10
	config.dRunKey = true
	config.dPersistence = true
	config.dTasks = true
	config.dServices = true
	config.dLocalPorts = true
	config.dShimCache = false
}

func (config *DiffConfig) LoadParams(data []byte) {
	config.Default()

	var tconfig DiffConfig
	json.Unmarshal(data, &tconfig)

	if tconfig.Endpoint != "" {
		config.Endpoint = tconfig.Endpoint
	}
	if tconfig.SearchLimit > 0 {
		config.SearchLimit = tconfig.SearchLimit
	}
	if tconfig.SubAggSize > 0 {
		config.SubAggSize = tconfig.SubAggSize
	}
}

type DiffData struct {
	DocId       string `json:"doc_id"`
	ChangedDate string `json:"changed_date"`
	ChangedType string `json:"changed_type,omitempty"` // installed || removed
	CaseName    string `json:"case_name"`
	AuditType   string `json:"audit_type"`
	ChangedData string `json:"changed_data"`
	Description string `json:"description"`
}

func TimelineEndpointDiff(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json, charset=UTF-8")
	var dc DiffConfig

	switch r.Method {
	case "GET":
		params := mux.Vars(r)
		dc.Default()
		dc.Endpoint = params["endpoint"]
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, "[+] POST /diff, failed to read request")
			fmt.Fprintf(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}
		dc.LoadParams(body)
	}

	numofcase, casedata := NumOfCaseByEndpoint(dc)

	var data []DiffData

	if dc.dRunKey {
		diffRunKey(dc, numofcase, casedata, &data)
	}
	if dc.dPersistence {
		diffPersistence(dc, numofcase, casedata, &data)
	}
	if dc.dTasks {
		diffTasks(dc, numofcase, casedata, &data)
	}
	if dc.dServices {
		diffServices(dc, numofcase, casedata, &data)
	}
	if dc.dLocalPorts {
		diffLocalPorts(dc, numofcase, casedata, &data)
	}
	if dc.dShimCache {
		diffShimCache(dc, numofcase, casedata, &data)
	}

	ret, _ := json.MarshalIndent(&data, "", " ")
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] %s /diff/%s", r.Method, dc.Endpoint))
	fmt.Fprintf(w, string(ret))
}

// getElasticsearchClient returns *elastic.Client
func getElasticsearchClient() (*elastic.Client, error) {
	server_url := config.ElasticUrl()
	if server_url == "" {
		return nil, errors.New("getElasticsearchClient - ERROR - Failed to get valid elasticsearch URL")
	}

	client, err := elastic.NewClient(elastic.SetURL(config.ElasticUrl()))
	if err != nil {
		return nil, err
	}

	return client, nil
}

type CaseData struct {
	CaseOrder        int
	CaseName         string
	DocCount         int64
	TriageDate       float64
	TriageDateString string
}

// NumOfCaseByEndpoint calculates number of cases for given endpoint
// and returns number of cases and case details
func NumOfCaseByEndpoint(configopt DiffConfig) (int, []CaseData) {
	// setting null
	numofcase := 0
	var casedata []CaseData

	var query elastic.Query
	agg_triage_date := elastic.NewTermsAggregation().
		Field("Record.JobCreated").
		Size(1)

	agg_case := elastic.NewTermsAggregation().
		Field("CaseInfo.CaseName.keyword").
		Size(configopt.SubAggSize).
		SubAggregation("triage_date", agg_triage_date)

	query = elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("ComputerName.keyword", configopt.Endpoint))

	client, err := getElasticsearchClient()
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Printf("NumOfCaseByEndpoint - ERROR - Failed to initialize Elasticserach client\n")
		return 0, nil
	}

	sr, err := client.Search().
		Index(config.ElasticIndex()).
		Type("audit_type").
		Query(query).
		Size(0).
		Aggregation("case_name", agg_case).
		Do(context.Background())
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Printf("NumOfCaseByEndpoint - ERROR - Failed to run Elasticsearch Search query\n")
		return 0, nil
	}

	cn, found := sr.Aggregations.Terms("case_name")
	if !found {
		api.LogDebug(api.DEBUG, fmt.Sprintf("[!] CaseName Aggregation not found for %s", configopt.Endpoint))
		fmt.Printf("NumOfCaseByEndpoint - ERROR - CaseName aggregation not found for %s\n", configopt.Endpoint)
		return 0, nil
	}

	numofcase = len(cn.Buckets)

	for _, cnBucket := range cn.Buckets {
		td, found := cnBucket.Aggregations.Terms("triage_date")
		if found {
			//fmt.Println(td.Buckets[0].Key)
			//fmt.Println(*td.Buckets[0].KeyAsString)
			//fmt.Println(td.Buckets[0].DocCount)

			cd := CaseData{
				CaseName:         cnBucket.Key.(string),
				DocCount:         cnBucket.DocCount,
				TriageDate:       td.Buckets[0].Key.(float64),
				TriageDateString: *td.Buckets[0].KeyAsString,
			}

			casedata = append(casedata, cd)
		}
	}

	//fmt.Println(casedata)
	// default return
	return numofcase, casedata
}

// getTriageDateStringForCase returns the triage collected date from []CaseData
func getTriageDateStringForCase(casedata []CaseData, casename string) string {
	for _, cd := range casedata {
		if cd.CaseName == casename {
			return cd.TriageDateString
		}
	}
	// Default should never be returned
	return ""
}

// getTimelineEventDescription returns Elasticsearch document_id and audit event
// description
func getTimelineEventDescription(audit_type string, audit_key string, audit_value string, casename string, endpoint string) (docid string, description string) {
	client, err := getElasticsearchClient()
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Printf("diffRunKey - ERROR - Failed to initialize Elasticserach client\n")
	}

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("CaseInfo.CaseName.keyword", casename),
			elastic.NewTermQuery("AuditType.Generator.keyword", audit_type),
			elastic.NewTermQuery("ComputerName.keyword", endpoint),
			elastic.NewTermQuery(audit_key, audit_value))

	sr, err := client.Search().
		Index(config.ElasticIndex()).
		Query(query).
		Do(context.Background())
	if err != nil {
		api.LogError(api.DEBUG, errors.New(fmt.Sprintf("getTimelineEventDescription - ERROR - Failed to search Elasticsearch")))
		fmt.Printf("getTimelineEventDescription - ERROR - Failed to search Elasticsearch")
	}

	hit := sr.Hits.Hits[0]
	var so gsapi.SearchOutput
	so.Unmarshal(hit, true)
	return so.Id, so.Summary
}

func diffPersistence(dc DiffConfig, numofcase int, casedata []CaseData, data *[]DiffData) {
	getTimelineDataByAudit("w32persistence", "Record.Path.keyword", dc, numofcase, casedata, data)
}

func diffTasks(dc DiffConfig, numofcase int, casedata []CaseData, data *[]DiffData) {
	getTimelineDataByAudit("w32tasks", "Record.Name.keyword", dc, numofcase, casedata, data)
}

func diffServices(dc DiffConfig, numofcase int, casedata []CaseData, data *[]DiffData) {
	getTimelineDataByAudit("w32services", "Record.Name.keyword", dc, numofcase, casedata, data)
}

func diffRunKey(dc DiffConfig, numofcase int, casedata []CaseData, data *[]DiffData) {
	client, err := getElasticsearchClient()
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Printf("diffRunKey - ERROR - Failed to initialize Elasticserach client\n")
	}

	agg_cn := elastic.NewTermsAggregation().
		Field("CaseInfo.CaseName.keyword").
		Size(dc.SubAggSize)
	agg_runkey := elastic.NewTermsAggregation().
		Field("Record.StackPath.keyword").
		Size(dc.SearchLimit).
		OrderByCount(true).
		SubAggregation("casename", agg_cn)

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("AuditType.Generator.keyword", "w32scripting-persistence"),
			elastic.NewTermQuery("Record.PersistenceType.keyword", "Registry"),
			elastic.NewWildcardQuery("Record.RegPath.keyword", "*Run*"),
			elastic.NewTermQuery("ComputerName.keyword", dc.Endpoint))

	sr, err := client.Search().
		Index(config.ElasticIndex()).
		Query(query).
		Size(0).
		Aggregation("runkey", agg_runkey).
		Do(context.Background())
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Printf("diffRunKey - ERROR - Failed to run Elasticsearch Search query")
		return
	}

	rk, found := sr.Aggregations.Terms("runkey")
	if !found {
		api.LogDebug(api.DEBUG, fmt.Sprintf("diffRunKey - INFO - Service aggregation not found"))
		fmt.Printf("diffRunKey - INFO - Service aggreation not found")
		return
	}

	for _, rkBucket := range rk.Buckets {
		if rkBucket.DocCount == int64(numofcase) {
			continue
		}

		cn, found := rkBucket.Aggregations.Terms("casename")
		if found {
			for _, cnBucket := range cn.Buckets {
				dd := DiffData{
					ChangedDate: getTriageDateStringForCase(casedata, cnBucket.Key.(string)),
					CaseName:    cnBucket.Key.(string),
					AuditType:   "w32runkey",
					ChangedData: rkBucket.Key.(string),
				}

				id, desc := getTimelineEventDescription("w32scripting-persistence", "Record.StackPath.keyword", dd.ChangedData, dd.CaseName, dc.Endpoint)
				dd.DocId = id
				dd.Description = desc

				//fmt.Println(dd)
				*data = append(*data, dd)
			}
		}
	}
}

func diffLocalPorts(dc DiffConfig, numofcase int, casedata []CaseData, data *[]DiffData) {
	getTimelineDataByAudit("w32ports", "Record.Process.keyword", dc, numofcase, casedata, data)
}

func diffShimCache(dc DiffConfig, numofcase int, casedata []CaseData, data *[]DiffData) {

}

func getTimelineDataByAudit(audit_type string, audit_field string, dc DiffConfig, numofcase int, casedata []CaseData, data *[]DiffData) {
	client, err := getElasticsearchClient()
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Printf("getTimelineDataByAudit:%s - ERROR - Failed to initialize Elasticserach client\n", audit_type)
		return
	}

	agg_cn := elastic.NewTermsAggregation().
		Field("CaseInfo.CaseName.keyword").
		Size(dc.SubAggSize)
	agg := elastic.NewTermsAggregation().
		Field(audit_field).
		Size(dc.SearchLimit).
		OrderByCount(true).
		SubAggregation("casename", agg_cn)

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("AuditType.Generator.keyword", audit_type),
			elastic.NewTermQuery("ComputerName.keyword", dc.Endpoint))

	sr, err := client.Search().
		Index(config.ElasticIndex()).
		Query(query).
		Size(0).
		Aggregation("agg", agg).
		Do(context.Background())
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Printf("getTimelineDataByAudit:%s - ERROR - Failed to run Elasticsearch Search query", audit_type)
		return
	}

	ags, found := sr.Aggregations.Terms("agg")
	if !found {
		api.LogDebug(api.DEBUG, fmt.Sprintf("getTimelineDataByAudit - INFO - Service aggregation not found"))
		fmt.Printf("getTimelineDataByAudit:%s - INFO - Service aggreation not found", audit_type)
		return
	}

	for _, agsBucket := range ags.Buckets {
		if agsBucket.DocCount >= int64(numofcase) {
			continue
		}

		cn, found := agsBucket.Aggregations.Terms("casename")
		if found {
			for _, cnBucket := range cn.Buckets {
				dd := DiffData{
					ChangedDate: getTriageDateStringForCase(casedata, cnBucket.Key.(string)),
					CaseName:    cnBucket.Key.(string),
					AuditType:   audit_type,
					ChangedData: agsBucket.Key.(string),
				}

				id, desc := getTimelineEventDescription(audit_type, audit_field, dd.ChangedData, dd.CaseName, dc.Endpoint)
				dd.DocId = id
				dd.Description = desc

				//fmt.Println(dd)
				*data = append(*data, dd)
			}
		}
	}
}
