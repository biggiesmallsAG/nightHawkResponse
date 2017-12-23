package config

import (
	"fmt"
	"net/http"
)

type StackDbConfig struct {
	Index           string
	LookupEnabled   bool
	LookupAvailable bool
	LookupChecked   bool
}

func (config *StackDbConfig) LoadDefaultConfig() {
	config.Index = "nighthawk"
	config.LookupEnabled = false
	config.LookupAvailable = false
	config.LookupChecked = false
}

// Initialize function check if elasticsearch index "nighthawk"
// is available to query IssuerCert and stacking information 
func (config *StackDbConfig) Initialize() {
	stackUrl := fmt.Sprintf("%s://%s:%d/%s", ElasticHttpScheme(), ElasticHost(), ElasticPort(), config.Index)
	res, err := http.Get(stackUrl)
	if err != nil {
		config.LookupChecked = true
		fmt.Println("StackDB::initialize - Failed to connect to StackDB - ", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Println("StackDB::initialize - OK")
		config.LookupEnabled = true
		config.LookupAvailable = true
	} else {
		fmt.Println("StackDB::initialize - Status code ", res.StatusCode)
	}
	config.LookupChecked = true
}

func StackDbIndex() string {
	return sdconfig.Index
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
