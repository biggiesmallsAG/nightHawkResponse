package log

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	nhc "nighthawk/common"
	nhconfig "nighthawk/config"
	"nighthawk/rabbitmq"
	"nighthawklogger/logger"

	"github.com/streadway/amqp"
)

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

var rconfig rabbitmq.RabbitMQConfig
var conn *amqp.Connection = nil
var ch *amqp.Channel = nil

func init() {
	if nhconfig.IsStandalone() {
		return
	}

	// Send logs to RabbitMQ
	rconfig = rabbitmq.LoadRabbitMQConfig(rabbitmq.RABBITMQ_CONFIG_FILE)
	conn = rabbitmq.Connect(rconfig.Server)

	tch, err := conn.Channel()
	if err != nil {
		nhc.ConsoleMessage("init", "ERROR", fmt.Sprintf("Failed to connect to channel. %s", err.Error()), true)
		os.Exit(nhc.ERROR_CHANNEL_CONNECT)
	}
	ch = tch
}

func LogMessage(sender string, level string, message string) {
	// Only send logs to subsystem if log level
	// matched with configuration file
	lsl := logSeverityLevel(level)
	if lsl > nhconfig.VerboseLevel() {
		return
	}

	if nhconfig.IsStandalone() {
		if nhconfig.Verbose() {
			nhc.ConsoleMessage(sender, level, message, true)
		}
		return // No further processing required in standalone mode
	}

	// Send logs to RabbitMQ
	// Initialize worker queue
	rabbitmq.RabbitQueueDeclare(ch, rconfig.Logger)
	msg := logger.Logger{
		Timestamp: time.Now().UTC(),
		LogLevel:  level,
		Body:      message,
	}

	data, _ := json.Marshal(&msg)
	err := rabbitmq.RabbitQueuePublisher(ch, rconfig.Logger, data)
	if err != nil {
		nhc.ConsoleMessage("LogMessage", "ERROR", fmt.Sprintf("Failed to pusblish log message. %s", err.Error()), true)
	}
}

func CreateJobMessage(message interface{}) {
	msg := logger.Logger{
		Timestamp: time.Now().UTC(),
		LogLevel:  "JOB",
		Body:      message,
	}

	data, _ := json.Marshal(&msg)
	err := rabbitmq.RabbitQueuePublisher(ch, rconfig.Jobqueue, data)
	if err != nil {
		nhc.ConsoleMessage("CreateJobMessage", "ERROR", fmt.Sprintf("Failed to publish job message. %s", err.Error()), true)
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
