/*
	nightHawkAPI.handler.config;
*/

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkapi/api/core"
)

type _Body struct {
	Cvars ConfigVars `json:"name"`
}

type ConfigVars struct {
	NightHawk NightHawkConf `json:"nighthawk"`
	Elastic   ElasticConf   `json:"elastic"`
}

type NightHawkConf struct {
	Ip_addr          string `json:"ip_addr"`
	Max_procs        int    `json:"max_procs"`
	Max_goroutine    int    `json:"max_goroutine"`
	Max_file_upload  int    `json:"max_file_upload"`
	Bulk_post_size   int    `json:"bulk_post_size"`
	Opcontrol        int    `json:"opcontrol"`
	Session_dir_size int    `json:"sessiondir_size"`
	Check_hash       bool   `json:"check_hash"`
	Check_stack      bool   `json:"check_stack"`
	Verbose          bool   `json:"verbose"`
	Verbose_level    int    `json:"verbose_level"`
}

type ElasticConf struct {
	Elastic_server string `json:"elastic_server"`
	Elastic_port   int    `json:"elastic_port"`
	Elastic_ssl    bool   `json:"elastic_ssl"`
	Elastic_user   string `json:"elastic_user"`
	Elastic_pass   string `json:"elastic_pass"`
	Elastic_index  string `json:"elastic_index"`
}

var configObject *ConfigVars

func init() {
	conf, _ := ReadConfFile()
	configObject = conf
}

func ReadConfFile() (c *ConfigVars, err error) {
	var conf ConfigVars
	configData, err := ioutil.ReadFile(fmt.Sprintf("%s%s", api.WORKING_DIR, api.CONF_FILE))

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	json.Unmarshal(configData, &conf)
	return &conf, err
}

func (config *ConfigVars) ServerHttpScheme() string {
	if config.Elastic.Elastic_ssl {
		return "https"
	}
	return "http"
}

func (config *ConfigVars) ServerHost() string {
	return config.Elastic.Elastic_server
}

func (config *ConfigVars) ServerPort() int {
	return config.Elastic.Elastic_port
}

func (config *ConfigVars) ServerIndex() string {
	return config.Elastic.Elastic_index
}

func ReturnSystemConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var ret string
	conf, err := ReadConfFile()

	if err != nil {
		api.LogError(api.DEBUG, err)
		ret = api.HttpFailureMessage("Failed to read config")
		if err != nil {
			api.LogError(api.DEBUG, err)
		}
		fmt.Fprintln(w, string(ret))
	}

	ret = api.HttpSuccessMessage("200", conf, 0)
	api.LogDebug(api.DEBUG, "[+] GET /config HTTP 200, config returned.")
	fmt.Fprintln(w, string(ret))
}

func UpdateSystemConfig(w http.ResponseWriter, r *http.Request) {
	var (
		conf _Body
		ret  string
	)

	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	if err := r.Body.Close(); err != nil {
		api.LogError(api.DEBUG, err)
	}

	json.Unmarshal(body, &conf)

	fin, err := json.MarshalIndent(&conf.Cvars, "", "    ")

	err = ioutil.WriteFile(fmt.Sprintf("%s%s", api.WORKING_DIR, api.CONF_FILE), fin, 0644)
	if err != nil {
		api.LogError(api.DEBUG, err)
		ret = api.HttpFailureMessage("Failed to write config")

		if err != nil {
			api.LogError(api.DEBUG, err)
		}
		fmt.Fprintln(w, string(ret))
	}

	ret = api.HttpSuccessMessage("200", fin, 0)
	api.LogDebug(api.DEBUG, "[+] POST /config HTTP 200, config updated.")
	fmt.Fprintln(w, string(ret))
}

func ElasticUrl() string {
	url := ""
	conf, err := ReadConfFile()
	if err != nil {
		api.LogError(api.DEBUG, err)
		return url
	}

	url = fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())
	return url
}

func ElasticIndex() string {
	return configObject.Elastic.Elastic_index
}
