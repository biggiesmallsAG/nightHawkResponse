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

const EsCommentType = "comments"
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
	Timestamp    string `json:"timestamp"`
	CreatedBy    string `json:"created_by"`
	CaseName     string `json:"casename,omitempty"`
	ComputerName string `json:"computername,omitempty"`
	Audit        string `json:"audit,omitempty"`
	DocId        string `json:"doc_id,,omitempty"`
	Comment      string `json:"comment"`
}

func loadCommentVars(comment *Comment, vars map[string]string) {
	if len(vars) == 0 {
		return
	}
	// Request URI parameters takes priority
	if vars["casename"] != "" {
		comment.CaseName = vars["casename"]
	}

	if vars["endpoint"] != "" {
		comment.ComputerName = vars["endpoint"]
	}
	if vars["audit"] != "" {
		comment.Audit = vars["audit"]
	}
	if vars["doc_id"] != "" {
		comment.DocId = vars["doc_id"]
	}
}

// AddComment create new comment about an artifact
// api_uri: /api/v1/comment/add/{casename}/{endpoint}/{audit}/{doc_id}
// post_data: {"casedate":"2018-01-01T01:01:01Z","created_by": "user1", "comment": "this is comment"}
func AddComment(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// For UnitTest comment auth
	isauth, message := auth.IsAuthenticatedSession(w, r)
	if !isauth {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

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

	if comment.DocId == "" {
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

	res, err := client.Index().Index(conf.ServerIndex()).Type(EsCommentType).BodyJson(string(jsonComment)).Do(context.Background())
	if err != nil || !res.Created {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}
	api.HttpResponseReturn(w, r, "success", "Comment added", res.Id)
}

// GetComment returns comment matching search condition
// api_uri: GET /api/v1/comment/show/{case}
// api_uri: GET /api/v1/comment/show/{case}/{endpoint}
// api_uri: GET /api/v1/comment/show/{case}/{endpoint}/{audit}
// api_uri: GET /api/v1/comment/show/{case}/{endpoint}/{audit}/{doc_id}
func GetComment(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// For UnitTest comment auth
	isauth, message := auth.IsAuthenticatedSession(w, r)
	if !isauth {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

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
		queries = append(queries, elastic.NewWildcardQuery("casename.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("casename.keyword", comment.CaseName))
	}

	if comment.ComputerName == "" {
		queries = append(queries, elastic.NewWildcardQuery("computername.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("computername.keyword", comment.ComputerName))
	}

	if comment.Audit == "" {
		queries = append(queries, elastic.NewWildcardQuery("audit.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("audit.keyword", comment.Audit))
	}

	if comment.DocId == "" {
		queries = append(queries, elastic.NewWildcardQuery("doc_id.keyword", "*"))
	} else {
		queries = append(queries, elastic.NewTermQuery("doc_id.keyword", comment.DocId))
	}

	if comment.CreatedBy != "" {
		queries = append(queries, elastic.NewTermQuery("created_by.keyword", comment.CreatedBy))
	}

	if comment.Comment != "" {
		queries = append(queries, elastic.NewWildcardQuery("comment.keyword", fmt.Sprintf("*%s*", comment.Comment)))
	}

	boolquery := elastic.NewBoolQuery().Must(queries...)
	sr, err := client.Search().
		Index(conf.ServerIndex()).
		Type(EsCommentType).
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
