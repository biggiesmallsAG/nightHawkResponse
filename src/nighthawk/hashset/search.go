package hashset

import (
	"fmt"
	"context"
	"encoding/json"

	elastic "gopkg.in/olivere/elastic.v5"
)


func SearchWhitelistHash(field string, value string, size int, bWriteOutput bool) bool { 
	if value == "" {
		return false
	}
	
	hsResult := QueryElasticsearch(field, value, size, WL_HASH_INDEX)

	if len(hsResult) >= 1 {
		if bWriteOutput {
			for _,h := range hsResult {
				h.WriteConsole()
			}
		}
		return true
	} 
	return false
}


func QueryElasticsearch(field string, value string, size int, index_type string) []HashSet {
	query := elastic.NewMatchQuery(field,value)
	sr,err := client.Search().
				Index(INDEX_NAME).
				Type(index_type).
				Size(size).
				Query(query).
				Do(context.Background())
	if err != nil {
		panic(err)
		return nil 
	}

	// Show elasticsearch query stats if DEBUG flag is 
	// enabled
	if DEBUG {
		fmt.Printf("SearchTime: %d ms\n", sr.TookInMillis)
		fmt.Printf("TotalHits: %d\n", sr.Hits.TotalHits)	
	}

	var hsResult []HashSet

	for _,s := range sr.Hits.Hits {
		var hs HashSet 
		json.Unmarshal(*s.Source, &hs)
		hsResult = append(hsResult, hs)
	}
	return hsResult 
}
