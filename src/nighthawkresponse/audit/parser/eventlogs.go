/* eventlogs.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Primary module to parse event windows event logs
 */

package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"nighthawkresponse/audit/parser/winevt"
	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"
	"nighthawkresponse/elastic"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
	nhutil "nighthawkresponse/util"
)

func ParseEventLogs(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseEventLogs", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
		os.Exit(nhc.ERROR_READING_AUDIT_FILE)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	count := 0
	total := 0
	var inElement string

	var esrecords []nhs.RlRecord

	for {

		if count == MAX_RECORD {
			cmsg := fmt.Sprintf("Uploading %s::%s %d event records", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParseEventLogs", "INFO", cmsg)
			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)

			esrecords = esrecords[:0]
			count = 0
		}

		t, _ := decoder.Token()
		if t == nil {
			cmsg := fmt.Sprintf("Uploading %s::%s %d event records", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParseEventLogs", "INFO", cmsg)
			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)

			esrecords = esrecords[:0]
			count = 0
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "EventLogItem" {
				//nighthawkresponse.ConsoleMessage("INFO", "Parse event item", nhconfig.VERBOSE)
				var item nhs.EventLogItem
				decoder.DecodeElement(&item, &se)

				// Fixing Empty timestamps
				if item.GenTime == "" {
					item.GenTime = nhutil.FixEmptyTimestamp()
				}

				if item.WriteTime == "" {
					item.GenTime = nhutil.FixEmptyTimestamp()
				}

				// Parsing message details
				item.MessageDetail = winevt.ProcessEventItem(item.Log, item.EID, item.Message)
				item.IsWhitelisted = false

				var rlrec nhs.RlRecord
				rlrec.ComputerName = caseinfo.ComputerName
				rlrec.CaseInfo = caseinfo
				rlrec.AuditType = auditinfo
				rlrec.Record = item

				esrecords = append(esrecords, rlrec)

				count++
				total++
			}
		default:
		}
	}

	cmsg := fmt.Sprintf("Number of Event log entries %d", total)
	nhlog.LogMessage("ParseEventLogs", "INFO", cmsg)
}
