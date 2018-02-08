package upload

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	api "nighthawkresponse/api/core"
	"nighthawkresponse/rabbitmq"
	"path/filepath"
	"time"
)

type Job struct {
	UID          string   `json:"uid"`
	TS           string   `json:"timestamp"`
	UserID       string   `json:"user_id"`
	CaseID       string   `json:"case_id"`
	ComputerName string   `json:"computer_name"`
	Audits       []string `json:"audit_file"`
	InProg       bool     `json:"in_progress"`
	Complete     bool     `json:"is_complete"`
	Cancelled    bool     `json:"is_cancelled"`
}

func JobDispatch(f *multipart.Form, c []string) {
	api.LogDebug(api.DEBUG, fmt.Sprintf("[-] Job dispatch recieved %d files to upload.", len(f.File)))

	// Make UID from random input
	uuid := api.GenUID()

	// Get MQ Client
	rconfig := rabbitmq.LoadRabbitMQConfig(rabbitmq.RABBITMQ_CONFIG_FILE)
	uploadq := rabbitmq.Connect(rconfig.Server)
	ch, err := uploadq.Channel()

	defer uploadq.Close()

	if err != nil {
		api.LogError(api.DEBUG, err)
	}

	// Get File Array
	var af []string
	for i, _ := range f.File {
		af = append(af, filepath.Join(api.MEDIA_DIR, i))
	}

	job_msg := Job{
		UID:       uuid,
		TS:        time.Now().UTC().Format(api.LAYOUT),
		CaseID:    c[0],
		Audits:    af,
		InProg:    false,
		Complete:  false,
		Cancelled: false,
	}

	_jm, _ := json.Marshal(&job_msg)

	err = rabbitmq.RabbitQueuePublisher(ch, rconfig.Worker, _jm)

	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] Dispatched %d jobs.", len(f.File)))
}
