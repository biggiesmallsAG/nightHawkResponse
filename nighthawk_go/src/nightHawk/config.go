/*
 *@package 	nightHawk
 *@file 	config.go
 *@author	roshan maskey <roshanmaskey@gmail.com>
 *@version	0.0.1
 *@updated	2016-06-15
 *
 *@description	nightHawk Configuration settings
 */


 package nightHawk


 import (
 	"path/filepath"
 	"io/ioutil"
 	"encoding/json"
 )

 const VERSION = "0.0.1"


 // Directory Settings for nightHawk 
 const BASEDIR = "/opt/nighthawk"
 var WORKSPACE = filepath.Join(BASEDIR, "var/workspace")
 var CONFIG = filepath.Join(BASEDIR, "etc/nighthawk.json")
 var TMP = filepath.Join(BASEDIR, "var/tmp")


 
 // nightHawk Configuration Default Settings
 var VERBOSE = false
 var VERBOSE_LEVEL = 6
 var SESSIONDIR_SIZE = 8

 /// Configuration File Loadable Configuration options
 // nightHawk Configuration
 var MAXPROCS = 2
 var MAX_GOROUTINE = 100
 var BULKPOST_SIZE = 10000

 // Controlling output
 const (
 	OP_CONSOLE_ONLY = 1
 	OP_DATASTORE_ONLY = 2
 	OP_CONSOLE_DATASTORE = 3
 	OP_WRITE_FILE = 4
 )
 var OPCONTROL = OP_DATASTORE_ONLY


/// Elastic Configuration
 var ELASTICHOST = ""
 var ELASTICPORT = 0
 var ELASTIC_INDEX = ""



/// Config Loading structure and functions
type nHConfig struct {
	MaxProcs 			int `json:"max_procs"`
	MaxGorouting 		int `json:"max_goroutine"`
	BulkPostSize 		int `json:"bulk_post_size"`
	OpControl 			int `json:"opcontrol"`
	SessionDirSize 		int `json:"sessiondir_size"` 
	Verbose 			bool `json:"verbose"`
	VerboseLevel 		int `json:"verbose_level"`
}

type nHElastic struct {
	ElasticHost 		string `json:"elastic_server"`
	ElasticPort 		int `json:"elastic_port"`
	ElasticIndex 		string `json:"elastic_index"`
}

type nightHawkConfig struct {
	NightHawk 			nHConfig `json:"nightHawk"`
	Elastic				nHElastic `json:"elastic"`
}


func LoadConfigFile(configfile string) bool {
	configData, err := ioutil.ReadFile(configfile)

	if err != nil {
		return false
	}

	var nhconfig nightHawkConfig 
	json.Unmarshal(configData, &nhconfig)

	MAXPROCS = nhconfig.NightHawk.MaxProcs
	MAX_GOROUTINE = nhconfig.NightHawk.MaxGorouting
	BULKPOST_SIZE = nhconfig.NightHawk.BulkPostSize
	OPCONTROL = nhconfig.NightHawk.OpControl

	if nhconfig.NightHawk.SessionDirSize > SESSIONDIR_SIZE {
		SESSIONDIR_SIZE = nhconfig.NightHawk.SessionDirSize
	}

	VERBOSE = nhconfig.NightHawk.Verbose
	if nhconfig.NightHawk.VerboseLevel > 0 {
		VERBOSE_LEVEL = nhconfig.NightHawk.VerboseLevel
	}


	ELASTICHOST = nhconfig.Elastic.ElasticHost
	ELASTICPORT = nhconfig.Elastic.ElasticPort
	ELASTIC_INDEX = nhconfig.Elastic.ElasticIndex

	return true
}










