package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Log level constant
const (
	EMERG   = 0
	ALERT   = 1
	CRIT    = 2
	ERROR   = 3
	WARNING = 4
	NOTICE  = 5
	INFO    = 6
	DEBUG   = 7
)

const nhrIndex = "nighthawk"
const nhrLog = "log"

var (
	client *elastic.Client
	err    error
)

func init() {
	// 0xredskull: Removing logs from RabbitMQ to Elasticsearch
	// Logs are send to Elasticsearch index: nighthawk, type: log
	conf, err := nhconfig.LoadConfigFile(nhconfig.CONFIG_FILE)
	if err != nil {
		nhc.ConsoleMessage("log.go", "ERROR", err.Error(), true)
	}

	// Elasticsearch client initialization
	httpSchema := "http"
	if conf.Elastic.Ssl {
		httpSchema = "https"
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", httpSchema, conf.Elastic.Host, conf.Elastic.Port)))
	if err != nil {
		nhc.ConsoleMessage("log.go", "ERROR", err.Error(), true)
		return
	}
}

// NHRLogMessage is common log structure
type NHRLogMessage struct {
	Timestamp string      // Timestamp of when log was generated
	Source    string      // Name of nighthawk component that generated log message
	Pid       int         // ProcessID of Source
	FuncName  string      // Function that generated this log
	Level     string      // Log level
	Message   interface{} // Log message
}

// LogMessage writes function Elasticsearch
func LogMessage(sender string, level string, message string) {
	logLevel := logSeverityLevel(level)

	//// Console printing logs
	if logLevel <= nhconfig.VerboseLevel() {
		nhc.ConsoleMessage(sender, level, message, true)
	}

	// If elasticsearch client is not initialized
	// do not proceed indexing
	//
	// If logLevel is above Error do not log in elasticsearch -- save storage
	if client == nil || logLevel >= INFO {
		return
	}
	_, procname := filepath.Split(os.Args[0])

	logentry := NHRLogMessage{
		Timestamp: time.Now().UTC().Format(nhc.Layout),
		Source:    procname,
		Pid:       os.Getpid(),
		FuncName:  sender,
		Level:     level,
		Message:   message,
	}

	jsonLogEntry, _ := json.Marshal(logentry)

	res, err := client.Index().Index(nhrIndex).Type(nhrLog).BodyJson(string(jsonLogEntry)).Do(context.Background())

	if err != nil {
		fmt.Printf("LogMessage - %s\n", err.Error())
		return
	}

	if res.Result != "created" {
		fmt.Printf("LogMessage - index status %s", res.Result)
		return
	}

}

// JobMessage passes MQ messages
func JobMessage(sender string, level string, message interface{}) {

	if client == nil {
		return
	}
	_, procname := filepath.Split(os.Args[0])

	logentry := NHRLogMessage{
		Timestamp: time.Now().UTC().Format(nhc.Layout),
		Source:    procname,
		Pid:       os.Getpid(),
		FuncName:  sender,
		Level:     level,
		Message:   message,
	}

	jsonLogEntry, _ := json.Marshal(logentry)

	res, err := client.Index().Index(nhrIndex).Type(nhrLog).BodyJson(string(jsonLogEntry)).Do(context.Background())

	if err != nil {
		fmt.Printf("LogMessage - %s\n", err.Error())
		return
	}

	if res.Result != "created" {
		fmt.Printf("JobMessage - index status %s", res.Result)
	}
}

func logSeverityLevel(level string) int {
	switch level {
	case "EMERG":
		return EMERG
	case "ALERT":
		return ALERT
	case "CRIT":
		return CRIT
	case "ERROR":
		return ERROR
	case "WARNING":
		return WARNING
	case "NOTICE":
		return NOTICE
	case "INFO":
		return INFO
	case "DEBUG":
		return DEBUG
	}

	// return default to INFO
	return INFO
}
