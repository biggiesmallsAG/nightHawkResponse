package config

import (
	"fmt"
	"net/http"
)

type HashSetConfig struct {
	Index           string
	LookupEnabled   bool
	LookupAvailable bool
	LookupChecked   bool
}

func (config *HashSetConfig) LoadDefaultConfig() {
	config.Index = "hashset"
	config.LookupEnabled = false
	config.LookupAvailable = false
	config.LookupChecked = false
}

// checking hashset
func (config *HashSetConfig) Initialize() {
	hashsetUrl := fmt.Sprintf("%s://%s:%d/%s", ElasticHttpScheme(), ElasticHost(), ElasticPort(), config.Index)
	res, err := http.Get(hashsetUrl)
	if err != nil {
		config.LookupChecked = true // HASHSET_LOOKUP_CHECKED = true
		fmt.Println("HashSet::Initialize - Failed to connect to HashSet - ", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		fmt.Println("HashSet::Initialize - OK")
		config.LookupEnabled = true
		config.LookupAvailable = true // HASHSET_LOOKUP_AVAILABLE = true
	} else {
		fmt.Println("HashSet::initialize - Status code ", res.StatusCode)
	}
	config.LookupChecked = true // HASHSET_LOOKUP_CHECKED = true
}

/******************************
*** External Functions
*******************************/

func HashSetIndex() string {
	return hsconfig.Index
}

func HashSetEnabled() bool {
	return hsconfig.LookupEnabled
}

func HashSetAvailable() bool {
	return hsconfig.LookupAvailable
}

func HashSetChecked() bool {
	return hsconfig.LookupChecked
}
