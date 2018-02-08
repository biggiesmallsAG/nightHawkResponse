package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestAddTag(t *testing.T) {

	casename := "CASE-ALFA-03"
	computername := "COMPUTER-01"
	audit := "w32tasks"
	doc_id := "aaaaaaaaaabbbbbbb"

	cbody := make(map[string]string)
	cbody["created_by"] = "roshan"
	cbody["tag_category"] = "network"
	cbody["tag_name"] = "malicious"
	data, _ := json.Marshal(cbody)

	uri := fmt.Sprintf("/api/v1/tag/add/%s/%s/%s/%s", casename, computername, audit, doc_id)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/tag/add/{casename}/{endpoint}/{audit}/{doc_id}", AddTag)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())

}

func TestGetTagAll(t *testing.T) {
	search := make(map[string]string)

	fmt.Println("Get all tags")
	getTag(t, search, nil)

	fmt.Println("Get all tags for a casename")
	search["casename"] = "CASE-ALFA-01"
	getTag(t, search, nil)

	fmt.Println("Get all tags for a casename/endpoint")
	search["computername"] = "COMPUTER-01"
	getTag(t, search, nil)

	fmt.Println("Get all tags for a casename/endpoint/audit")
	search["audit"] = "w32services"
	getTag(t, search, nil)

	fmt.Println("Get all tags for a casename/endpoint/audit/doc_id")
	search["doc_id"] = "Abcdefgh1234567890"
	getTag(t, search, nil)
}

// Test GetTag by using
// POST /api/v1/tag/show
func TestGetTagByPost(t *testing.T) {
	search := Tag{
		//CaseName:  "CASE-ALFA-01",
		CreatedBy: "",
		TagName:   "malicious",
	}

	data, _ := json.Marshal(search)
	getTag(t, nil, data)
}

func getTag(t *testing.T, search map[string]string, data []byte) {
	// Dynamically build uri as per search parameter
	uri := "/api/v1/tag/show"

	if search["casename"] != "" {
		uri = fmt.Sprintf("%s/%s", uri, search["casename"])
	}

	if search["computername"] != "" {
		uri = fmt.Sprintf("%s/%s", uri, search["computername"])
	}

	if search["audit"] != "" {
		uri = fmt.Sprintf("%s/%s", uri, search["audit"])
	}

	if search["doc_id"] != "" {
		uri = fmt.Sprintf("%s/%s", uri, search["doc_id"])
	}

	req, _ := http.NewRequest("GET", uri, nil)

	if data != nil {
		req, _ = http.NewRequest("POST", uri, bytes.NewReader(data))
	}

	fmt.Printf("%s %s\n", req.Method, req.URL)
	if data != nil {
		fmt.Println(string(data))
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/tag/show", GetTagData)
	r.HandleFunc("/api/v1/tag/show/{casename}", GetTagData)
	r.HandleFunc("/api/v1/tag/show/{casename}/{endpoint}", GetTagData)
	r.HandleFunc("/api/v1/tag/show/{casename}/{endpoint}/{audit}", GetTagData)
	r.HandleFunc("/api/v1/tag/show/{casename}/{endpoint}/{audit}/{doc_id}", GetTagData)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}

func TestGetTagCaseAndComputer(t *testing.T) {
	uri := fmt.Sprintf("/api/v1/tag/case-and-computer")
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/tag/case-and-computer", GetTagCaseAndComputer)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}
