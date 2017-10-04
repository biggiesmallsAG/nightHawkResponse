// Package delete provides functions to delete case, endpoints and documents
package delete

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

	"github.com/gorilla/mux"
	elastic "gopkg.in/olivere/elastic.v5"
)

// DeleteConfig is API configuration option for delete operation
type DeleteConfig struct {
	CaseName     string   `json:"case_name"`
	Endpoint     string   `json:"endpoint"`
	EndpointList []string `json:"endpoint_list"`
}

// Default sets default values to DeleteConfig
// CaseName: *          - Delete from all cases
// Endpoint: NULL       - User must provide
// EndpointList: nil    - Optional. Use to delete multiple endpoints
func (config *DeleteConfig) Default() {
	config.CaseName = "*"     // default: delete from all cases
	config.Endpoint = ""      // default detele single endpoint
	config.EndpointList = nil // Used if need to delete more than one endpoint
}

// LoadParams maps post data to DeleteConfig object.
// (config *DeleteConfig) Default() is called to initialize
// DeleteConfig object before parsing post data.
func (config *DeleteConfig) LoadParams(data []byte) {
	config.Default()

	var tconfig DeleteConfig
	json.Unmarshal(data, &tconfig)

	if tconfig.CaseName != "" {
		config.CaseName = tconfig.CaseName
	}
	if tconfig.EndpointList != nil {
		config.EndpointList = tconfig.EndpointList
	}

	// Endpoint must be processed only after EndpointList has been processed
	// to avoid removal of previously appended entry by assignment operation
	if tconfig.Endpoint != "" {
		config.Endpoint = tconfig.Endpoint
		config.EndpointList = append(config.EndpointList, config.Endpoint)
	}
}

// DeletecCaseEndpoint function deletes endpoint for a given case.
// This function is invoked when case and endpoint are passed as
// /api/v1/delete/{case_name}/{endpoint_name}
//
func DeleteCaseEndpoint(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json, charset=UTF-8")
	var dc DeleteConfig

	if r.Method != "GET" {
		api.LogDebug(api.DEBUG, fmt.Sprintf("[+] %s %s - Invalid HTTP Method", r.Method, r.URL.Path))
		fmt.Fprintln(w, api.HttpFailureMessage("Invalid HTTP request method"))
		return
	}

	params := mux.Vars(r)
	dc.Default()
	dc.CaseName = params["case_name"]
	dc.Endpoint = params["endpoint_name"]

	var query elastic.Query
	query = elastic.NewBoolQuery().
		Must(elastic.NewWildcardQuery("CaseInfo.case_name", dc.CaseName),
			elastic.NewTermQuery("ComputerName.keyword", dc.Endpoint))

	deleteEndpointByQuery(w, query, "DeleteCaseEndpoint")
}

// DeleteCase deletes case and all the endpoints in a given case.
// A valid case_name must be provided in HTTP POST request.
func DeleteCase(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json, charset=UTF-8")
	var dc DeleteConfig

	switch r.Method {
	case "GET":
		params := mux.Vars(r)
		dc.Default()
		dc.CaseName = params["case_name"]
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, "[+] POST /delete/endpoint, failed to read request")
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}
		dc.LoadParams(body)
	}

	if dc.CaseName == "*" || dc.CaseName == "" {
		api.LogDebug(api.DEBUG, "[+] POST /delete/case, valid casename required")
		fmt.Fprintln(w, api.HttpFailureMessage("Valid casename is required. * or NULL can not be used"))
		return
	}

	var query elastic.Query
	query = elastic.NewBoolQuery().Must(elastic.NewTermQuery("CaseInfo.case_name", dc.CaseName))
	deleteEndpointByQuery(w, query, "DeleteCase")
}

// DeleteEndpoint deletes endpoint record from case.
// If CaseName is not specified, it will delete endpoint from
// all the cases in the given index.
func DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json, charset=UTF-8")
	var dc DeleteConfig

	switch r.Method {
	case "GET":
		params := mux.Vars(r)
		dc.Default()
		dc.Endpoint = params["endpoint_name"]
		dc.EndpointList = append(dc.EndpointList, dc.Endpoint)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.LogDebug(api.DEBUG, "[+] POST /delete/endpoint, failed to read request")
			fmt.Fprintln(w, api.HttpFailureMessage("Failed to read HTTP request"))
			return
		}
		dc.LoadParams(body)
	}

	// Verify Endpoint is provided in request body
	if len(dc.EndpointList) == 0 {
		api.LogDebug(api.DEBUG, "[+] POST /delete/endpoint - endpoint is required")
		fmt.Fprintln(w, api.HttpFailureMessage("Endpoint required"))
		return
	}

	var query elastic.Query

	// Convert EndpointList to interface{} slice
	endpoints := make([]interface{}, len(dc.EndpointList))
	for i, v := range dc.EndpointList {
		endpoints[i] = v
	}

	query = elastic.NewBoolQuery().
		Must(elastic.NewWildcardQuery("CaseInfo.case_name", dc.CaseName),
			elastic.NewTermsQuery("ComputerName.keyword", endpoints...))

	deleteEndpointByQuery(w, query, "DeleteEndpoint")

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

// deleteGlobalEndpoint is private function that deletes endpoint from
// elasticsearch server
func deleteEndpointByQuery(w http.ResponseWriter, query elastic.Query, caller string) {
	fmt.Println("Deleting endpoint: ", caller)

	client, err := getElasticsearchClient()
	if err != nil {
		api.LogError(api.DEBUG, err)
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to initialized query. ERROR 0x40000012"))
		return
	}

	res, err := client.DeleteByQuery().
		Index(config.ElasticIndex()).
		Type("audit_type").
		Query(query).
		Do(context.TODO())

	if err != nil {
		fmt.Fprintf(w, api.HttpFailureMessage("Failed to delete endpoint"))
		return
	}
	json_res, _ := json.MarshalIndent(&res, " ", "  ")

	// fmt.Fprintf(w, api.HttpSuccessMessage(fmt.Sprintf("%s: Endpoint deleted\n", caller)))
	fmt.Fprintf(w, string(json_res))
}
