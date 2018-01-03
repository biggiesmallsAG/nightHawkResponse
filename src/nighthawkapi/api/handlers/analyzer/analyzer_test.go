package analyzer

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)


func TestBlacklistEntry(t *testing.T) {
	data := AnalyzeItem {
        Title: "Blacklist Driver",
        Description: "This driver is used by xyz",
		AuditType: "w32services",
		Name: "1394ohci01",
		Path: "C:\\Windows\\System32\\drivers\\1394ohci.sys",
		ServiceDescriptiveName: "1394 OHCI Compliant Host Controller",
	}
	jdata,_ := json.Marshal(&data)

	req, err := http.NewRequest("POST", "/api/v1/analyze/blacklist", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(AddBlacklistInformation)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}


func TestWhitelistEntry(t *testing.T) {
	data := AnalyzeItem {
		AuditType: "w32services",
		Name: "1394ohci01",
		Path: "C:\\Windows\\System32\\drivers\\1394ohci.sys",
		ServiceDescriptiveName: "1394 OHCI Compliant Host Controller",
	}
	jdata,_ := json.Marshal(&data)

	req, err := http.NewRequest("POST", "/api/v1/analyze/whitelist", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(AddWhitelistInformation)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}



func TestStackCommonEntry(t *testing.T) {
	data := AnalyzeItem {
		AuditType: "w32services",
		Name: "1394ohci03",
		Path: "C:\\Windows\\System32\\drivers\\1394ohci.sys",
		ServiceDescriptiveName: "1394 OHCI Compliant Host Controller",
	}
	jdata,_ := json.Marshal(&data)

	req, err := http.NewRequest("POST", "/api/v1/analyze/stack", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(AddStackInformation)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}



func TestDeleteAnalyzerItemByID(t *testing.T) {
	
	req, err := http.NewRequest("GET", "/api/v1/analyze/delete/blacklist/AWCHsJXj3YG5pniWbNjZ", nil)

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(DeleteAnalyzerItemByID)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}


func TestDeleteAnalyzerItemByQuery(t *testing.T) {
	data := AnalyzeItem {
		AuditType: "w32services",
		Name: "1394ohci01",
		Path: "C:\\Windows\\System32\\drivers\\1394ohci.sys",
		ServiceDescriptiveName: "1394 OHCI Compliant Host Controller",
	}
	jdata,_ := json.Marshal(&data)

	req, err := http.NewRequest("POST", "/api/v1/analyze/delete/blacklist", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(AddStackInformation)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}