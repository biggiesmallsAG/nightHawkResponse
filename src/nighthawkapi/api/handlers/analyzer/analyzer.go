package analyzer

import (
	"encoding/json"
	"context"
	"fmt"
	"net/http"
	"io/ioutil"

	elastic "gopkg.in/olivere/elastic.v5"
	"github.com/gorilla/mux"

	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"
	
)


// Constant
const (
	IndexName = "nighthawk"
	BlTable = "blacklist"
	WlTable = "whitelist"
	StackTable = "stack"
)


type AnalyzeItem struct {
	Title 					string `json:"title,omitempty"`						// common name for the item
	Description 			string `json:"description,omitempty"`				// description about the blacklist/whitelist/stack item
	AuditType 				string `json:"audit_type"`							// audit type
	Name 					string `json:"name,omitempty"`						// Record.Name 
	Path					string `json:"path,omitempty"`						// Record.Path
	Md5sum 					string `json:"md5,omitempty"`						// Record.MD5sum - only applicable to w32apifiles and w32rawfiles
	Arguments				string	`json:"arguments,omitempty"`				// Record.Arguments 
	RegPath					string `json:"reg_path,omitempty"`					// Record.StackPath OR Record.RegPath - only applicable to registry based items
	PersistenceType			string `json:"persistence_type,omitempty"`			// Record.PersistenceType - Only applicable to w32scripting-persistences
	ServiceDescriptiveName 	string `json:"service_descriptive_name,omitempty"`	// Record.DescriptiveName - Only applicable to w32services
	TaskCreator				string `json:"task_creator,omitempty"`				// Record.Creator - only applicable to w32tasks
}


func AddBlacklistInformation(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to POST /analyze/blacklist")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	conf, err := config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	entry, err := client.Index().Index(IndexName).Type(BlTable).BodyJson(string(body)).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to add blacklist entry")
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to write to index " + err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, "Blacklist entry successfully added with _id " + entry.Id)
	fmt.Fprintf(w, api.HttpSuccessMessage("200", entry.Id, 1))
}


func AddWhitelistInformation(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to POST /analyze/whitelist")
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}

	conf, err := config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	entry, err := client.Index().Index(IndexName).Type(WlTable).BodyJson(string(body)).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to add whitelist entry")
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to write to index " + err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, "Whitelist entry successfully added with _id " + entry.Id)
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

	conf, err := config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	entry, err := client.Index().Index(IndexName).Type(StackTable).BodyJson(string(body)).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to add StackCommonItem entry")
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to write to index " + err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, "StackCommonItem entry successfully added with _id " + entry.Id)
	fmt.Fprintf(w, api.HttpSuccessMessage("200", entry.Id, 1))
}



func DeleteAnalyzerItemByID(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	AnalyzerType := vars["analyzer_type"]
	AnalyzerID := vars["analyzer_id"]

	conf, err := config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	res, err := client.Delete().Index(IndexName).Type(AnalyzerType).Id(AnalyzerID).Do(context.Background())
	if err != nil {
		api.LogDebug(api.DEBUG, fmt.Sprintf("Failed to delete %s entry with Id %s", AnalyzerType, AnalyzerID))
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to delete analyzer entry, " + err.Error()))
		return
	}

	api.LogDebug(api.DEBUG, fmt.Sprintf("%s entry %s is successfully deleted", AnalyzerType, res.Id))
	fmt.Fprintf(w, api.HttpSuccessMessage("200", res.Id, 1))
}





func DeleteAnalyzerItemByQuery(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	AnalyzerType := vars["analyzer_type"]

	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.LogDebug(api.DEBUG, fmt.Sprintf("Failed to POST /analyze/delete/%s", AnalyzerType))
		fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
		return
	}
	

	conf, err := config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}


	// Constructing BoolQuery
	var item AnalyzeItem
	json.Unmarshal(body, &item)
	fmt.Println(item)
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("audit_type", item.AuditType))
	if item.Title != "" {
		queries = append(queries, elastic.NewTermQuery("title", item.Title))
	}
	if item.Description != "" {
		queries = append(queries, elastic.NewTermQuery("description", item.Description))
	}
	if item.Name != "" {
		queries = append(queries, elastic.NewTermQuery("name", item.Name))
	}
	if item.Path != "" {
		queries = append(queries, elastic.NewTermQuery("path.keyword", item.Path))
	}
	if item.ServiceDescriptiveName != "" {
		queries = append(queries, elastic.NewTermQuery("service_descriptive_name", item.ServiceDescriptiveName))
	}

	boolquery := elastic.NewBoolQuery().Must(queries...)

	boolQueryMap,_ := boolquery.Source()
	jsonBoolQuery,_ := json.Marshal(boolQueryMap)
	fmt.Println(string(jsonBoolQuery))

	res, err := client.DeleteByQuery().Index(IndexName).Type(AnalyzerType).Query(boolquery).Do(context.TODO())
	if err != nil {
		api.LogDebug(api.DEBUG, fmt.Sprintf("Failed to delete %s entry", AnalyzerType))
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to delete analyzer entry " + err.Error()))
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






