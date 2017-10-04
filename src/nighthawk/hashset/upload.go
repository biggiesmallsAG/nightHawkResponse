package hashset 


import (
	"fmt"
	"context"
	"encoding/json"
	elastic "gopkg.in/olivere/elastic.v5"
)


func UploadWhitelistHashSet(data []HashSet) {
	bulk := client.Bulk()

	for _, record := range data {
		json_record,_ := json.Marshal(record)

		req := elastic.NewBulkIndexRequest().
				Index(INDEX_NAME).
				Type(WL_HASH_INDEX).
				Doc(string(json_record))
		bulk.Add(req)
	}

	res, err := bulk.Do(context.Background())
	if err != nil {
		panic(err)
	}

	i := res.Indexed()[0]
	if i.Status != 201 {
		fmt.Printf("Elastic - ERROR - StatusCode: %d, Reason: %s\n", i.Status, i.Error.Reason)
	} else {
		fmt.Printf("Elastic - INFO - Uploaded %d records\n", len(data))
	}
}
