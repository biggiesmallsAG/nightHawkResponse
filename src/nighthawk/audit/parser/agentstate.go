/* nighthawk.audit.parser.agentstate.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for HX agentstateinspection
 */

package parser

import (
 	"fmt"
 	"os"
 	"encoding/xml"

 	nhconfig "nighthawk/config"
 	nhs "nighthawk/nhstruct"
 	nhlog "nighthawk/log"
 	nhc "nighthawk/common"
)


func ParseAgentStateInspection(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	if !nhconfig.ParserConfigSetting("agentstate") {
		nhlog.LogMessage("ParseAgentStateInspection", "WARNING", "Parsing agentstateinspection not enabled")
		return
	} 

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseAgentStateInspection", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
		os.Exit(nhc.ERROR_READING_AUDIT_FILE)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	count := 0
	total := 0

	var inElement string 

	for {
		// Upload records to Elasticsearch 
		if count == MAX_RECORD  {
			cmsg := fmt.Sprintf("Uploading %s:%s %d stage agent events to elasticsearch", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParseAgentStateInspection", "INFO", cmsg)
			count = 0
		}

		t,_ := decoder.Token()

		// Upload data to Elasticsearch if not further xml record found
		if t == nil {
			cmsg := fmt.Sprintf("Uploading %s:%s %d stage agent events to elasticsearch", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParseAgentStateInspection", "INFO", cmsg)
			count = 0

			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "eventItem" {
				count++
				total++
			}
		}
	}

	cmsg := fmt.Sprintf("Total number of AgentStateInpsection events %d", total+1)
	nhlog.LogMessage("ParseAgentStateInspection", "INFO", cmsg)

}