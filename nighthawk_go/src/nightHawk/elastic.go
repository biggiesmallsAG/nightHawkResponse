/*
 * @package 	nightHawk
 * @file		elastic.go
 * @author		Daniel Eden
 * @version		0.0.2
 * @updated 	2016-06-02
 */

package nightHawk

import (
 	"fmt"
 	"time"
 	"bytes"
 	"net/http"
 	"encoding/json"
 	"io/ioutil"
)


const Layout = "2006-01-02T15:04:05Z"


type ESSource struct {
	DateCreated 	string `json:"date_created"`
}

type ParentCheck struct {
	Found 			bool `json:"found"`
	Source			ESSource `json:"_source"`
}

type CB struct {
	Type 			string `json:"type"`
	Reason 			string `json:"reason"`
}

type RootCause struct {
	Type 			string `json:"type"`
	Reason 			string `json:"reason"`
	CausedBy 		CB 		`json:"caused_by"`
}

type ESErrorCheck struct {
	Error 			RootCause `json:"error"`
}

type Shard struct {
	Total 			int `json:"total"`
	Successful 		int `json:"successful"`
	Failed 			int `json:"failed"`
}

type ESSuccess struct {
	Shards 			Shard 	`json:"_shards"`
	Created 		bool 	`json:"created"`
}



func ExportToElastic(computername string, auditgenerator string, data []byte) {
	/*	First: understand if the parent exists, if it does just append the child.
	 *	Second: build string concat of ComputerName for index parent
	 *	Third: check if case date is same, if so dont update.
	 *  Fourth: build child document and post to parent
	 */
	 var ElasticHost = fmt.Sprintf("http://%s:%d", ELASTICHOST, ELASTICPORT)
	 var HostnameIndex = fmt.Sprintf("%s/hostname", ELASTIC_INDEX)

	 var parent ParentCheck
	 //var essuccess ESSuccess
	 //var eserr ESErrorCheck

	 es_url := fmt.Sprintf("%s/%s/%s", ElasticHost, HostnameIndex, computername)

	 res, err := http.Get(es_url)
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
 		pres,err := http.Post(post_url, "application/json", post_data)
 		if err != nil {
 			panic(err.Error())
 		}

 		ConsoleMessage("INFO", "New parent created for "+computername, VERBOSE)
 		if pres.StatusCode != 200 {
 			ConsoleMessage("ERROR", "Error creating parent node "+computername, true)
 		}
 		
	 }

	 // Adding record to Elastic
	audit_data := string(data[:])
	post_data := bytes.NewBufferString(audit_data)

	post_url := fmt.Sprintf("%s/investigations/_bulk", ElasticHost)
	pres, err := http.Post(post_url, "application/json", post_data)
	if err != nil {
		panic(err.Error())
	}

	prbody,_ := ioutil.ReadAll(pres.Body)
	if pres.StatusCode != 200 && pres.StatusCode != 201 {
		fmt.Println(string(prbody))
	}		
}