package analyze

import (
	"strings"
	"fmt"
	"context"

	elastic "gopkg.in/olivere/elastic.v5"

	nhat "nighthawk/audit/audittype"
	nhs "nighthawk/nhstruct"
	nhconfig "nighthawk/config"
	nhlog "nighthawk/log"

)

const StackDB = "nighthawk"
const CertTable = "cert"
const StackTable = "stack"

func QueryCAInformation(cainfo *nhs.IssuingCA) bool {
	method := nhconfig.ElasticHttpScheme()
	conn_str := elastic.SetURL(fmt.Sprintf("%s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))
	
	nhlog.LogMessage("QueryCAInformation", "DEBUG", fmt.Sprintf("Setting new elasticsearch client to %s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))

	client, err := elastic.NewClient(conn_str, elastic.SetSniff(false))
	if err != nil {
		nhlog.LogMessage("QueryCAInformation", "ERROR", fmt.Sprintf("Failed to connected to server %s", err.Error()))
		return false
	}

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
	method := nhconfig.ElasticHttpScheme()
	conn_str := elastic.SetURL(fmt.Sprintf("%s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))
	
	nhlog.LogMessage("IsCommonStackItem", "DEBUG", fmt.Sprintf("Setting new elasticsearch client to %s://%s:%d", method, nhconfig.ElasticHost(), nhconfig.ElasticPort()))

	client, err := elastic.NewClient(conn_str, elastic.SetSniff(false))
	if err != nil {
		nhlog.LogMessage("IsCommonStackItem", "ERROR", fmt.Sprintf("Failed to connected to server %s", err.Error()))
		return false
	}

	var query elastic.Query


	switch strings.ToLower(si.AuditType) {
	case nhat.RL_SERVICES:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", si.AuditType),
				elastic.NewTermQuery("name.keyword", si.Name),
				elastic.NewTermQuery("path.keyword", si.Path),
				elastic.NewTermQuery("service_descriptive_name.keyword", si.ServiceDescriptiveName))

	case nhat.RL_TASKS:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", si.AuditType),
				elastic.NewTermQuery("name.keyword", si.Name),
				elastic.NewTermQuery("task_creator.keyword", si.TaskCreator))

	case nhat.RL_PERSISTENCE:
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", si.AuditType),
				elastic.NewTermQuery("name", si.Name),
				elastic.NewTermQuery("persistence_type", si.PersistenceType),
				elastic.NewTermQuery("path.keyword", si.Path),
				elastic.NewTermQuery("reg_path.keyword",si.RegPath))

	case "w32processes":
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("audit_type.keyword", "w32processes"),
				elastic.NewTermQuery("name.keyword", si.Name),
				elastic.NewTermQuery("path.keyword", si.Path),
				elastic.NewTermQuery("arguments.keyword", si.Arguments))

	}


	sr, err := client.Search().Index(StackDB).Type(StackTable).Query(query).Do(context.Background())
	if err != nil {
		nhlog.LogMessage("IsCommonStackItem", "ERROR", fmt.Sprintf(err.Error()))
		return false
	}
	
	if sr.TotalHits() < 1 {
		return false
	}

	//fmt.Println(sr.TotalHits(), si)
	return true
}
