/*
 * file:        nighthawk-initdb.go
 * author:      roshan maskey <roshanmaskey@gmail.com>
 * version:     0.0.2
 * updated:     2016-07-10
 *
 * description: 
 * This file initializes Elasticsearch database.
 */

package main

import (
    "fmt"
    "io/ioutil"
    "os"
    //"bufio"
    //"strings"
    "bytes"
    "time"
    "flag"

    "nighthawk"
    nhconfig "nighthawk/config"
    nhelastic "nighthawk/elastic"
)


type RuntimeOptions struct {
    IndexName           string
    CreateNewIndex      bool
    ConfigFile          string
    ElasticMapFile      string
    Verbose             bool
}


func main() {

    var runopt RuntimeOptions

    flag.StringVar(&runopt.IndexName,"index-name","investigation1", "Name of index to be created")
    flag.StringVar(&runopt.ConfigFile, "c", "/opt/nighthawk/etc/nighthawk.json", "nighthawk response configuration file")
    flag.StringVar(&runopt.ElasticMapFile, "m", "/opt/nighthawk/lib/elastic/ElasticMapping.json", "Elasticsearch mapping file")
    flag.BoolVar(&runopt.CreateNewIndex, "create-new",false, "Overwrite the index. Delete and create new index")
    flag.BoolVar(&runopt.Verbose, "v", false, "Show verbose messages")

    flag.Parse()


    esmap,_ := ioutil.ReadFile(runopt.ElasticMapFile)
    nhconfig.LoadConfigFile(runopt.ConfigFile)

    var esIndexUrl string = "" 
    esIndexUrl = GetElasticUrl(runopt.IndexName)

    resp, err := nhelastic.HttpOperation("GET", esIndexUrl, nhconfig.ELASTIC_AUTHCODE, nhconfig.ELASTIC_SSL, nil)
    if err != nil {
        nighthawk.ConsoleMessage("ERROR", "Error occured connecting to Elasticsearch Instance ["+esIndexUrl+"]", true)
        os.Exit(1)
    }
    
    if resp.StatusCode == 200 || resp.StatusCode == 201 {
        if runopt.CreateNewIndex {
            msg := fmt.Sprintf("Deleting existing index %s", runopt.IndexName)
            nighthawk.ConsoleMessage("INFO", msg, true)
            esIndexUrl = GetElasticUrl(runopt.IndexName)
            resp, err := nhelastic.HttpOperation("DELETE", esIndexUrl, nhconfig.ELASTIC_AUTHCODE, nhconfig.ELASTIC_SSL, nil)
            if err != nil {
                msg := fmt.Sprintf("Error deleting exising index %s: %s", runopt.IndexName, esIndexUrl)
                nighthawk.ConsoleMessage("ERROR", msg, true)
                os.Exit(1)
            }

            if resp.StatusCode != 200 && resp.StatusCode != 201 {
                msg := fmt.Sprintf("Failed to create index %s with status code %d", esIndexUrl, resp.StatusCode)
                nighthawk.ConsoleMessage("ERROR", msg, true)

                body,_ := ioutil.ReadAll(resp.Body)
                fmt.Println(string(body))

                os.Exit(1)
            }
            nighthawk.ConsoleMessage("DEBUG", "Deleting elasticsearch.init", runopt.Verbose)
            os.Remove("/opt/nighthawk/var/run/elasticsearch.init")

            msg = fmt.Sprintf("Re-creating index %s", runopt.IndexName)
            nighthawk.ConsoleMessage("INFO", msg, true)
            CreateElasticIndex(esIndexUrl, esmap, runopt.Verbose)
        } else {
            msg := fmt.Sprintf("Index %s already created. To delete and re-create use -create-new option", runopt.IndexName)
            nighthawk.ConsoleMessage("INFO", msg, true)
        }   
    } else {
        msg := fmt.Sprintf("Creating new index %s", runopt.IndexName)
        nighthawk.ConsoleMessage("INFO", msg, true)
        esIndexUrl = GetElasticUrl("investigation1")
        CreateElasticIndex(esIndexUrl, esmap, runopt.Verbose)
    }
}


/*
 * This function returns the full uri as http://127.0.0.1:9200/investigation1 or https://127.0.0.1:9200/investigation1
 * The decision to pick http or https is based on nighthawk.json configuration file
 */

func GetElasticUrl(IndexName string) string {
    var esIndexUrl string = ""

    if nhconfig.ELASTIC_SSL {
        esIndexUrl = fmt.Sprintf("https://%s:%d/%s", nhconfig.ELASTICHOST, nhconfig.ELASTICPORT, IndexName)
    } else {
        esIndexUrl = fmt.Sprintf("http://%s:%d/%s", nhconfig.ELASTICHOST, nhconfig.ELASTICPORT, IndexName)
    }

    return esIndexUrl
}

/* 
 * This function creates a new elasticsearch index
 * elasticsearch.init : 
 *      File created to indicate Elasticsearch index has been initialized. This does not take consideration
 *      of multiple index.
 *      Full path for this file: /opt/nighthawk/var/run/elasticsearch.init 
 * 
 * The function also changes the owner of file "elasticsearch.init" to nighthawk:nighthawk i.e. uid: 3728, gid:3728
 */
func CreateElasticIndex(esIndexUrl string, data []byte, flgVerbose bool) {
    postData := bytes.NewBuffer(data)

    //resp, err := nhelastic.HttpOperation("POST", esIndexUrl, nhconfig.ELASTIC_AUTHCODE, nhconfig.ELASTIC_SSL, postData)
    resp, err := nhelastic.ElasticPUT(esIndexUrl, nhconfig.ELASTIC_AUTHCODE, nhconfig.ELASTIC_SSL, postData)
    if err != nil {
        nighthawk.ConsoleMessage("ERROR", "Error occured posting data to Elasticsearch instance ["+esIndexUrl+"]", true)
        os.Exit(1)
    }
    
    if resp.StatusCode != 200 && resp.StatusCode != 201 {
        msg := fmt.Sprintf("Failed to create index %s with status code %d", esIndexUrl, resp.StatusCode)
        nighthawk.ConsoleMessage("ERROR", msg, true)
        body,_ := ioutil.ReadAll(resp.Body)
        fmt.Println(string(body))
        os.Exit(2)
    }
    
    err = ioutil.WriteFile("/opt/nighthawk/var/run/elasticsearch.init", []byte(time.Now().UTC().String()), 0644)
    if err != nil {
        nighthawk.ConsoleMessage("ERROR", "Failed to write elastic init file", true)
        os.Exit(2)
    }

    nighthawk.ConsoleMessage("DEBUG", "Creating elasticsearch.init", flgVerbose)
    os.Chown("/opt/nighthawk/var/run/elasticsearch.init", 3728, 3728)
}
