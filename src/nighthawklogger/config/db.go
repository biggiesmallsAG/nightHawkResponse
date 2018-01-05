package config

import (
	"database/sql"
	"encoding/json"
	"log"
	nhconfig "nighthawk/config"
	"nighthawklogger/logger"

	_ "github.com/mattn/go-sqlite3"
)

type DBFactory interface {
	InitDB() *sql.DB
	CreateTable(db *sql.DB)
	StoreLogs(db *sql.DB, logs []logger.Logger)
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", nhconfig.DBDIR+"/"+nhconfig.LOGGER_DB_NAME)
	if err != nil {
		log.Panic(err)
	}
	return db
}

func CreateTable(db *sql.DB) {
	table_schema := `
	CREATE TABLE IF NOT EXISTS logs(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Timestamp DATETIME NOT NULL,
		Worker TEXT,
		Loglevel TEXT NOT NULL,
		Body TEXT NOT NULL)
	`

	_, err := db.Exec(table_schema)
	if err != nil {
		log.Panic(err)
	}

}

func StoreLogs(db *sql.DB, logs *logger.Logger) {
	log_insert := `
	INSERT INTO logs(
		Timestamp,
		Worker,
		Loglevel,
		Body) values (?, ?, ?, ?)
	`

	p, err := db.Prepare(log_insert)
	if err != nil {
		log.Panic(err)
	}

	defer p.Close()

	switch logs.Body.(type) {
	case map[string]interface{}:
		message_body, err := json.Marshal(&logs.Body)
		_, err = p.Exec(logs.Timestamp, logs.Worker, logs.LogLevel, message_body)
		if err != nil {
			log.Println(err)
		}
	default:
		_, err = p.Exec(logs.Timestamp, logs.Worker, logs.LogLevel, logs.Body)
		if err != nil {
			log.Println(err)
		}
	}
}

func ReadLogs(db *sql.DB, read_query string) []logger.Logger {
	rows, err := db.Query(read_query)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	var res []logger.Logger
	for rows.Next() {
		_log := logger.Logger{}
		err := rows.Scan(&_log.Timestamp, &_log.LogLevel, &_log.Worker, &_log.Body)
		if err != nil {
			log.Println(err)
		}
		res = append(res, _log)
	}

	return res
}
