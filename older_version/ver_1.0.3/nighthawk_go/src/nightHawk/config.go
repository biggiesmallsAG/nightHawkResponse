/*
 *@package  nightHawk
 *@file     config.go
 *@author   roshan maskey <roshanmaskey@gmail.com>
 *
 *@description  nightHawk Response Configuration settings
 */


 package nightHawk


 import (
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
 )

 const VERSION = "1.0.3"


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

// Redis Configuration
const REDIS_SERVER = "localhost"
const REDIS_PORT = 6379
const REDIS_CHAN = "broadcast:uploadstatus"
var REDIS_PUB = false

/// Elastic Configuration
 var ELASTICHOST = ""
 var ELASTICPORT = 0
 var ELASTIC_INDEX = ""
 var ELASTIC_SSL = false
 var ELASTIC_AUTHCODE = ""
 var ELASTIC_USER = ""
 var ELASTIC_PASS = ""



/// Config Loading structure and functions
type nHConfig struct {
    MaxProcs            int `json:"max_procs"`
    MaxGorouting        int `json:"max_goroutine"`
    BulkPostSize        int `json:"bulk_post_size"`
    OpControl           int `json:"opcontrol"`
    SessionDirSize      int `json:"sessiondir_size"` 
    Verbose             bool `json:"verbose"`
    VerboseLevel        int `json:"verbose_level"`
}

type nHElastic struct {
    ElasticHost         string `json:"elastic_server"`
    ElasticPort         int `json:"elastic_port"`
    ElasticIndex        string `json:"elastic_index"`
    ElasticSsl          bool `json:"elastic_ssl"`
    ElasticUser         string `json:"elastic_user"`
    ElasticPass         string `json:"elastic_pass"`
    Authcode            string `json:"authcode"`
}

type nightHawkConfig struct {
    NightHawk           nHConfig `json:"nightHawk"`
    Elastic             nHElastic `json:"elastic"`
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
    
    if nhconfig.Elastic.ElasticSsl {
        ELASTIC_SSL = true
    }
    
    // added in ver1.0.3
    ELASTIC_USER = nhconfig.Elastic.ElasticUser
    ELASTIC_PASS = nhconfig.Elastic.ElasticPass 
    //ELASTIC_AUTHCODE = nhconfig.Elastic.Authcode 
    ELASTIC_AUTHCODE = base64.StdEncoding.EncodeToString([]byte(ELASTIC_USER+":"+ELASTIC_PASS))

    return true
}










