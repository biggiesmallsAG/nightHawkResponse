package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/config"
	"time"

	"github.com/gorilla/mux"
	elastic "gopkg.in/olivere/elastic.v5"
)

const esTagType = "tags"

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

type Tag struct {
	Timestamp    string
	CreatedBy    string
	CaseName     string `json:"CaseName,omitempty"`
	ComputerName string `json:"ComputerName,omitempty"`
	Audit        string `json:"Audit,omitempty"`
	DocID        string `json:"DocID,omitempty"`
	TagCategory  string
	TagName      string
}

func loadTagVars(tag *Tag, vars map[string]string) {
	if len(vars) == 0 {
		return
	}
	// Request URI parameters takes priority
	if vars["casename"] != "" {
		tag.CaseName = vars["casename"]
	}

	if vars["case_name"] != "" {
		tag.CaseName = vars["case_name"]
	}

	if vars["case"] != "" {
		tag.CaseName = vars["case"]
	}

	if vars["endpoint"] != "" {
		tag.ComputerName = vars["endpoint"]
	}
	if vars["audit"] != "" {
		tag.Audit = vars["audit"]
	}
	if vars["doc_id"] != "" {
		tag.DocID = vars["doc_id"]
	}

	if vars["tag_category"] != "" {
		tag.TagCategory = vars["tag_category"]
	}

	if vars["tag_category"] == "" {
		tag.TagCategory = "default"
	}

	if vars["tag_name"] != "" {
		tag.TagName = vars["tag_name"]
	}
}

func AddTag(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	var tag Tag

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}
	json.Unmarshal(body, &tag)

	loadTagVars(&tag, mux.Vars(r))

	// Check for mandatory data
	if tag.CaseName == "" {
		api.HttpResponseReturn(w, r, "failed", "Casename is required", nil)
		return
	}

	if tag.ComputerName == "" {
		api.HttpResponseReturn(w, r, "failed", "Computer name is required", nil)
		return
	}

	if tag.Audit == "" {
		api.HttpResponseReturn(w, r, "failed", "Audit is required", nil)
		return
	}

	if tag.DocID == "" {
		api.HttpResponseReturn(w, r, "failed", "DocumentID is required", nil)
		return
	}

	if tag.TagCategory == "" {
		api.HttpResponseReturn(w, r, "failed", "Tag category is required", nil)
		return
	}

	if tag.TagName == "" {
		api.HttpResponseReturn(w, r, "failed", "Tag name is required", nil)
		return
	}

	if tag.CreatedBy == "" {
		api.HttpResponseReturn(w, r, "failed", "Creator is required", nil)
		return
	}

	// Sett Tag creation time
	tag.Timestamp = time.Now().UTC().Format(Layout)

	jsonTag, _ := json.Marshal(tag)
	res, err := client.Index().
		Index(conf.ServerIndex()).
		Type(esTagType).
		BodyJson(string(jsonTag)).
		Do(context.Background())

	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	api.HttpResponseReturn(w, r, "success", "Tag added successfully", res.Id)
}

// GetTag returns tag matching search condition
// api_uri: GET /api/v1/tag/sqhow/{case}
// api_uri: GET /api/v1/tag/show/{case}/{endpoint}
// api_uri: GET /api/v1/tag/show/{case}/{endpoint}/{audit}
// api_uri: GET /api/v1/tag/show/{case}/{endpoint}/{audit}/{doc_id}
func GetTagData(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var tag Tag

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
			return
		}

		json.Unmarshal(body, &tag)
	}

	loadTagVars(&tag, mux.Vars(r))

	var queries []elastic.Query
	if tag.CaseName == "" {
		queries = append(queries, elastic.NewWildcardQuery("CaseName.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("CaseName.keyword", tag.CaseName))
	}

	if tag.ComputerName == "" {
		queries = append(queries, elastic.NewWildcardQuery("ComputerName.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("ComputerName.keyword", tag.ComputerName))
	}

	if tag.Audit == "" {
		queries = append(queries, elastic.NewWildcardQuery("Audit.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("Audit.keyword", tag.Audit))
	}

	if tag.DocID == "" {
		queries = append(queries, elastic.NewWildcardQuery("DocID.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("DocID.keyword", tag.DocID))
	}

	if tag.CreatedBy != "" {
		queries = append(queries, elastic.NewTermQuery("CreatedBy.keyword", tag.CreatedBy))
	}

	if tag.TagName != "" {
		queries = append(queries, elastic.NewWildcardQuery("TagName.keyword", fmt.Sprintf("*%s*", tag.TagName)))
	}

	boolquery := elastic.NewBoolQuery().Must(queries...)

	// Print JsonQuery
	// Uncomment the code below to JsonQuery to Elasticsearch
	/*
		boolQueryMap, _ := boolquery.Source()
		jsonBoolQuery, _ := json.Marshal(boolQueryMap)
		fmt.Println(string(jsonBoolQuery))
	*/

	sr, err := client.Search().
		Index(conf.ServerIndex()).
		Type(esTagType).
		Query(boolquery).
		Do(context.Background())

	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	if sr.TotalHits() < 1 {
		api.HttpResponseReturn(w, r, "failed", "No result found", nil)
		return
	}

	// Populating array of tags
	var tags []Tag
	for _, hit := range sr.Hits.Hits {
		var tag Tag
		json.Unmarshal(*hit.Source, &tag)
		tags = append(tags, tag)
	}

	api.HttpResponseReturn(w, r, "success", "Tag search completed", tags)
}

// GetTagCaseAndComputer provides list of CaseName and
// list of ComputerName
// GET /api/v1/comment/cases-and-computers
func GetTagCaseAndComputer(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var casenameList []string
	var computernameList []string

	data := make(map[string]interface{})

	aggCasename := elastic.NewTermsAggregation().
		Field("CaseName.keyword")
	sr, err := client.Search().
		Index(conf.ServerIndex()).Type(esTagType).
		Aggregation("agg_casename", aggCasename).
		Do(context.Background())
	if err != nil {
		api.LogError(api.DEBUG, err)
	}
	ao, found := sr.Aggregations.Terms("agg_casename")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Aggregation not found")
	}

	for _, casenameBucket := range ao.Buckets {
		casenameList = append(casenameList, casenameBucket.Key.(string))
	}

	aggComputer := elastic.NewTermsAggregation().Field("ComputerName.keyword")
	sr, err = client.Search().Index(conf.ServerIndex()).Type(esTagType).Aggregation("agg_computer", aggComputer).Do(context.Background())
	if err != nil {
		api.LogError(api.DEBUG, err)
	}
	ao, found = sr.Aggregations.Terms("agg_computer")
	if !found {
		api.LogDebug(api.DEBUG, "[!] Computer aggregation not found")
	}

	for _, computerBucket := range ao.Buckets {
		computernameList = append(computernameList, computerBucket.Key.(string))
	}

	data["casenames"] = casenameList
	data["endpoints"] = computernameList

	api.HttpResponseReturn(w, r, "success", "Tag CaseNames and ComputerNames", data)
}
