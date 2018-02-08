package upload

import (
	"log"
	"net/http"
	api "nighthawkresponse/api/core"
	"nighthawkresponse/rabbitmq"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func SubscribeJobs(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	rconfig := rabbitmq.LoadRabbitMQConfig(rabbitmq.RABBITMQ_CONFIG_FILE)
	uploadq := rabbitmq.Connect(rconfig.Server)

	defer uploadq.Close()
	ch, err := uploadq.Channel()
	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	defer ch.Close()

	// Init worker queue
	_ = rabbitmq.RabbitQueueDeclare(ch, rconfig.Jobqueue)
	// Bind to exchange
	_ = rabbitmq.RabbitBindQueue(ch, rconfig.Jobqueue)

	job := rabbitmq.RabbitQueueConsumer(ch, rconfig.Jobqueue)

	client := &api.Client{Id: api.GenUID(), Socket: c, Send: make(chan []byte)}

	api.Manager.Register <- client

	go client.Read()
	go client.Write()

	forever := make(chan bool)

	go func() {
		for m := range job {
			api.Manager.Broadcast <- m.Body
		}
	}()

	api.LogDebug(api.DEBUG, "[+] GET /subscribe/uploadjobs WS 101, inside websocket job queue loop.")
	<-forever
}
