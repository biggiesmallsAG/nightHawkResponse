package audit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetCaseList(t *testing.T) {
	uri := fmt.Sprintf("/api/v1/list/cases")
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/list/cases", GetCaseList)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("test result: ", rr.Body.String())
}
