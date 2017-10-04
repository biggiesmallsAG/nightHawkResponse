/*
 * nighthawk.config.env
 *
 */

package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Global Environmental Variables
// Directory Setup
var (
	BASEDIR     = "/opt/nighthawk"
	CONFDIR     = ""
	STATEDIR    = ""
	DBDIR       = ""
	MEDIA       = ""
	TMP         = ""
	RUNDIR      = ""
	WORKSPACE   = ""
	CONFIG_FILE = ""
)

// STANALONE: Run binary without MessageQueue and database
// nighthawk binary can be run in standalone mode
// i.e. Binary directly interacting with Elasticsearch
// and not using RabbitMQ for notification
//var STANDALONE bool = true

var nhrconfig NHRConfig
var hsconfig HashSetConfig
var sdconfig StackDbConfig
var nhdb *sql.DB = nil // nighthawk stack database

const (
	CONFIG_FILE_NAME = "nighthawk.json"
	LOGGER_DB_NAME   = "logger.sqlite3"
)

// Controlling output
const (
	OP_CONSOLE_ONLY      = 1
	OP_DATASTORE_ONLY    = 2
	OP_CONSOLE_DATASTORE = 3
	OP_SPLUNK_FILE       = 4
	OP_SPLUNK            = 5
)

// Initialze all global configuration
func init() {
	if runtime.GOOS == "windows" {
		BASEDIR = "C:\\ProgramData\\nighthawk"
	}

	CONFDIR = filepath.Join(BASEDIR, "etc")
	STATEDIR = filepath.Join(BASEDIR, "var")
	DBDIR = filepath.Join(STATEDIR, "db")
	MEDIA = filepath.Join(STATEDIR, "media")
	TMP = filepath.Join(STATEDIR, "tmp")
	RUNDIR = filepath.Join(STATEDIR, "run")
	WORKSPACE = filepath.Join(STATEDIR, "workspace")
	CONFIG_FILE = filepath.Join(CONFDIR, CONFIG_FILE_NAME)

	// Currently Only supported for Windows
	SetNighthawkDirectory()

	// Load default configuration at startup
	// configuration options will be updated once
	// appropriate configuration file is used
	nhrconfig.LoadDefaultConfig()

	// Setting default Index names for HashSet and StackDb
	//hsconfig.Index = "hashset"
	hsconfig.LoadDefaultConfig()

	//sdconfig.Index = "nighthawk.db"
	sdconfig.LoadDefaultConfig()
	//fmt.Println("Initializing nighthawk configuration")
}

// Create supporting directory structure
// for Windows and Linux to run in standalone mode
func SetNighthawkDirectory() {
	if !fileExists(BASEDIR) {
		os.MkdirAll(BASEDIR, 0755)
	}

	if !fileExists(CONFDIR) {
		os.MkdirAll(CONFDIR, 0755)
	}

	if !fileExists(STATEDIR) {
		os.MkdirAll(STATEDIR, 0755)
	}

	if !fileExists(WORKSPACE) {
		os.MkdirAll(WORKSPACE, 0755)
	}

	if !fileExists(TMP) {
		os.MkdirAll(TMP, 0755)
	}

	if !fileExists(DBDIR) {
		os.MkdirAll(DBDIR, 0755)
	}

	if !fileExists(RUNDIR) {
		os.MkdirAll(RUNDIR, 0755)
	}

	if !fileExists(CONFIG_FILE) {
		var nhr NHRConfig
		nhr.LoadDefaultConfig()
		nhr.SaveConfigFile(CONFIG_FILE)
	}
}

// Local fileExists copy to avoid import cycle
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}

// Statically defined rabbitmq filename and configuration directory
func isRabbitMQConfigured() bool {
	fh, err := os.Open(filepath.Join(CONFDIR, "rabbitmq.json"))
	if err != nil {
		return true
	}
	defer fh.Close()
	return false
}

/******************************************************
*** External API
******************************************************/

func SetConfigFile(configfile string) {
	CONFIG_FILE = configfile
}

func LoadNighthawkConfig() error {
	tnhrconfig, err := LoadConfigFile(CONFIG_FILE)
	if err != nil {
		return err
	}

	// Setting global NHRconfig
	nhrconfig = tnhrconfig
	//hsconfig.LookupEnabled = CheckHashSet()
	//sdconfig.LookupEnabled = CheckStack()

	if CheckHashSet() {
		fmt.Println("LoadNighthawkConfig - Initializing HastSet")
		hsconfig.Initialize()
	}
	if CheckStack() {
		fmt.Println("LoadNighthawkConfig - Initialzing StackDb")
		sdconfig.Initialize()
	}

	// default return
	return nil
}
