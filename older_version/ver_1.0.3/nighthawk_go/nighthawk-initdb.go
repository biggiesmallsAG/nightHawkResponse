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
    "bufio"
    "strings"
    "bytes"
    
    "time"
    
    "nightHawk"
)


func main() {
    esmap,_ := ioutil.ReadFile("/opt/nighthawk/lib/elastic/ElasticMapping.json")

    nightHawk.LoadConfigFile("/opt/nighthawk/etc/nighthawk.json")

    var esIndexUrl string = ""
    esIndexUrl = GetElasticUrl(nightHawk.ELASTIC_INDEX)

    resp, err := nightHawk.HttpOperation("GET", esIndexUrl, nightHawk.ELASTIC_AUTHCODE, nightHawk.ELASTIC_SSL, nil)
    if err != nil {
        nightHawk.ConsoleMessage("ERROR", "Error occured connecting to Elasticsearch Instance ["+esIndexUrl+"]", true)
        os.Exit(1)
        //panic(err.Error())
    }
    
    if resp.StatusCode == 200 || resp.StatusCode == 201 {
        reader := bufio.NewReader(os.Stdin)

        fmt.Printf("Index Already exists. Do you want to recreate [y/n]: ")
        user_input,_ := reader.ReadString('\n')

        if strings.ToLower(user_input[:1]) == "y" {
            esIndexUrl = GetElasticUrl("investigation1")
            // Delete current index
            resp, err := nightHawk.HttpOperation("DELETE", esIndexUrl, nightHawk.ELASTIC_AUTHCODE, nightHawk.ELASTIC_SSL, nil)
            if err != nil {
                nightHawk.ConsoleMessage("ERROR", "Error deleting default index investigation1 ["+esIndexUrl+"]", true)
                //panic(err.Error())
                os.Exit(1)
            }

            if resp.StatusCode != 200 && resp.StatusCode != 201 {
                body,_ := ioutil.ReadAll(resp.Body)
                fmt.Println(string(body))

                os.Exit(1)
            }
            os.Remove("/opt/nighthawk/var/run/elasticsearch/elastic.init")
            CreateElasticIndex(esIndexUrl, esmap)
        }       
    } else {
        esIndexUrl = GetElasticUrl("investigation1")
        CreateElasticIndex(esIndexUrl, esmap)
    }
}

func GetElasticUrl(IndexName string) string {
    var esIndexUrl string = ""

    if nightHawk.ELASTIC_SSL {
        esIndexUrl = fmt.Sprintf("https://%s:%d/%s", nightHawk.ELASTICHOST, nightHawk.ELASTICPORT, IndexName)
    } else {
        esIndexUrl = fmt.Sprintf("http://%s:%d/%s", nightHawk.ELASTICHOST, nightHawk.ELASTICPORT, IndexName)
    }

    return esIndexUrl
}

func CreateElasticIndex(esIndexUrl string, data []byte) {
    postData := bytes.NewBuffer(data)

    resp, err := nightHawk.HttpOperation("POST", esIndexUrl, nightHawk.ELASTIC_AUTHCODE, nightHawk.ELASTIC_SSL, postData)
    if err != nil {
        nightHawk.ConsoleMessage("ERROR", "Error occured posting data to Elasticsearch instance ["+esIndexUrl+"]", true)
        //panic(err.Error())
        os.Exit(1)
    }
    
    if resp.StatusCode != 200 && resp.StatusCode != 201 {
        body,_ := ioutil.ReadAll(resp.Body)
        fmt.Println(string(body))
        os.Exit(2)
    }
    
    err = ioutil.WriteFile("/opt/nighthawk/var/run/elasticsearch/elastic.init", []byte(time.Now().UTC().String()), 0644)
    if err != nil {
        nightHawk.ConsoleMessage("ERROR", "Failed to write elastic init file", true)
        //panic(err.Error())
        os.Exit(2)
    }
    os.Chown("/opt/nighthawk/var/run/elasticsearch/elastic.init", 3728, 3728)
}
