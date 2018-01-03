/* nighthawk.elastic.elastic
 * author: 0xredskull && biggiesmalls
 * Team nightHawk.
 *
 * Elasticsearch interface to Upload triage data
 */

package elastic

import (
 	"os"
	"context"
	"encoding/json"
	"fmt"
	"time"

	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
	nhc "nighthawk/common"
	nhlog "nighthawk/log"
	

	elastic "gopkg.in/olivere/elastic.v5"
)

type DateCreated struct {
	DC time.Time `json:"date_created"`
}


func ExportToElasticsearch(computername string, auditgenerator string, data []nhs.RlRecord) {
	var method string 		// Elasticsearch Http Method
	var con_str elastic.ClientOptionFunc

	method = nhconfig.ElasticHttpScheme()
	con_str = elastic.SetURL(fmt.Sprintf("%s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))

	nhlog.LogMessage("ExportToElasticsearch", "DEBUG", fmt.Sprintf("Elasticsearch Server %s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))
	
	client, err := elastic.NewClient(con_str, elastic.SetSniff(false))
	if err != nil {
		nhlog.LogMessage("ExportToElasticsearch", "ERROR", fmt.Sprintf("Failed to connect to server. %s", err.Error()))
		os.Exit(nhc.ERROR_ELASTIC_CONNECT)
	}

	// Check elasticsearch TypeID as computername
	parentExists, err := client.Get().
		Index(nhconfig.ElasticIndex()).
		Type("hostname").
		Id(computername).
		Do(context.Background())

	// Query for ComputerName TypeID failed
	// err contains "Error 404 (Not Found)"
	if err != nil {
		nhlog.LogMessage("ExportToElasticsearch", "INFO", fmt.Sprintf("Parent not found for %s, creating parent..", computername))

		// Creating Parent TypeID
		parent, err := client.Index().
			Index(nhconfig.ElasticIndex()).
			Type("hostname").
			Id(computername).
			BodyJson(DateCreated{
				DC: time.Now().UTC(),
			}).
			Do(context.Background())

		if err != nil {
			nhlog.LogMessage("ExportToElasticsearch", "ERROR", fmt.Sprintf("Failed creating parent object %s", computername))
			nhlog.LogMessage("ExportToElasticsearch", "ERROR", err.Error())
			os.Exit(nhc.ERROR_ELASTIC_CREATE_PARENT)
		}

		nhlog.LogMessage("ExportToElasticsearch", "INFO", fmt.Sprintf("Parent created ID %s", parent.Id))
		uploadBulkData(client, computername, data)
		return
	}

	if parentExists.Found {
		uploadBulkData(client, computername, data)
	}
}


func uploadBulkData(client *elastic.Client, computername string, data []nhs.RlRecord) {
	bulk := client.Bulk()
	for _, record := range data {
		d, _ := json.Marshal(record)

		req := elastic.NewBulkIndexRequest().
			Index(nhconfig.ElasticIndex()).
			Type("audit_type").
			Parent(computername).
			Doc(string(d))

		bulk.Add(req)
	}

	res, err := bulk.Do(context.Background())
	if err != nil {
		nhlog.LogMessage("uploadBulkData", "ERROR", fmt.Sprintf("Failed to bulk index. %s", err.Error()))
		return 
	}

	// If no error encountered uploading data
	i := res.Indexed()[0]
	if i.Status != 201 {
		nhlog.LogMessage("uploadBulkData","ERROR", fmt.Sprintf("ES Returned Status: %d, %s", i.Status, i.Error.Reason))
	} else {
		nhlog.LogMessage("uploadBulkData", "DEBUG", fmt.Sprintf("Upload successful. Server responded with status: %d", i.Status))
	}
}