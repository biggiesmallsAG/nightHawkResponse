package analyze

import (
	"context"
	"fmt"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"

	nhat "nighthawkresponse/audit/audittype"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
)

const StackDB = "nighthawk"
const CertTable = "cert"
const StackTable = "stack"

func QueryCAInformation(cainfo *nhs.IssuingCA) bool {

	query := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("common_name", cainfo.CommonName))
	sr, err := client.Search().Index(StackDB).Type(CertTable).Query(query).Do(context.Background())
	if err != nil {
		nhlog.LogMessage("QueryCAInformation", "ERROR", fmt.Sprintf(err.Error()))
		return false
	}

	if sr.TotalHits() < 1 {
		return false
	}

	return true
}

func IsKnownCertIssuer(certissuer string) bool {
	if certissuer == "" {
		return false
	}

	var ca nhs.IssuingCA
	ca.CommonName = certissuer
	return QueryCAInformation(&ca)
}

func IsCommonStackItem(si *nhs.StackItem) bool {

	switch strings.ToLower(si.AuditType) {
	case nhat.RL_SERVICES:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", si.AuditType),
				elastic.NewTermQuery("Name.keyword", si.Name),
				elastic.NewTermQuery("Path.keyword", si.Path),
				elastic.NewTermQuery("ServiceDescriptiveName.keyword", si.ServiceDescriptiveName))

	case nhat.RL_TASKS:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", si.AuditType),
				elastic.NewTermQuery("Name.keyword", si.Name),
				elastic.NewTermQuery("TaskCreator.keyword", si.TaskCreator))

	case nhat.RL_PERSISTENCE:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", si.AuditType),
				elastic.NewTermQuery("Name.keyword", si.Name),
				elastic.NewTermQuery("PersistenceType.keyword", si.PersistenceType),
				elastic.NewTermQuery("Path.keyword", si.Path),
				elastic.NewTermQuery("RegPath.keyword", si.RegPath))

	case "w32processes":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", "w32processes"),
				elastic.NewTermQuery("Name.keyword", si.Name),
				elastic.NewTermQuery("Path.keyword", si.Path))

	case nhat.RL_APIFILES, nhat.RL_RAWFILES, "w32files":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("AuditType.keyword", "w32files"),
				elastic.NewTermQuery("Name.keyword", si.Name),
				elastic.NewTermQuery("Md5sum.keyword", si.Md5sum))
	default:
		// Unhandled StackItem audit_type
		fmt.Printf("Unhandled audit_type for StackItem: %s\n", si.AuditType)
		return false

	}

	sr, err := client.Search().Index(StackDB).Type(StackTable).Query(query).Do(context.Background())
	if err != nil {
		nhlog.LogMessage("IsCommonStackItem", "ERROR", fmt.Sprintf(err.Error()))
		return false
	}

	if sr.TotalHits() < 1 {
		return false
	}

	return true
}
