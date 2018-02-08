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

const esCommentType = "comments"
const Layout = "2006-01-02T15:04:05Z"

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

type Comment struct {
	Timestamp    string
	CreatedBy    string
	CaseName     string `json:"CaseName,omitempty"`
	ComputerName string `json:"ComputerName,omitempty"`
	Audit        string `json:"Audit,omitempty"`
	DocID        string `json:"DocID,,omitempty"`
	Comment      string
}

func loadCommentVars(comment *Comment, vars map[string]string) {
	if len(vars) == 0 {
		return
	}
	// Request URI parameters takes priority
	if vars["casename"] != "" {
		comment.CaseName = vars["casename"]
	}

	if vars["case_name"] != "" {
		comment.CaseName = vars["case_name"]
	}

	if vars["case"] != "" {
		comment.CaseName = vars["case"]
	}

	if vars["endpoint"] != "" {
		comment.ComputerName = vars["endpoint"]
	}
	if vars["audit"] != "" {
		comment.Audit = vars["audit"]
	}
	if vars["doc_id"] != "" {
		comment.DocID = vars["doc_id"]
	}
}

// AddComment create new comment about an artifact
// api_uri: /api/v1/comment/add/{casename}/{endpoint}/{audit}/{doc_id}
// post_data: {"casedate":"2018-01-01T01:01:01Z","created_by": "user1", "comment": "this is comment"}
func AddComment(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	var comment Comment
	json.Unmarshal(body, &comment)
	loadCommentVars(&comment, mux.Vars(r))

	// Comment Data validation
	if comment.CreatedBy == "" {
		api.HttpResponseReturn(w, r, "failed", "Comment creator required", nil)
		return
	}

	if comment.CaseName == "" {
		api.HttpResponseReturn(w, r, "failed", "Casename is required", nil)
		return
	}

	if comment.ComputerName == "" {
		api.HttpResponseReturn(w, r, "failed", "Computer name is required", nil)
		return
	}

	if comment.DocID == "" {
		api.HttpResponseReturn(w, r, "failed", "Artifact document ID  required", nil)
		return
	}

	if comment.Comment == "" {
		api.HttpResponseReturn(w, r, "failed", "Comment data required", nil)
		return
	}

	// Sett Comment creation time
	comment.Timestamp = time.Now().UTC().Format(Layout)

	jsonComment, _ := json.Marshal(comment)

	res, err := client.Index().Index(conf.ServerIndex()).Type(esCommentType).BodyJson(string(jsonComment)).Do(context.Background())
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}
	api.HttpResponseReturn(w, r, "success", "Comment added", res.Id)
}

// GetComment returns comment matching search condition
// api_uri: GET /api/v1/comment/show/{case_name}
// api_uri: GET /api/v1/comment/show/{case_name}/{endpoint}
// api_uri: GET /api/v1/comment/show/{case_name}/{endpoint}/{audit}
// api_uri: GET /api/v1/comment/show/{case_name}/{endpoint}/{audit}/{doc_id}
func GetComment(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var comment Comment

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
			return
		}

		json.Unmarshal(body, &comment)
	}

	loadCommentVars(&comment, mux.Vars(r))

	var queries []elastic.Query
	if comment.CaseName == "" {
		queries = append(queries, elastic.NewWildcardQuery("CaseName.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("CaseName.keyword", comment.CaseName))
	}

	if comment.ComputerName == "" {
		queries = append(queries, elastic.NewWildcardQuery("ComputerName.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("ComputerName.keyword", comment.ComputerName))
	}

	if comment.Audit == "" {
		queries = append(queries, elastic.NewWildcardQuery("Audit.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("Audit.keyword", comment.Audit))
	}

	if comment.DocID == "" {
		queries = append(queries, elastic.NewWildcardQuery("DocID.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("DocID.keyword", comment.DocID))
	}

	if comment.CreatedBy != "" {
		queries = append(queries, elastic.NewTermQuery("CreatedBy.keyword", comment.CreatedBy))
	}

	if comment.Comment != "" {
		queries = append(queries, elastic.NewWildcardQuery("Comment.keyword", fmt.Sprintf("*%s*", comment.Comment)))
	}

	boolquery := elastic.NewBoolQuery().Must(queries...)
	sr, err := client.Search().
		Index(conf.ServerIndex()).
		Type(esCommentType).
		Query(boolquery).
		Do(context.Background())

	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	if sr.TotalHits() < 1 {
		api.HttpResponseReturn(w, r, "failed", "Comment not found", nil)
		return
	}

	// Populating array of comments
	var comments []Comment
	for _, hit := range sr.Hits.Hits {
		var comment Comment
		json.Unmarshal(*hit.Source, &comment)
		comments = append(comments, comment)
	}

	api.HttpResponseReturn(w, r, "success", "Comment search completed", comments)
}

// GetCommentCaseAndComputer provides list of CaseName and
// list of ComputerName
// GET /api/v1/comment/cases-and-computers
func GetCommentCaseAndComputer(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var casenameList []string
	var computernameList []string

	data := make(map[string]interface{})

	aggCasename := elastic.NewTermsAggregation().
		Field("CaseName.keyword")
	sr, err := client.Search().
		Index(conf.ServerIndex()).Type(esCommentType).
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
	sr, err = client.Search().Index(conf.ServerIndex()).Type(esCommentType).Aggregation("agg_computer", aggComputer).Do(context.Background())
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

	api.HttpResponseReturn(w, r, "success", "Comment CaseNames and ComputerNames", data)
}
