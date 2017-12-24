package analyzer

import (
	"context"
	"fmt"
	"net/http"
	"io/ioutil"

	elastic "gopkg.in/olivere/elastic.v5"
	
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
	AuditType 				string `json:"audit_type"`
	Name 					string `json:"name,omitempty"`
	Path					string `json:"path,omitempty"`
	Md5sum 					string `json:"md5,omitempty"`
	Arguments				string	`json:"arguments,omitempty"`
	RegPath					string `json:"reg_path,omitempty"`
	PersistenceType			string `json:"persistence_type,omitempty"`
	ServiceDescriptiveName 	string `json:"service_descriptive_name,omitempty"`
	TaskCreator				string `json:"task_creator,omitempty"`
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




