package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
)

type AuditParserConfig struct {
	agentstate     bool
	arp            bool
	disk           bool
	dns            bool
	eventlogs      bool
	filedlhistory  bool
	files          bool
	persistence    bool
	ports          bool
	prefetch       bool
	process_memory bool
	processtree    bool
	registry       bool
	route          bool
	services       bool
	system         bool
	tasks          bool
	urlhistory     bool
	useraccounts   bool
	volumes        bool
}

func (pc *AuditParserConfig) SetDefault() {
	pc.agentstate = true
	pc.arp = true
	pc.disk = true
	pc.dns = true
	pc.eventlogs = true
	pc.filedlhistory = true
	pc.files = true
	pc.persistence = true
	pc.ports = true
	pc.prefetch = true
	pc.process_memory = true
	pc.processtree = true
	pc.registry = true
	pc.route = true
	pc.services = true
	pc.system = true
	pc.tasks = true
	pc.urlhistory = true
	pc.useraccounts = true
	pc.volumes = true
}

type NighthawkConfig struct {
	MaxProcs       int  `json:"max_procs"`
	MaxGorouting   int  `json:"max_goroutine"`
	BulkPostSize   int  `json:"bulk_post_size"`
	OpControl      int  `json:"opcontrol"`
	SessionDirSize int  `json:"sessiondir_size"`
	CheckHashSet   bool `json:"check_hash"`
	CheckStack     bool `json:"check_stack"`
	Verbose        bool `json:"verbose"`
	VerboseLevel   int  `json:"verbose_level,omitempty"`
	Standalone     bool `json:"standalone"`
}

type ElasticConfig struct {
	Host     string `json:"elastic_server"`
	Port     int    `json:"elastic_port"`
	Index    string `json:"elastic_index"`
	Ssl      bool   `json:"elastic_ssl"`
	User     string `json:"elastic_user"`
	Pass     string `json:"elastic_pass"`
	Authcode string `json:"-"`
}

type NHRConfig struct {
	Nighthawk    NighthawkConfig   `json:"nighthawk"`
	Elastic      ElasticConfig     `json:"elastic"`
	ParserConfig AuditParserConfig `json:"parser_config,omitempty"`
}

func (config *NHRConfig) SetVerbose() {
	config.Nighthawk.Verbose = true
}

func (config *NHRConfig) UnsetVerbose() {
	config.Nighthawk.Verbose = false
}

func LoadConfigFile(configfile string) (NHRConfig, error) {
	var nhrconfig NHRConfig
	nhrconfig.LoadDefaultConfig()

	configData, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nhrconfig, errors.New("Error Opening config file")
	}

	json.Unmarshal(configData, &nhrconfig)
	nhrconfig.Elastic.Authcode = base64.StdEncoding.EncodeToString([]byte(nhrconfig.Elastic.User + ":" + nhrconfig.Elastic.Pass))

	if nhrconfig.Nighthawk.VerboseLevel == 0 {
		nhrconfig.Nighthawk.VerboseLevel = 6
	}

	nhrconfig.ParserConfig.SetDefault()
	return nhrconfig, nil
}

func (config *NHRConfig) SaveConfigFile(configfile string) error {
	json_config, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return errors.New("Error Marshaling configuration object")
	}

	return ioutil.WriteFile(configfile, json_config, 0644)
}

func (config *NHRConfig) LoadDefaultConfig() {
	config.Nighthawk.MaxProcs = 2
	config.Nighthawk.MaxGorouting = 100
	config.Nighthawk.BulkPostSize = 1000
	config.Nighthawk.OpControl = 2
	config.Nighthawk.SessionDirSize = 8
	config.Nighthawk.CheckHashSet = false
	config.Nighthawk.CheckStack = false
	config.Nighthawk.Verbose = false
	config.Nighthawk.VerboseLevel = 4
	config.Nighthawk.Standalone = true
	config.Elastic.Host = "127.0.0.1"
	config.Elastic.Port = 9200
	config.Elastic.Index = "investigations"
	config.Elastic.Ssl = false
}

/**********************************************************************
** Functions Exported for External Access
***********************************************************************/

func NHRConfigObject() *NHRConfig {
	return &nhrconfig
}

func MaxProcs() int {
	return nhrconfig.Nighthawk.MaxProcs
}

func MaxGoRoutine() int {
	return nhrconfig.Nighthawk.MaxGorouting
}

func SessionDirSize() int {
	return nhrconfig.Nighthawk.SessionDirSize
}

func GetVerbose() bool {
	return nhrconfig.Nighthawk.Verbose
}

func SetVerbose() {
	nhrconfig.Nighthawk.Verbose = true
}

func Verbose() bool {
	return nhrconfig.Nighthawk.Verbose
}

func VerboseLevel() int {
	return nhrconfig.Nighthawk.VerboseLevel
}

func IsStandalone() bool {
	return nhrconfig.Nighthawk.Standalone
}

func OpControl() int {
	return nhrconfig.Nighthawk.OpControl
}

func CheckHashSet() bool {
	return nhrconfig.Nighthawk.CheckHashSet
}

func CheckStack() bool {
	return nhrconfig.Nighthawk.CheckStack
}

func BulkPostSize() int {
	return nhrconfig.Nighthawk.BulkPostSize
}

func ElasticHost() string {
	return nhrconfig.Elastic.Host
}

func ElasticPort() int {
	return nhrconfig.Elastic.Port
}

func ElasticIndex() string {
	return nhrconfig.Elastic.Index
}

func ElasticUser() string {
	return nhrconfig.Elastic.User
}

func ElasticPassword() string {
	return nhrconfig.Elastic.Pass
}

func ElasticSSL() bool {
	return nhrconfig.Elastic.Ssl
}

func ElasticBasicAuth() string {
	return nhrconfig.Elastic.Authcode
}

func ElasticHttpScheme() string {
	if ElasticSSL() {
		return "https"
	}
	return "http"
}

func SetAuthCode(authcode string) {
	nhrconfig.Elastic.Authcode = authcode
}

func ParserConfigSetting(auditname string) bool {
	switch auditname {
	case "agentstate":
		return nhrconfig.ParserConfig.agentstate
	case "arp":
		return nhrconfig.ParserConfig.arp
	case "disk":
		return nhrconfig.ParserConfig.disk
	case "dns":
		return nhrconfig.ParserConfig.dns
	case "eventlogs":
		return nhrconfig.ParserConfig.eventlogs
	case "filedlhistory":
		return nhrconfig.ParserConfig.filedlhistory
	case "files":
		return nhrconfig.ParserConfig.files
	case "ports":
		return nhrconfig.ParserConfig.ports
	case "prefetch":
		return nhrconfig.ParserConfig.prefetch
	case "process_memory":
		return nhrconfig.ParserConfig.process_memory
	case "processtree":
		return nhrconfig.ParserConfig.processtree
	case "registry":
		return nhrconfig.ParserConfig.registry
	case "route":
		return nhrconfig.ParserConfig.route
	case "services":
		return nhrconfig.ParserConfig.services
	case "system":
		return nhrconfig.ParserConfig.system
	case "tasks":
		return nhrconfig.ParserConfig.tasks
	case "urlhistory":
		return nhrconfig.ParserConfig.urlhistory
	case "useraccounts":
		return nhrconfig.ParserConfig.useraccounts
	case "volumes":
		return nhrconfig.ParserConfig.volumes
	}
	return true
}
