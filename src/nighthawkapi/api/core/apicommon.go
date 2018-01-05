/*
	nightHawkAPI.common;
*/

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	BASEDIR   = "/opt/nighthawk"
	CONFDIR   = ""
	STATEDIR  = ""
	MEDIA_DIR = ""
	CONF_FILE = ""
)

const (
	VERSION    = "1.0-rc"
	API_NAME   = "/api"
	API_VER    = "/v1"
	HTTP       = "http://"
	HTTPS      = "https://"
	UPLOAD_MEM = 100000
	LAYOUT     = "2006-01-02T15:04:05Z"
)

var (
	DEBUG  = false
	STDOUT = false
)

func init() {
	if runtime.GOOS == "windows" {
		BASEDIR = "C:\\ProgramData\\nighthawk"
	}

	CONFDIR = filepath.Join(BASEDIR, "etc")
	STATEDIR = filepath.Join(BASEDIR, "var")
	MEDIA_DIR = filepath.Join(STATEDIR, "media")
	CONF_FILE = filepath.Join(CONFDIR, "nighthawk.json")
}

type Err struct {
	Reason   string `json:"reason"`
	Response string `json:"response"`
}

type Success struct {
	Reason   string      `json:"reason"`
	Response string      `json:"response"`
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
}

func HttpFailureMessage(reason string) string {
	s := Err{
		Reason:   reason,
		Response: "failed",
	}
	data, _ := json.MarshalIndent(&s, "", "    ")
	return string(data)
}

func HttpSuccessMessage(reason string, data interface{}, total int64) string {
	s := Success{
		Reason:   reason,
		Response: "success",
		Data:     &data,
		Total:    total,
	}
	ret, _ := json.MarshalIndent(&s, "", "    ")
	return string(ret)
}

type RuleDefaults struct {
	ESHost         string `yaml:"es_host"`
	ESPort         int    `yaml:"es_port"`
	Timestampfield string `yaml:"timestamp_field"`
	Index          string `yaml:"index"`
}

type RB struct {
	RuleBase `json:"ruleBase"`
}

type RuleBase struct {
	Type            string `json:"rule_type"`
	Name            string `json:"rule_name"`
	RealertDuration int    `json:"realert_duration,omitempty"`
	RealertTimelen  string `json:"realert_timelength,omitempty"`
	RuleMeta        struct {
		ListTerms  []string `json:"list_terms"`
		CompareKey []string `json:"compare_key"`
	} `json:"rule_meta"`
}

type BLWLRule struct {
	Type    string `yaml:"type"`
	Name    string `yaml:"name"`
	Realert struct {
		Minutes int `yaml:"minutes,omitempty"`
		Days    int `yaml:"days,omitempty"`
		Hours   int `yaml:"hours,omitempty"`
	} `yaml:"realert"`
	Blacklist  []string `yaml:"blacklist"`
	CompareKey []string `yaml:"compare_key"`
}

type WatcherResults struct {
	Timestamp string `json:"alert_time"`
	MatchBody struct {
		CaseInfo struct {
			CaseAnalyst string `json:"case_analyst"`
			CaseDate    string `json:"case_date"`
			Endpoint    string `json:"computer_name"`
			CaseName    string `json:"case_name"`
		}
		NumMatches int    `json:"num_matches"`
		Id         string `json:"_id"`
	} `json:"match_body"`
	RuleName string `json:"rule_name"`
}

func LogError(DEBUG bool, err error) {
	if DEBUG {
		log.Printf("ERROR - %s\n", err)
	}
}

func LogConsole(STDOUT bool, message string) {
	if STDOUT {
		fmt.Println(message)
	}
}

func LogDebug(DEBUG bool, message string) {
	if DEBUG {
		log.Printf("DEBUG - %s\n", message)
	}
}

func GenUID() string {
	file, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	file.Read(b)
	file.Close()
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
