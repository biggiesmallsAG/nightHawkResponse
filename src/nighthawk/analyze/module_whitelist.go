package analyze

import (
	"fmt"
	"strings"
	"context"

	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
	nhlog "nighthawk/log"
	nhat "nighthawk/audit/audittype"

	elastic "gopkg.in/olivere/elastic.v5"
)

func QueryWhitelistInformation(wl *nhs.WhitelistItem) bool {
	
	method := nhconfig.ElasticHttpScheme()
	conn_str := elastic.SetURL(fmt.Sprintf("%s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))
	
	nhlog.LogMessage("QueryWhitelistInformation", "DEBUG", fmt.Sprintf("Setting new elasticsearch client to %s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))

	client, err := elastic.NewClient(conn_str, elastic.SetSniff(false))
	if err != nil {
		nhlog.LogMessage("QueryWhitelistInformation", "ERROR", fmt.Sprintf("Failed to connected to server %s", err.Error()))
		return false
	}

	var query elastic.Query


	switch strings.ToLower(wl.AuditType) {
	case nhat.RL_SERVICES:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", wl.AuditType),
				elastic.NewTermQuery("name.keyword", wl.Name),
				elastic.NewTermQuery("path.keyword", wl.Path),
				elastic.NewTermQuery("service_descriptive_name.keyword", wl.ServiceDescriptiveName))

	case nhat.RL_TASKS:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", wl.AuditType),
				elastic.NewTermQuery("name.keyword", wl.Name),
				elastic.NewTermQuery("task_creator.keyword", wl.TaskCreator))

	case nhat.RL_PERSISTENCE:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", wl.AuditType),
				elastic.NewTermQuery("name", wl.Name),
				elastic.NewTermQuery("persistence_type", wl.PersistenceType),
				elastic.NewTermQuery("path.keyword", wl.Path),
				elastic.NewTermQuery("reg_path.keyword",wl.RegPath))

	case "w32processes":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", "w32processes"),
				elastic.NewTermQuery("name.keyword", wl.Name),
				elastic.NewTermQuery("path.keyword", wl.Path))
	}


	sr, err := client.Search().Index(StackDB).Type("whitelist").Query(query).Do(context.Background())
	if err != nil {
		nhlog.LogMessage("QueryWhitelistInformation", "ERROR", fmt.Sprintf(err.Error()))
		return false
	}
	
	if sr.TotalHits() < 1 {
		return false
	}

	return true
}