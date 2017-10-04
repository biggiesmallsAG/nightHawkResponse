package hashset

import (
	"fmt"
	"strconv"
	"strings"

	nhconfig "nighthawk/config"
	nhlog "nighthawk/log"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	UPLOAD_SIZE 	= 2000			// Bulk upload size for Elasticsearch
	INDEX_NAME		= "hashset"		
	WL_HASH_INDEX	= "whitelist"	// Whitelisted MD5 hashes
	BL_HASH_INDEX	= "blacklist"	// Blacklisted MD5 hashes

	// Elasticsearch Server Variables
	elastic_server = "127.0.0.1"
	elastic_port = 9200
	elastic_user = ""
	elastic_pass = ""

	client *elastic.Client = nil		// Elasticsearch client

	DEBUG = false
)



func LoadHashSetConfig(){

	// Try loading configuration from nighthawk.json 
	elastic_server = nhconfig.ElasticHost()
	elastic_port = nhconfig.ElasticPort()
	elastic_user = nhconfig.ElasticUser()
	elastic_pass = nhconfig.ElasticPassword()
	

	client = SetElasticClient(elastic_server, elastic_port, elastic_user, elastic_pass)
}


type HashSet struct {
	Sha1 		string `json:"sha1"`
	Md5 		string `json:"md5"`
	FileName	string `json:"filename"`
	FileSize	int    `json:"filesize"`
}


func (h *HashSet) ReadLine(line string) {
	line = strings.Replace(line, "\"", "", -1)
	data := strings.SplitN(line, ",", -1)

	h.Sha1 = data[0]
	h.Md5 = data[1]
	h.FileName = data[2]
	h.FileSize,_ = strconv.Atoi(data[3])
}


func (h *HashSet) WriteConsole() {
	fmt.Printf("SHA1=%s, MD5=%s, FileName=%s, FileSize=%d\n", h.Sha1, h.Md5, h.FileName, h.FileSize)
}


func SetElasticClient(server string, port int, user string, pass string) *elastic.Client {
	var http_scheme string = "http"
	conn_str := elastic.SetURL(fmt.Sprintf("%s://%s:%d", http_scheme, server, port))
	client, err := elastic.NewClient(conn_str, elastic.SetSniff(false))

	if err != nil {
		nhlog.LogMessage("SetElasticClient", "WARNING", "Elasticsearch server is unavailable")
		return nil
	}

	return client 
}


func UpdateElasticClient(server string, port int, user string, pass string) {
	client = SetElasticClient(server, port, user, pass)
}

func SetDebug(enabled bool) {
	DEBUG = enabled
}