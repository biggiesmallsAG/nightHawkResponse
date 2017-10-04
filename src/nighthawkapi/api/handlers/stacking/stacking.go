package stacking

import (
	"encoding/json"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"
)

type StackingConfig struct {
	CaseName     string `json:"case_name"`
	ContextItem  string `json:"context_item"`
	Type         string `json:"type"`
	SearchLimit  int    `json:"search_limit"`
	SortDesc     bool   `json:"sort_desc"`
	IgnoreGood   bool   `json:"ignore_good"`
	SubAggSize   int    `json:"sub_agg_size,omitempty"`
	ShowComputer bool   `json:"show_computer"`
	EndpointName string `json:"endpoint"`
}

func (sc *StackingConfig) Default() {
	sc.CaseName = "*"
	sc.ContextItem = "*"
	sc.Type = "*"
	sc.SearchLimit = 100
	sc.SortDesc = true
	sc.IgnoreGood = true
	sc.SubAggSize = 20
	sc.ShowComputer = false
	sc.EndpointName = ""
}

func (sc *StackingConfig) LoadParams(data []byte) {
	sc.Default()

	var tsc StackingConfig
	json.Unmarshal(data, &tsc)

	if tsc.CaseName != "" {
		sc.CaseName = tsc.CaseName
	}

	if tsc.ContextItem != "" {
		sc.ContextItem = tsc.ContextItem
	}

	if tsc.Type != "" {
		sc.Type = tsc.Type
	}

	if tsc.SearchLimit > 0 {
		sc.SearchLimit = tsc.SearchLimit
	}

	if tsc.SortDesc {
		sc.SortDesc = false
	}

	if !tsc.IgnoreGood {
		sc.IgnoreGood = false
	}

	if tsc.SubAggSize > 0 {
		sc.SubAggSize = tsc.SubAggSize
	}

	if tsc.ShowComputer {
		sc.ShowComputer = true
	}

	if tsc.EndpointName != "" {
		sc.EndpointName = tsc.EndpointName
	}
}

func PrintQueryObject(query elastic.Query) {

	jq, _ := query.Source()
	jquery, _ := json.Marshal(jq)
	fmt.Println(string(jquery))
}

func PrintAggObject(agg elastic.Aggregation) {
	jq, _ := agg.Source()
	jquery, _ := json.Marshal(jq)
	fmt.Println(string(jquery))
}
