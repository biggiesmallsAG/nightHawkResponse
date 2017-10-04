package logger

import (
	"nighthawk/rabbitmq"
	"time"

	"github.com/streadway/amqp"
)

type Logger struct {
	Timestamp time.Time   `json:"timestamp"`
	LogLevel  string      `json:"log_level"`
	Worker    string      `json:"worker"`
	Body      interface{} `json:"body"`
}

type LoggerFactory interface {
	InitMQLogger()
	ConsumeMQLogger(ch *amqp.Channel, rconfig *rabbitmq.RabbitMQConfig)
}

func InitMQLogger() (*amqp.Channel, rabbitmq.RabbitMQConfig) {
	rconfig := rabbitmq.LoadRabbitMQConfig(rabbitmq.RABBITMQ_CONFIG_FILE)
	conn := rabbitmq.Connect(rconfig.Server)
	ch, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	return ch, rconfig
}

func ConsumeMQLogger(ch *amqp.Channel, rconfig *rabbitmq.RabbitMQConfig) <-chan amqp.Delivery {
	_ = rabbitmq.RabbitQueueDeclare(ch, rconfig.Logger)

	messages := rabbitmq.RabbitQueueConsumer(ch, rconfig.Logger)
	return messages
}
