package stacking

import (
    "fmt"
    "testing"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "net/http/httptest"
)

func getDiffConfig() DiffConfig {
    config := DiffConfig{Endpoint:""}
    return config
}

func printConfig(config DiffConfig) {
    json_config,_ := json.Marshal(config)
    fmt.Println(string(json_config))
}


func TestEndpointDiff(t *testing.T) {
    var config DiffConfig

    fmt.Println("TestCase-01: Diffing single endpoint")
    config = getDiffConfig()
    config.Endpoint="WHISKEY01"
    printConfig(config)
    diffEndpoint(t, config)

    fmt.Println("TestCase-02: Diffing Endpoint by name: /api/v1/diff/{endpoint}")
    getEndpointDiffApi("/api/v1/diff/WHISKEY01")
}

func diffEndpoint(t *testing.T, config DiffConfig) {
    jconfig,_ := json.Marshal(&config)
    req, err := http.NewRequest("POST", "/api/v1/diff", bytes.NewReader(jconfig))
    if err != nil {t.Fatal(err)}

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(TimelineEndpointDiff)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: %v want %v", status, http.StatusOK)
    }

    fmt.Println(rr.Body.String())
}

func getEndpointDiffApi(uri string) {
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
