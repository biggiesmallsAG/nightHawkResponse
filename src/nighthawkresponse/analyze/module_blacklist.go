package analyze

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	nhat "nighthawkresponse/audit/audittype"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"

	elastic "gopkg.in/olivere/elastic.v5"
)

func QueryBlacklistInformation(bl *nhs.BlacklistItem) bool {

	switch strings.ToLower(bl.AuditType) {
	case nhat.RL_SERVICES:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", bl.AuditType),
				elastic.NewTermQuery("Name.keyword", bl.Name),
				elastic.NewTermQuery("Path.keyword", bl.Path),
				elastic.NewTermQuery("ServiceDescriptiveName.keyword", bl.ServiceDescriptiveName))

	case nhat.RL_TASKS:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", bl.AuditType),
				elastic.NewTermQuery("Name.keyword", bl.Name),
				elastic.NewTermQuery("TaskCreator.keyword", bl.TaskCreator))

	case nhat.RL_PERSISTENCE:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", bl.AuditType),
				elastic.NewTermQuery("Name.keyword", bl.Name),
				elastic.NewTermQuery("PersistenceType.keyword", bl.PersistenceType),
				elastic.NewTermQuery("Path.keyword", bl.Path),
				elastic.NewTermQuery("RegPath.keyword", bl.RegPath))

	case "w32processes":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", "w32processes"),
				elastic.NewTermQuery("Name.keyword", bl.Name),
				elastic.NewTermQuery("Path.keyword", bl.Path))

	case nhat.RL_APIFILES, nhat.RL_RAWFILES, "w32files":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", "w32files"),
				elastic.NewTermQuery("Name.keyword", bl.Name),
				elastic.NewTermQuery("Md5sum.keyword", bl.Md5sum))
	default:
		// Unhandled AuditType return false
		fmt.Printf("Unhandled audit_type %s\n", bl.AuditType)
		return false

	}

	// Show elastcisearch query in debug
	querySource, _ := query.Source()
	jsonQuery, _ := json.Marshal(querySource)
	nhlog.LogMessage("QueryBlacklistInformatin", "DEBUG", fmt.Sprintf("Blacklist query :%s", string(jsonQuery)))

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
