package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func getTimelineConfig() TimelineConfig {
	// Basic TimelineConfig
	config := TimelineConfig{
		SearchLimit: 10,
		SortDesc:    false,
		IgnoreGood:  false,
	}

	return config
}

func getGlobalSearchConfig() GlobalSearchConfig {
	config := GlobalSearchConfig{
		CaseName:   "*",
		Type:       "*",
		SearchSize: 10,
	}
	return config
}

func TestTimelineSearch(t *testing.T) {
	var config TimelineConfig

	// Base TestCases
	// TestCase-01: Timeline everything without any filter.
	// This timeline search will search everything in index
	fmt.Println("TestCase-01: Timeline index")
	config = getTimelineConfig()
	timelineSearch(t, config)

	// TestCase-02: Timeline everything in a given Case
	fmt.Println("TestCase-02: Timeline a Case")
	config = getTimelineConfig()
	config.CaseName = "CASE-BRAVO"
	timelineSearch(t, config)

	// TestCase-03: Timeline AuditType
	fmt.Println("TestCase-03: Timeline AuditType")
	config = getTimelineConfig()
	config.Type = "w32tasks"
	timelineSearch(t, config)

	// TestCase-04: Timeline single endpoint
	// If there are mulple cases it will timeline all the cases
	fmt.Println("TestCase-04: Timeline single endpoint")
	config = getTimelineConfig()
	config.Endpoint = "COMPUTER-01"
	timelineSearch(t, config)

	// TestCase-05: Timeline multiple endpoints
	fmt.Println("TestCase-05: Timeline multiple endpoints")
	computers := []string{"COMPUTER-01", "COMPUTER-02"}
	config = getTimelineConfig()
	config.ComputerList = computers
	timelineSearch(t, config)

	// TestCase-06: Timeline StartTime and EndTime
	fmt.Println("TestCase-06: Timeline by StartTime and EndTime")
	config = getTimelineConfig()
	config.StartTime = "2015-01-01T00:00:00Z"
	config.EndTime = "2015-02-01T00:00:00Z"
	timelineSearch(t, config)

	// TestCase-07: Timeline by StartTime and TimeDelta
	// StartTime is taken as EventTime
	fmt.Println("TestCase-07: Timeline by StartTime and TimeDelta")
	config = getTimelineConfig()
	config.StartTime = "2015-01-04:T09:00:00Z"
	config.TimeDelta = 15
	timelineSearch(t, config)

	// Compound TestCases
	// TestCase-08: Timeline by CaseName and Endpoint
	fmt.Println("TestCase-08: Timeline by Case and Endpoint")
	config = getTimelineConfig()
	config.CaseName = "SICHUANPEPPER"
	config.Endpoint = "COMPUTER-01"
	timelineSearch(t, config)

	// TestCase-09: Timeline by Case and multiple endpoints
	fmt.Println("TestCase-09: Timeline by Case and multiple endpoints")
	config = getTimelineConfig()
	config.CaseName = "SICHUANPEPPER"
	config.ComputerList = []string{"COMPUTER-01", "COMPUTER-02"}
	timelineSearch(t, config)
}

func TestGlobalSearch(t *testing.T) {
	var config GlobalSearchConfig

	// TestCase-01: Basic search
	fmt.Println("TestCase-01: Basic global search")
	config = getGlobalSearchConfig()
	config.SearchTerm = "audio"
	globalsearch(t, config)

	// TestCase-02: Search by casename
	fmt.Println("UserCase-02: Search term by case")
	config = getGlobalSearchConfig()
	config.CaseName = "CASE-BRAVO"
	config.SearchTerm = "audio"
	globalsearch(t, config)

	// TestCase-03: Search by AuditType
	fmt.Println("TestCase-03: Search term by audittype")
	config = getGlobalSearchConfig()
	config.Type = "w32scripting-persistence"
	config.SearchTerm = "google"
	globalsearch(t, config)

	// TestCase-04: Search by Case and AuditType
	fmt.Println("TestCase-04: Search term by Case and AuditType")
	config = getGlobalSearchConfig()
	config.CaseName = "CASE-BRAVO"
	config.Type = "w32scripting-persistence"
	config.SearchTerm = "google"
	globalsearch(t, config)

	// TestCase-05: Search by Endpoint
	fmt.Println("TestCase-05: Search term by endpoint")
	config = getGlobalSearchConfig()
	config.Endpoint = "COMPUTER-03"
	config.SearchTerm = "google"
	globalsearch(t, config)

	// TestCase-06: Search by Case and Endpoint
	fmt.Println("TestCase-06: Search term by Case and Endpoint")
	config = getGlobalSearchConfig()
	config.CaseName = "CASE-BRAVO"
	config.Endpoint = "COMPUTER-03"
	config.SearchTerm = "citrix"
	globalsearch(t, config)

	// TestCase-07: Search by AuditType and Endpoint
	fmt.Println("TestCase-07: Search term by AuditType and Endpoint")
	config = getGlobalSearchConfig()
	config.Type = "w32scripting-persistence"
	config.Endpoint = "COMPUTER-03"
	config.SearchTerm = "microsoft"
	globalsearch(t, config)

	// TestCase-08: Search by Case, Audit and Endpoint
	fmt.Println("TestCase-08: Search term by Case, AuditType and Endpoint")
	config = getGlobalSearchConfig()
	config.CaseName = "SICHUANPEPPER"
	config.Type = "w32scripting-persistence"
	config.Endpoint = "COMPUTER-03"
	config.SearchTerm = "oracle"
	globalsearch(t, config)
}

func timelineSearch(t *testing.T, config_data TimelineConfig) {
	jcdata, _ := json.Marshal(&config_data)
	req, err := http.NewRequest("POST", "/api/v1/search/timeline", bytes.NewReader(jcdata))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTimelineSearch)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	fmt.Println(rr.Body.String())
}

func globalsearch(t *testing.T, config_data GlobalSearchConfig) {
	jcdata, _ := json.Marshal(&config_data)
	req, err := http.NewRequest("POST", "/api/v1/search/globalsearch", bytes.NewReader(jcdata))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetGlobalSearch)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: %v want %v", status, http.StatusOK)
	}

	fmt.Println(rr.Body.String())
}

func TestTimelineSearchBug(t *testing.T) {
	config := TimelineConfig{
		CaseName:    "CASE-ALFA-02",
		Endpoint:    "D0023245BE4BE",
		Type:        "w32tasks",
		StartTime:   "1970-01-01T01:01:01Z",
		EndTime:     "2018-01-01T01:01:10Z",
		SearchLimit: 1000,
		SortDesc:    false,
		IgnoreGood:  false,
	}
	data, _ := json.Marshal(config)

	uri := fmt.Sprintf("/api/v1/search/timeline")
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/search/timeline", GetTimelineSearch)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}
