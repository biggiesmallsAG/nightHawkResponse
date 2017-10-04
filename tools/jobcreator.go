/*
 * This tools creates job message and 
 * publishes to RabbitMQ server
 *
 * Syntax: jobcreator -N CASE-DELTA -f /opt/nighthawk/var/media/mytriage.mans
 *
 */

package main 

import (
	"fmt"
	"flag"
	"time"
	"encoding/json"

	"nighthawk/rabbitmq"
	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/upload"
)



func main() {
	var CaseName string
	var TriageFile string 

	flag.StringVar(&CaseName,"N","", "Casename")
	flag.StringVar(&TriageFile,"f", "", "Triage file")
	flag.Parse()

	var af []string 
	af = append(af, TriageFile)

	uuid := api.GenUID()

	rconfig := rabbitmq.LoadRabbitMQConfig(rabbitmq.RABBITMQ_CONFIG_FILE)
	uploadq := rabbitmq.Connect(rconfig.Server)
	ch, err := uploadq.Channel()
	if err != nil {
		panic(err.Error())
	}
	defer uploadq.Close()

	job_msg := upload.Job{
		UID: uuid,
		TS: time.Now().UTC().Format(api.LAYOUT),
		CaseID: CaseName,
		Audits: af,
		InProg: false,
		Complete: false,
		Cancelled: false,
	}

	jb,_ := json.Marshal(&job_msg)
	err = rabbitmq.RabbitQueuePublisher(ch, rconfig.Worker, jb)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created new job %s\n", uuid)

}
