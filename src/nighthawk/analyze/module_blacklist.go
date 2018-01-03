package analyze

import (
	"fmt"
	"strings"
	"context"
	"encoding/json"

	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
	nhlog "nighthawk/log"
	nhat "nighthawk/audit/audittype"

	elastic "gopkg.in/olivere/elastic.v5"
)

func QueryBlacklistInformation(bl *nhs.BlacklistItem) bool {
	
	method := nhconfig.ElasticHttpScheme()
	conn_str := elastic.SetURL(fmt.Sprintf("%s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))
	
	nhlog.LogMessage("QueryBlacklistInformation", "DEBUG", fmt.Sprintf("Setting new elasticsearch client to %s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))

	client, err := elastic.NewClient(conn_str, elastic.SetSniff(false))
	if err != nil {
		nhlog.LogMessage("QueryBlacklistInformation", "ERROR", fmt.Sprintf("Failed to connected to server %s", err.Error()))
		return false
	}

	var query elastic.Query


	switch strings.ToLower(bl.AuditType) {
	case nhat.RL_SERVICES:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", bl.AuditType),
				elastic.NewTermQuery("name.keyword", bl.Name),
				elastic.NewTermQuery("path.keyword", bl.Path),
				elastic.NewTermQuery("service_descriptive_name.keyword", bl.ServiceDescriptiveName))

	case nhat.RL_TASKS:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", bl.AuditType),
				elastic.NewTermQuery("name.keyword", bl.Name),
				elastic.NewTermQuery("task_creator.keyword", bl.TaskCreator))

	case nhat.RL_PERSISTENCE:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", bl.AuditType),
				elastic.NewTermQuery("name", bl.Name),
				elastic.NewTermQuery("persistence_type", bl.PersistenceType),
				elastic.NewTermQuery("path.keyword", bl.Path),
				elastic.NewTermQuery("reg_path.keyword",bl.RegPath))

	case "w32processes":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", "w32processes"),
				elastic.NewTermQuery("name.keyword", bl.Name),
				elastic.NewTermQuery("path.keyword", bl.Path))

	}

	// Show elastcisearch query in debug 
	querySource,_ := query.Source()
	jsonQuery,_ := json.Marshal(querySource)
	nhlog.LogMessage("QueryBlacklistInformatin","DEBUG", fmt.Sprintf("Blacklist query :%s", string(jsonQuery)))
	
	sr, err := client.Search().Index(StackDB).Type("blacklist").Query(query).Do(context.Background())
	if err != nil {
		nhlog.LogMessage("QueryBlacklistInformation", "ERROR", fmt.Sprintf(err.Error()))
		return false
	}
	
	if sr.TotalHits() < 1 {
		return false
	}

	fmt.Println("Blacklist Positive Hitcount: ", sr.TotalHits())
	return true
}