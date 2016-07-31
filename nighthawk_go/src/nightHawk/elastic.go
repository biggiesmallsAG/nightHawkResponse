/*
 * @package     nightHawk
 * @file        elastic.go
 * @author      Daniel Eden, 0xredskull
 */

package nightHawk

import (
    "fmt"
    "time"
    "bytes"
    "net/http"
    "encoding/json"
    "io"
    "io/ioutil"
    "crypto/tls"
)


const Layout = "2006-01-02T15:04:05Z"


type ESSource struct {
    DateCreated     string `json:"date_created"`
}

type ParentCheck struct {
    Found           bool `json:"found"`
    Source          ESSource `json:"_source"`
}

type CB struct {
    Type            string `json:"type"`
    Reason          string `json:"reason"`
}

type RootCause struct {
    Type            string `json:"type"`
    Reason          string `json:"reason"`
    CausedBy        CB      `json:"caused_by"`
}

type ESErrorCheck struct {
    Error           RootCause `json:"error"`
}

type Shard struct {
    Total           int `json:"total"`
    Successful      int `json:"successful"`
    Failed          int `json:"failed"`
}

type ESSuccess struct {
    Shards          Shard   `json:"_shards"`
    Created         bool    `json:"created"`
}



func ExportToElastic(computername string, auditgenerator string, data []byte) {
    /*  First: understand if the parent exists, if it does just append the child.
     *  Second: build string concat of ComputerName for index parent
     *  Third: check if case date is same, if so dont update.
     *  Fourth: build child document and post to parent
     */

     var ElasticHost string 
     if ELASTIC_SSL {
         ElasticHost = fmt.Sprintf("https://%s:%d", ELASTICHOST, ELASTICPORT)
     } else {
         ElasticHost = fmt.Sprintf("http://%s:%d", ELASTICHOST, ELASTICPORT)
     }
    
     var HostnameIndex = fmt.Sprintf("%s/hostname", ELASTIC_INDEX)

     var parent ParentCheck

     es_url := fmt.Sprintf("%s/%s/%s", ElasticHost, HostnameIndex, computername)
     
     res, err := HttpOperation("GET", es_url, ELASTIC_AUTHCODE, ELASTIC_SSL, nil)
     if err != nil {
        panic(err.Error())
     }
        
     defer res.Body.Close()
     body,_ := ioutil.ReadAll(res.Body)
     json.Unmarshal(body, &parent)

     // Create new parent if it does not exist
     if !parent.Found {

        cur_time := fmt.Sprintf("%s", time.Now().UTC().Format(Layout))
        parent_data := fmt.Sprintf("{\"index\": {\"_id\": \"%s\"}}\n {\"date_created\":\"%s\"}\n", computername, cur_time)

        post_data := bytes.NewBufferString(parent_data)

        post_url := fmt.Sprintf("%s/%s/_bulk", ElasticHost, HostnameIndex)
        pres,err := HttpOperation("POST", post_url, ELASTIC_AUTHCODE, ELASTIC_SSL, post_data)

        if err != nil {
            panic(err.Error())
        }

        ConsoleMessage("INFO", "New parent created for "+computername, VERBOSE)
        RedisPublish("INFO", "New parent created for "+computername, REDIS_PUB)
        if pres.StatusCode != 200 {
            ConsoleMessage("ERROR", "Error creating parent node "+computername, true)
            RedisPublish("ERROR", "Error creating parent node "+computername, REDIS_PUB)
        }
        
     }

     // Adding record to Elastic
    audit_data := string(data[:])
    post_data := bytes.NewBufferString(audit_data)

    post_url := fmt.Sprintf("%s/investigations/_bulk", ElasticHost)
    pres, err := HttpOperation("POST", post_url, ELASTIC_AUTHCODE, ELASTIC_SSL, post_data)
    if err != nil {
        panic(err.Error())
    }

    prbody,_ := ioutil.ReadAll(pres.Body)
    if pres.StatusCode != 200 && pres.StatusCode != 201 {
        esErr := &ESErrorCheck{}
        json.Unmarshal(prbody, esErr)

        ConsoleMessage("ERROR", esErr.Error.Reason, VERBOSE)
        RedisPublish("ERROR", esErr.Error.Reason, REDIS_PUB)
    }       
}


func HttpOperation(method string, url string, authcode string, sslenabled bool, data io.Reader) (resp *http.Response, err error) {
    req, err := http.NewRequest(method,url, data)
    req.Header.Add("Authorization", "Basic " + authcode)
    
    if sslenabled {
        transcfg := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
        client := &http.Client{Transport: transcfg}
        return client.Do(req)
    } else {
        client := &http.Client{}
        return client.Do(req)
    }
} // __HttpOperation__
