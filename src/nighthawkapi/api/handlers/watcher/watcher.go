package watcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

	"github.com/gorilla/mux"

	elastic "gopkg.in/olivere/elastic.v5"
	yaml "gopkg.in/yaml.v2"
)

func GetWatcherResults(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method string
		wr     api.WatcherResults
		_wr    []api.WatcherResults
		query  elastic.Query
		ret    string
		conf   *config.ConfigVars
		err    error
		client *elastic.Client
		alerts *elastic.SearchResult
	)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if conf.Elastic.Elastic_ssl {
		method = api.HTTPS
	} else {
		method = api.HTTP
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s%s:%d", method, conf.Elastic.Elastic_server, conf.Elastic.Elastic_port)))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	query = elastic.NewTermQuery("alert_sent", true)

	alerts, err = client.Search().
		Index("elastalert_status").
		Query(query).
		Size(1000).
		Do(context.Background())

	if err != nil {
		api.LogError(api.DEBUG, err)
		ret := api.HttpFailureMessage(err.Error())
		fmt.Fprint(w, string(ret))
		return
	}
	
	for _, hit := range alerts.Hits.Hits {
		json.Unmarshal(*hit.Source, &wr)
		_wr = append(_wr, wr)
	}

	ret = api.HttpSuccessMessage("200", &_wr, alerts.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /watcher/results HTTP 200, returned watcher results")
	fmt.Fprintln(w, string(ret))

}

func GetWatcherResultById(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var (
		method, ret string
		vars        = mux.Vars(r)
		vars_id     = vars["id"]
		query       elastic.Query
		conf        *config.ConfigVars
		err         error
		client      *elastic.Client
		res         *elastic.SearchResult
	)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if conf.Elastic.Elastic_ssl {
		method = api.HTTPS
	} else {
		method = api.HTTP
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s%s:%d", method, conf.Elastic.Elastic_server, conf.Elastic.Elastic_port)))
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	query = elastic.NewTermQuery("_id", vars_id)

	res, err = client.Search().
		Index(conf.Elastic.Elastic_index).
		Query(query).
		Size(1).
		Do(context.Background())

	ret = api.HttpSuccessMessage("200", &res.Hits.Hits[0], res.TotalHits())
	api.LogDebug(api.DEBUG, "[+] GET /watcher/results/id HTTP 200, returned watcher result id doc")
	fmt.Fprintln(w, ret)
}

func GenerateWatcherRule(w http.ResponseWriter, r *http.Request) {
	var (
		rb api.RB
		rd api.RuleDefaults
	)

	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if err := r.Body.Close(); err != nil {
		api.LogError(api.DEBUG, err)
	}

	json.Unmarshal(body, &rb)

	conf, err := config.ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}
	// Setup RuleDefaults
	m := make(map[interface{}]interface{}) // map is needed to avoid struct names as becoming yaml

	rd.ESHost = conf.Elastic.Elastic_server
	rd.ESPort = conf.Elastic.Elastic_port
	rd.Index = conf.Elastic.Elastic_index
	rd.Timestampfield = "Record.TlnTime"

	switch rb.RealertTimelen {
	case "minutes":
		m["realert"] = struct {
			Minutes int
		}{
			rb.RealertDuration,
		}
		break
	case "hours":
		m["realert"] = struct {
			Hours int
		}{
			rb.RealertDuration,
		}
		break
	case "days":
		m["realert"] = struct {
			Days int
		}{
			rb.RealertDuration,
		}
	}

	// make the rule
	m["es_host"] = rd.ESHost
	m["es_port"] = rd.ESPort
	m["index"] = rd.Index
	m["timestamp_field"] = rd.Timestampfield
	m["name"] = rb.Name
	m["type"] = rb.Type

	m["compare_key"] = rb.RuleMeta.CompareKey
	switch rb.Type {
	case "blacklist":
		m["blacklist"] = rb.RuleMeta.ListTerms
		break
	case "whitelist":
		m["whitelist"] = rb.RuleMeta.ListTerms
		break
	}

	data, err := yaml.Marshal(&m)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(data))
}
