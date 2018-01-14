package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/auth"
	"nighthawkapi/api/handlers/config"
	"time"

	"github.com/gorilla/mux"
	elastic "gopkg.in/olivere/elastic.v5"
)

const EsTagType = "tags"

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
	Timestamp    string `json:"timestamp"`
	CreatedBy    string `json:"created_by"`
	CaseName     string `json:"casename,omitempty"`
	ComputerName string `json:"computername,omitempty"`
	Audit        string `json:"audit,omitempty"`
	DocId        string `json:"doc_id,,omitempty"`
	TagCategory  string `json:"tag_category"`
	TagName      string `json:"tag_name"`
}

func loadTagVars(tag *Tag, vars map[string]string) {
	if len(vars) == 0 {
		return
	}
	// Request URI parameters takes priority
	if vars["casename"] != "" {
		tag.CaseName = vars["casename"]
	}

	if vars["endpoint"] != "" {
		tag.ComputerName = vars["endpoint"]
	}
	if vars["audit"] != "" {
		tag.Audit = vars["audit"]
	}
	if vars["doc_id"] != "" {
		tag.DocId = vars["doc_id"]
	}
	if vars["tag_category"] != "" {
		tag.TagCategory = vars["tag_category"]
	}
	if vars["tag_name"] != "" {
		tag.TagName = vars["tag_name"]
	}
}

func AddTag(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// For UnitTest comment auth
	isauth, message := auth.IsAuthenticatedSession(w, r)
	if !isauth {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

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

	if tag.DocId == "" {
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
		Type(EsTagType).
		BodyJson(string(jsonTag)).
		Do(context.Background())

	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	api.HttpResponseReturn(w, r, "success", "Tag added successfully", res.Id)
}

// GetTag returns tag matching search condition
// api_uri: GET /api/v1/tag/show/{case}
// api_uri: GET /api/v1/tag/show/{case}/{endpoint}
// api_uri: GET /api/v1/tag/show/{case}/{endpoint}/{audit}
// api_uri: GET /api/v1/tag/show/{case}/{endpoint}/{audit}/{doc_id}
func GetTagData(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// For UnitTest comment auth
	isauth, message := auth.IsAuthenticatedSession(w, r)
	if !isauth {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

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
		queries = append(queries, elastic.NewWildcardQuery("casename.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("casename.keyword", tag.CaseName))
	}

	if tag.ComputerName == "" {
		queries = append(queries, elastic.NewWildcardQuery("computername.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("computername.keyword", tag.ComputerName))
	}

	if tag.Audit == "" {
		queries = append(queries, elastic.NewWildcardQuery("audit.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("audit.keyword", tag.Audit))
	}

	if tag.DocId == "" {
		queries = append(queries, elastic.NewWildcardQuery("doc_id.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("doc_id.keyword", tag.DocId))
	}

	if tag.CreatedBy != "" {
		queries = append(queries, elastic.NewTermQuery("created_by.keyword", tag.CreatedBy))
	}

	if tag.TagName != "" {
		queries = append(queries, elastic.NewWildcardQuery("tag_name.keyword", fmt.Sprintf("*%s*", tag.TagName)))
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
		Type(EsTagType).
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
