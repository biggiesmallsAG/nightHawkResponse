package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	elastic "gopkg.in/olivere/elastic.v5"

	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/config"
)

// Constant
const (
	IndexName  = "nighthawk"
	BlTable    = "blacklist"
	WlTable    = "whitelist"
	StackTable = "stack"
)

var (
	conf   *config.ConfigVars
	err    error
	client *elastic.Client
	query  elastic.Query
)

func init() {
	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to initialize config read")
		return
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to initialize elastic client")
		return
	}
}

type AnalyzeItem struct {
	Title                  string `json:"Title,omitempty"`                  // common name for the item
	Description            string `json:"Description,omitempty"`            // description about the blacklist/whitelist/stack item
	AuditType              string `json:"AuditType"`                        // audit type
	Name                   string `json:"Name,omitempty"`                   // Record.Name
	Path                   string `json:"Path,omitempty"`                   // Record.Path
	Md5sum                 string `json:"Md5,omitempty"`                    // Record.MD5sum - only applicable to w32apifiles and w32rawfiles
	Arguments              string `json:"Arguments,omitempty"`              // Record.Arguments
	RegPath                string `json:"RegPath,omitempty"`                // Record.StackPath OR Record.RegPath - only applicable to registry based items
	PersistenceType        string `json:"PersistenceType,omitempty"`        // Record.PersistenceType - Only applicable to w32scripting-persistences
	ServiceDescriptiveName string `json:"ServiceDescriptiveName,omitempty"` // Record.DescriptiveName - Only applicable to w32services
	TaskCreator            string `json:"TaskCreator,omitempty"`            // Record.Creator - only applicable to w32tasks
}

func AddBlacklistInformation(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to POST /analyze/add/blacklist")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	entry, err := client.Index().Index(IndexName).Type(BlTable).BodyJson(string(body)).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to add blacklist entry")
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to write to index "+err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, "Blacklist entry successfully added with _id "+entry.Id)
	fmt.Fprintf(w, api.HttpSuccessMessage("200", entry.Id, 1))
}

func AddWhitelistInformation(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to POST /analyze/add/whitelist")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	entry, err := client.Index().Index(IndexName).Type(WlTable).BodyJson(string(body)).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to add whitelist entry")
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to write to index "+err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, "Whitelist entry successfully added with _id "+entry.Id)
	fmt.Fprintf(w, api.HttpSuccessMessage("200", entry.Id, 1))
}

func AddStackInformation(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to POST /analyze/stack")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	entry, err := client.Index().Index(IndexName).Type(StackTable).BodyJson(string(body)).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to add StackCommonItem entry")
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to write to index "+err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, "StackCommonItem entry successfully added with _id "+entry.Id)
	fmt.Fprintf(w, api.HttpSuccessMessage("200", entry.Id, 1))
}

func DeleteAnalyzerItemByID(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	AnalyzerType := vars["analyzer_type"]
	AnalyzerID := vars["analyzer_id"]

	res, err := client.Delete().Index(IndexName).Type(AnalyzerType).Id(AnalyzerID).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, fmt.Sprintf("Failed to delete %s entry with Id %s", AnalyzerType, AnalyzerID))
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to delete analyzer entry, "+err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, fmt.Sprintf("%s entry %s is successfully deleted", AnalyzerType, res.Id))
	fmt.Fprintf(w, api.HttpSuccessMessage("200", res.Id, 1))
}

func DeleteAnalyzerItemByQuery(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	AnalyzerType := vars["analyzer_type"]

	var item AnalyzeItem

	if r.ContentLength > 0 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, fmt.Sprintf("Failed to POST /analyze/delete/%s", AnalyzerType))
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}
		json.Unmarshal(body, &item)
	}

	queries := generateAnalyzeItemElasticQuery(item)
	boolquery := elastic.NewBoolQuery().Must(queries...)

	/*
		// Code snippet to show Elasticsearch query
		boolQueryMap, _ := boolquery.Source()
		jsonBoolQuery, _ := json.Marshal(boolQueryMap)
		fmt.Println(string(jsonBoolQuery))
	*/

	res, err := client.DeleteByQuery().Index(IndexName).Type(AnalyzerType).Query(boolquery).Do(context.TODO())
	if err != nil {
		api.LogDebug(api.DEBUG, fmt.Sprintf("Failed to delete %s entry", AnalyzerType))
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to delete analyzer entry "+err.Error()))
		return
	}

	if res.Deleted > 0 {
		api.LogDebug(api.DEBUG, fmt.Sprintf("%s entry is successfully deleted", AnalyzerType))
		fmt.Fprintf(w, api.HttpSuccessMessage("200", AnalyzerType, res.Deleted))
	} else {
		api.LogDebug(api.DEBUG, fmt.Sprintf("%s entry is not deleted", AnalyzerType))
		fmt.Fprintf(w, api.HttpFailureMessage("Delete unsuccessfull"))
	}

}

// ShowAnalyzeItemByType shows all AnalyzeItem of type blacklist
// whitelist and stack.
// api_uri: GET /api/v1/analyze/show/{analyzer_type}
func ShowAnalyzerItemByType(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	analyzerType := vars["analyzer_type"]

	var item AnalyzeItem

	if r.Method == "POST" && r.ContentLength > 0 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, fmt.Sprintf("Failed to POST /analyze/show/%s", analyzerType))
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}

		json.Unmarshal(body, &item)
	}

	queries := generateAnalyzeItemElasticQuery(item)
	boolquery := elastic.NewBoolQuery().Must(queries...)
	sr, err := client.Search().
		Index(IndexName).
		Type(analyzerType).
		Query(boolquery).
		Do(context.Background())

	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	if sr.TotalHits() < 1 {
		api.HttpResponseReturn(w, r, "failed", fmt.Sprintf("%s entries not found", analyzerType), nil)
		return
	}

	var items []AnalyzeItem
	for _, hit := range sr.Hits.Hits {
		var hitItem AnalyzeItem
		json.Unmarshal(*hit.Source, &hitItem)
		items = append(items, hitItem)
	}
	api.HttpResponseReturn(w, r, "success", fmt.Sprintf("%s entries returned", analyzerType), items)
}

func generateAnalyzeItemElasticQuery(item AnalyzeItem) []elastic.Query {
	var queries []elastic.Query

	if item.AuditType == "" {
		queries = append(queries, elastic.NewWildcardQuery("AuditType", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("AuditType", item.AuditType))
	}

	if item.Title != "" {
		queries = append(queries, elastic.NewTermQuery("Title", item.Title))
	}
	if item.Description != "" {
		queries = append(queries, elastic.NewTermQuery("Description", item.Description))
	}
	if item.Name != "" {
		queries = append(queries, elastic.NewTermQuery("Name", item.Name))
	}
	if item.Path != "" {
		queries = append(queries, elastic.NewTermQuery("Path.keyword", item.Path))
	}
	if item.ServiceDescriptiveName != "" {
		queries = append(queries, elastic.NewTermQuery("ServiceDescriptiveName", item.ServiceDescriptiveName))
	}

	return queries
}
