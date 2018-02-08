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

func TestGetCasedateByEndpoint(t *testing.T) {
	uri := fmt.Sprintf("/api/v1/show/CASE-ALFA-02/COMPUTER-1")
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/show/{case}/{endpoint}", GetCasedateByEndpoint)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}

func TestGetDocById(t *testing.T) {
	//uri := fmt.Sprintf("/api/v1/show/doc/AWFAdtEXFQ83gjENqLRE/DTS01W006")
	uri := fmt.Sprintf("/api/v1/show/doc/AWFAdtEXFQ83gjENqLRE/CASE-ALFA-001")
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/show/doc/{doc_id}/{endpoint}", GetDocById)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}

func hfGetAuditDataByAuditGenerator(t *testing.T, uri string, data []byte) {
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/show/{case}/{endpoint}/{case_date}/{audittype}", GetAuditDataByAuditGenerator)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}

func TestGetAuditDataByAuditGenerator(t *testing.T) {
	var pd AuditFilter

	uri := fmt.Sprintf("/api/v1/show/CASE-ALFA-001/DTS01W006/2018-01-31T00:10:30.000Z/w32eventlogs?from=0&size=100&sort=Record.TlnTime&order=desc")

	pd.ColId = "Record.EID"
	pd.Filter.MatchType = "contains"
	pd.Filter.FilterType = "text"
	pd.Filter.Value = "60"

	data, _ := json.Marshal(pd)
	fmt.Println(string(data))
	hfGetAuditDataByAuditGenerator(t, uri, data)
}

func TestGetAuditDataGetTime(t *testing.T) {
	var pd AuditFilter

	uri := fmt.Sprintf("/api/v1/show/CASE-ALFA-001/DTS01W006/2018-01-31T00:10:30.000Z/w32eventlogs?from=0&size=100&sort=Record.TlnTime&order=desc")

	pd.ColId = "Record.GenTime"
	pd.Filter.MatchType = "lessThan"
	pd.Filter.FilterType = "date"
	pd.Filter.StartDate = "2016-09-13"
	//pd.Filter.EndDate = "2016-10-13"

	data, _ := json.Marshal(pd)
	fmt.Println(string(data))
	hfGetAuditDataByAuditGenerator(t, uri, data)
}

func TestGetAuditDataIndex(t *testing.T) {
	var pd AuditFilter

	uri := fmt.Sprintf("/api/v1/show/CASE-ALFA-001/DTS01W006/2018-01-31T00:10:30.000Z/w32eventlogs?from=0&size=100&sort=Record.TlnTime&order=desc")

	pd.ColId = "Record.Index"
	pd.Filter.MatchType = "equal"
	pd.Filter.FilterType = "string"
	pd.Filter.Value = "100"

	data, _ := json.Marshal(pd)
	fmt.Println(string(data))
	hfGetAuditDataByAuditGenerator(t, uri, data)
}

func TestGetAuditDataHttpGet(t *testing.T) {
	uri := fmt.Sprintf("/api/v1/show/CASE-ALFA-001/DTS01W006/2018-01-31T00:10:30.000Z/w32eventlogs?from=0&size=100&sort=Record.TlnTime&order=desc")

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/show/{case}/{endpoint}/{case_date}/{audittype}", GetAuditDataByAuditGenerator)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}
