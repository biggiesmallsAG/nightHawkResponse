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

func TestAddComment(t *testing.T) {
	comment := Comment{
		CreatedBy:    "roshan",
		CaseName:     "CASE-ALFA-04",
		ComputerName: "COMPUTER-01",
		Audit:        "w32tasks",
		DocID:        "aaaaaaaaaabbbbbbb",
		Comment:      "This is comment about task",
	}

	cBody := make(map[string]string)
	cBody["comment"] = comment.Comment
	cBody["created_by"] = comment.CreatedBy
	data, _ := json.Marshal(cBody)

	uri := fmt.Sprintf("/api/v1/comment/add/%s/%s/%s/%s", comment.CaseName, comment.ComputerName, comment.Audit, comment.DocID)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/comment/add/{casename}/{endpoint}/{audit}/{doc_id}", AddComment)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}

func TestGetCommentByCase(t *testing.T) {
	search := make(map[string]string)

	fmt.Println("Test GetComment by valid casename")
	search["casename"] = "CASE-ALFA-01"
	getComment(t, search, nil)

	fmt.Println("Test GetComment by invalid casename")
	search["casename"] = "CASE-ALFA-02"
	getComment(t, search, nil)

}

func TestGetCommentByComputer(t *testing.T) {

	search := make(map[string]string)

	fmt.Println("Test GetComment by valid casename and computername")
	search["casename"] = "CASE-ALFA-01"
	search["computername"] = "COMPUTER-01"
	getComment(t, search, nil)

	fmt.Println("Test GetComment by casename and invalid computername")
	search["casename"] = "CASE-ALFA-01"
	search["computername"] = "COMPUTER-02"
	getComment(t, search, nil)
}

func TestGetCommentByAudit(t *testing.T) {

	search := make(map[string]string)

	fmt.Println("Test GetComment by valid audit")
	search["casename"] = "CASE-ALFA-01"
	search["computername"] = "COMPUTER-02"
	search["audit"] = "w32services"
	getComment(t, search, nil)

	fmt.Println("Test GetComment by audit")
	search["casename"] = "CASE-ALFA-01"
	search["computername"] = "COMPUTER-02"
	getComment(t, search, nil)
}

func TestGetCommentAll(t *testing.T) {
	search := make(map[string]string)

	fmt.Println("Get all comments")
	getComment(t, search, nil)

	fmt.Println("Get all comments for a case")
	search["casename"] = "CASE-ALFA-01"
	getComment(t, search, nil)

	fmt.Println("Get all comments for a case/endpoint")
	search["computername"] = "COMPUTER-01"
	getComment(t, search, nil)

	fmt.Println("Get all comments for a case/endpoint/audit")
	search["audit"] = "w32services"
	getComment(t, search, nil)

	fmt.Println("Get all comments for a case/endpoint/audit/doc_id")
	search["doc_id"] = "Abcdefgh1234567890"
	getComment(t, search, nil)
}

func TestGetCommentByPost(t *testing.T) {
	search := Comment{
		//CaseName:  "CASE-ALFA-01",
		CreatedBy: "",
		Comment:   "first",
	}

	data, _ := json.Marshal(search)
	getComment(t, nil, data)
}

func getComment(t *testing.T, search map[string]string, data []byte) {
	// Dynamically build uri as per search parameter
	uri := "/api/v1/comment/show"

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
	r.HandleFunc("/api/v1/comment/show", GetComment)
	r.HandleFunc("/api/v1/comment/show/{case}", GetComment)
	r.HandleFunc("/api/v1/comment/show/{case}/{endpoint}", GetComment)
	r.HandleFunc("/api/v1/comment/show/{case}/{endpoint}/{audit}", GetComment)
	r.HandleFunc("/api/v1/comment/show/{case}/{endpoint}/{audit}/{doc_id}", GetComment)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}

func TestGetCommentCaseAndComputer(t *testing.T) {
	uri := fmt.Sprintf("/api/v1/comment/case-and-computer")
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/comment/case-and-computer", GetCommentCaseAndComputer)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}
