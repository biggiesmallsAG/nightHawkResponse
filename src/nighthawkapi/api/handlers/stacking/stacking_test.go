package stacking

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestStacking(t *testing.T) {
    //serviceStacking(t)
    //taskStacking(t)
    prefetchStacking(t)
    //registryPersistenceStacking(t)
    //servicePersistenceStacking(t)
    //serviceDllPersistenceStacking(t)
    //linkPersistenceStacking(t)
    //allPersistenceStacking(t)
    //localListenPortStacking(t)
    //runKeyStacking(t)
    //dnsStacking(t)
    //urlDomainStacking(t)
}


func serviceStacking(t *testing.T) {
    data := StackingConfig{
        CaseName: "*",
        Type: "",
        SearchLimit: 5,
        SortDesc: true,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata,_ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/service", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackServices)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}




func taskStacking(t *testing.T) {
        data := StackingConfig{
        CaseName: "*",
        Type: "",
        SearchLimit: 5,
        SortDesc: false,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata, _ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/tasks", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackTasks)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}

func prefetchStacking(t *testing.T) {
        data := StackingConfig{
        CaseName: "*",
        Type: "",
        SearchLimit: 5,
        SortDesc: true,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata, _ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/prefetch", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackPrefetch)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
    fmt.Println(rr.Body.String())
}




func registryPersistenceStacking(t *testing.T) {
    persistenceStacking(t, "Registry")
}

func servicePersistenceStacking(t *testing.T) {
    persistenceStacking(t, "Service")
}

func serviceDllPersistenceStacking(t *testing.T) {
    persistenceStacking(t, "ServiceDll")
}

func linkPersistenceStacking(t *testing.T) {
    persistenceStacking(t, "Link")
}

func allPersistenceStacking(t *testing.T) {
    persistenceStacking(t, "*")
}

func persistenceStacking(t *testing.T, ptype string) {
    data := StackingConfig{
        CaseName: "*",
        Type: ptype,
        SearchLimit: 5,
        SortDesc: true,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata,_ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/persistence", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackPersistence)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}


func localListenPortStacking(t *testing.T) {
    data := StackingConfig{
        CaseName: "*",
        Type: "",
        SearchLimit: 15,
        SortDesc: true,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata,_ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/locallistport", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackLocalListenPort)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}


func runKeyStacking(t *testing.T) {
    data := StackingConfig{
        CaseName: "*",
        Type: "",
        SearchLimit: 15,
        SortDesc: true,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata,_ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/runkey", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackRunKey)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}


func dnsStacking(t *testing.T) {
    data := StackingConfig{
        CaseName: "*",
        Type: "",
        SearchLimit: 15,
        SortDesc: false,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata,_ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/dns/a", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackDnsARequest)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}

func urlDomainStacking(t *testing.T) {
    data := StackingConfig{
        CaseName: "*",
        Type: "",
        SearchLimit: 15,
        SortDesc: false,
        IgnoreGood: true,
        SubAggSize: 10,
    }

    jdata,_ := json.Marshal(&data)

    req, err := http.NewRequest("POST", "/api/v1/stacking/url/domain", bytes.NewReader(jdata))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(StackUrlDomain)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}