package upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nighthawklogger/config"
)

const (
	job_query = `
	SELECT Timestamp, Loglevel, Worker, Body FROM logs
	WHERE Loglevel = "JOB"
	ORDER BY DATETIME(Timestamp) Desc
	`
)

func ListCompletedJobs(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	db := config.InitDB()
	jobs := config.ReadLogs(db, job_query)

	var (
		j  Job
		_j []Job
	)

	for _, job := range jobs {
		json.Unmarshal(job.Body.([]byte), &j)
		if j.Complete {
			_j = append(_j, j)
		}
	}
	ret, _ := json.MarshalIndent(&_j, "", "    ")
	fmt.Fprintln(w, string(ret))
}
