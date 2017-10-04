/* nighthawk.rabbitmq.rabbitmq.go
 * author: 0xredskull
 *
 * RabbitMQ messaging
 */

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	nhc "nighthawk/common"
	nhconfig "nighthawk/config"

	"github.com/streadway/amqp"
)

const RABBITMQ_CONFIG_FILENAME = "rabbitmq.json"

var RABBITMQ_CONFIG_FILE string = filepath.Join(nhconfig.CONFDIR, RABBITMQ_CONFIG_FILENAME)

type ServerConfig struct {
	Host     string `json:"host"`
	Protocol string `json:"protocol"` // amqp || amqps
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type QueueConfig struct {
	Process    int        `json:"process,omitempty"` // Number of worker processes
	Name       string     `json:"name,omitempty"`
	RoutingKey string     `json:"routing_key,omitempty"` // Name of Message Queue
	Exchange   string     `json:"exchange,omitempty"`
	AutoAck    bool       `json:"auto_ack,omitempty"`
	Durable    bool       `json:"durable,omitempty"`
	AutoDelete bool       `json:"auto_delete,omitempty"`
	Exclusive  bool       `json:"exclusive,omitempty"`
	Wait       bool       `json:"no_wait,omitempty"`
	Local      bool       `json:"no_local,omitempty"`
	Mandatory  bool       `json:"mandatory,omitempty"`
	Immediate  bool       `json:"immediate,omitempty"`
	Argument   amqp.Table `json:"-"`
}

type ExchangeConfig struct {
	Name       string     `json:"name,omitempty"`
	Type       string     `json:"type,omitempty"`
	Durable    bool       `json:"durable,omitempty"`
	AutoDelete bool       `json:"auto_delete,omitempty"`
	Internal   bool       `json:"internal,omitempty"`
	Wait       bool       `json:"no_wait,omitempty"`
	Arguments  amqp.Table `json:"-"`
}

type RabbitMQConfig struct {
	Server            ServerConfig   `json:"server"`
	LogTopicExchange  ExchangeConfig `json:"nighthawk_logexchange"`
	WorkTopicExchange ExchangeConfig `json:"nighthawk_workexchange"`
	Worker            QueueConfig    `json:"worker"`
	Jobqueue          QueueConfig    `json:"jobqueue"`
	Hunter            QueueConfig    `json:"hunter"`
	Logger            QueueConfig    `json:"logger"`
}

func LoadRabbitMQConfig(config_file string) RabbitMQConfig {
	config_data, err := ioutil.ReadFile(config_file)
	if err != nil {
		nhc.ConsoleMessage("LoadRabbitMQConfig", "ERROR", "Error opening RabbitMQ configuration file", true)
		//if nhconfig.VerboseLevel() == 7 {
		//	panic(err.Error())
		//}
	}
	var rconfig RabbitMQConfig
	json.Unmarshal(config_data, &rconfig)

	return rconfig
}

func Connect(config ServerConfig) *amqp.Connection {
	rabbit_server := fmt.Sprintf("%s://%s:%s@%s:%d/", config.Protocol, config.User, config.Password, config.Host, config.Port)
	conn, err := amqp.Dial(rabbit_server)
	if err != nil {
		panic(err.Error())
	}

	return conn
}

/*
 * This function returns an exchange binding
 */
func RabbitExchangeDeclare(ch *amqp.Channel, config ExchangeConfig) error {
	err := ch.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.Wait,
		config.Arguments)

	return err
}

/*
 * This function returns nighthawk worker queue
 */
func RabbitQueueDeclare(ch *amqp.Channel, config QueueConfig) amqp.Queue {
	q, err := ch.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.Wait,
		nil,
	)

	if err != nil {
		panic(err.Error())
	}

	return q
}

/*
 * RabbitMQ Bind wrapper function
 */
func RabbitBindQueue(ch *amqp.Channel, config QueueConfig) error {
	err := ch.QueueBind(
		config.Name,
		config.RoutingKey,
		config.Exchange,
		config.Wait,
		config.Argument)

	return err
}

/*
 * RabbitMQ Consumer wrapper function
 */
func RabbitQueueConsumer(ch *amqp.Channel, config QueueConfig) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		config.Name,
		config.Exchange,
		config.AutoAck,
		config.Exclusive,
		config.Local,
		config.Wait,
		nil,
	)

	if err != nil {
		panic(err.Error())
	}

	return msgs
}

/*
 * RabbitMQ Publisher wrapper function
 */

func RabbitQueuePublisher(ch *amqp.Channel, config QueueConfig, body []byte) error {
	err := ch.Publish(
		config.Exchange,
		config.RoutingKey,
		config.Mandatory,
		config.Immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	return err
}
