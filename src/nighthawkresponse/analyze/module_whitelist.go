package analyze

import (
	"context"
	"fmt"
	"strings"

	nhat "nighthawkresponse/audit/audittype"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"

	elastic "gopkg.in/olivere/elastic.v5"
)

func QueryWhitelistInformation(wl *nhs.WhitelistItem) bool {

	switch strings.ToLower(wl.AuditType) {
	case nhat.RL_SERVICES:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", wl.AuditType),
				elastic.NewTermQuery("Name.keyword", wl.Name),
				elastic.NewTermQuery("Path.keyword", wl.Path),
				elastic.NewTermQuery("ServiceDescriptiveName.keyword", wl.ServiceDescriptiveName))

	case nhat.RL_TASKS:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", wl.AuditType),
				elastic.NewTermQuery("Name.keyword", wl.Name),
				elastic.NewTermQuery("TaskCreator.keyword", wl.TaskCreator))

	case nhat.RL_PERSISTENCE:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", wl.AuditType),
				elastic.NewTermQuery("Name.keyword", wl.Name),
				elastic.NewTermQuery("PersistenceType.keyword", wl.PersistenceType),
				elastic.NewTermQuery("Path.Keyword", wl.Path),
				elastic.NewTermQuery("RegPath.keyword", wl.RegPath))

	case "w32processes":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", "w32processes"),
				elastic.NewTermQuery("Name.keyword", wl.Name),
				elastic.NewTermQuery("Path.keyword", wl.Path))

	case nhat.RL_APIFILES, nhat.RL_RAWFILES, "w32files":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", "w32files"),
				elastic.NewTermQuery("Name.keyword", wl.Name),
				elastic.NewTermQuery("Md5sum.keyword", wl.Md5sum))
	default:
		// Unhandled audit_type for Whitelist
		fmt.Printf("Unhandled audit_type %s for Whitelist\n", wl.AuditType)
		return false
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
