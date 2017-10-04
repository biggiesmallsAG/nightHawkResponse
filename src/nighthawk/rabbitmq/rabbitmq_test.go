package rabbitmq

import (
	"fmt"
	"testing"
)


func TestRabbitMQConfig(t *testing.T) {
	const rabbitmq_config_file = "/opt/nighthawk/etc/rabbitmq.json"

	rconfig := LoadRabbitMQConfig(rabbitmq_config_file)
	fmt.Println(rconfig.Server)
	conn := Connect(rconfig.Server)

	ch, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}

	huntQ := RabbitQueueDeclare(ch, rconfig.Hunter)
	fmt.Println(huntQ)

	msgs := RabbitQueueConsumer(ch, rconfig.Hunter)
	fmt.Println(msgs)

}