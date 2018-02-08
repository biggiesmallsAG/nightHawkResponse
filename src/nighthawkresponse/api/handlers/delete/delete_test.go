package delete

import (
    "fmt"
    "testing"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "net/http/httptest"
)


func getDeleteConfig() DeleteConfig {
    config := DeleteConfig{}

    return config
}


func printConfig(config DeleteConfig) {
    json_config,_ := json.Marshal(config)
    fmt.Println(string(json_config))
}


func TestDelete(t *testing.T) {
    var config DeleteConfig

    fmt.Println("TestCase-01: Delete single endpoint")
    config = getDeleteConfig()
    config.Endpoint = "COMPUTER-C0"
    printConfig(config)
    deleteEndpoint(t, config)

    fmt.Println("TestCase-02: Delete multiple endpoints")
    el := []string{"COMPUTER-C4", "COMPUTER-C5"}
    config = getDeleteConfig()
    config.EndpointList = el
    printConfig(config)
    deleteEndpoint(t, config)

    fmt.Println("TestCase-03: Delete single endpoint by case")
    config = getDeleteConfig()
    config.CaseName = "CASE-CHARLIE"
    config.Endpoint = "COMPUTER-C11"
    printConfig(config)
    deleteEndpoint(t, config)

    fmt.Println("TestCase-04: Delete multiple endpoints by case")
    config = getDeleteConfig()
    config.CaseName = "CASE-CHARLIE"
    config.EndpointList = []string{"COMPUTER-C15", "COMPUTER-C17"}
    printConfig(config)
    deleteEndpoint(t, config)

    fmt.Println("TestCase-05: Delete case by name")
    config = getDeleteConfig()
    config.CaseName = "CASE-DELTA"
    printConfig(config)
    deleteCase(t, "POST", "/api/v1/delete/case", config)


    fmt.Println("TestCase-06: Delete Endpoint: /api/v1/delete/endpoint/{endpoint_name}")
    getDeleteApi("/api/v1/delete/endpoint/COMPUTER-E17")

    fmt.Println("TestCase-07: Delete CaseEndpoint: /api/v1/delete/{case_name}/{endpoint_name}")
    getDeleteApi("/api/v1/delete/CASE-ECHO/COMPUTER-E11")

    fmt.Println("TestCase-10: Delete CaseByName: /api/v1/delete/case/{case_name}")
    getDeleteApi("/api/v1/delete/case/CASE-FOXTROT")
}

func deleteCase(t *testing.T, method string, uri string, config DeleteConfig) {
    jconfig,_ := json.Marshal(&config)
    req, err := http.NewRequest(method, uri, bytes.NewReader(jconfig))
    if err != nil {t.Fatal(err)}

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(DeleteCase)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}

func deleteEndpoint(t *testing.T, config DeleteConfig) {
    jconfig,_ := json.Marshal(&config)
    req, err := http.NewRequest("POST", "/api/v1/delete/endpoint", bytes.NewReader(jconfig))
    if err != nil {t.Fatal(err)}

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(DeleteEndpoint)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}

func getDeleteApi(uri string) {
    full_uri := fmt.Sprintf("http://localhost:8080%s", uri)
    fmt.Println("Calling API URL: ", full_uri)
    res, err := http.Get(full_uri)
    if err != nil {
        fmt.Println("Verify nhapi is running")
        return
    }

    body,_ := ioutil.ReadAll(res.Body)
    res.Body.Close()
    fmt.Println(string(body))
}
