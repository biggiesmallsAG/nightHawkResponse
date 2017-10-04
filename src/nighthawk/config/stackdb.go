package config

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type StackDbConfig struct {
	Index           string
	Path            string
	LookupEnabled   bool
	LookupAvailable bool
	LookupChecked   bool
}

func (config *StackDbConfig) LoadDefaultConfig() {
	config.Index = "nighthawk.db"
	config.Path = filepath.Join(CONFDIR, config.Index)
	config.LookupEnabled = false
	config.LookupAvailable = false
	config.LookupChecked = false
}

func (config *StackDbConfig) Initialize() {
	//sdconfig.Path = filepath.Join(CONFDIR, sdconfig.Index)
	tDb, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		//sdconfig.LookupAvailable = false
		config.LookupChecked = true
		fmt.Println("StackDbConfig::Initialize - Failed to open nighthakw.db - ", err.Error())
		fmt.Println("StackDbConfig - Disabling StackDb Lookup")
		return
	}

	config.LookupEnabled = true
	config.LookupAvailable = true
	config.LookupChecked = true
	nhdb = tDb
}

func StackDbObject() *sql.DB {
	return nhdb
}

func StackDbIndex() string {
	return sdconfig.Index
}

func StackDbPath() string {
	return sdconfig.Path
}

func StackDbEnabled() bool {
	return sdconfig.LookupEnabled
}

func StackDbAvailable() bool {
	return sdconfig.LookupAvailable
}

func StackDbChecked() bool {
	return sdconfig.LookupChecked
}
