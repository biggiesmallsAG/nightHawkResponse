/* nighthawk.audit.parser.tasks.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for Windows Scheduled tasks 
 */

package parser

import (
 	"fmt"
 	"os"
 	"encoding/xml"

 	"nighthawk/elastic"
 	nhconfig "nighthawk/config"
 	nhs "nighthawk/nhstruct"
 	nhutil "nighthawk/util"
 	nhlog "nighthawk/log"
 	nhc "nighthawk/common"
 	"nighthawk/analyze"
)


func ParseTasks(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile,err := os.Open(auditfile) 
	if err != nil {
		nhlog.LogMessage("ParseTasks", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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
			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0
		}

		t,_ := decoder.Token()
		if t == nil {
			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0

			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local

			if inElement == "TaskItem" {
				var item nhs.TaskItem 
				decoder.DecodeElement(&item, &se)

				// Fixing empty timestamp
				if item.CreationDate == "" { item.CreationDate = nhutil.FixEmptyTimestamp() }
				if item.MostRecentRunTime == "" {item.MostRecentRunTime = nhutil.FixEmptyTimestamp()}
				if item.NextRunTime == "" {item.NextRunTime = nhutil.FixEmptyTimestamp()}


				// Update the Path and Arguments information from ActionList
				// If there are more than one action item then use the first Path from
				// ActionList
				for i,_ := range item.ActionList {
					if item.Path == "" && item.ActionList[i].Path != "" {
						item.Path = item.ActionList[i].Path
						item.Arguments = item.ActionList[i].ExecArguments
					}
				}

				for i,_ := range item.TriggerList {
					if item.TriggerList[i].TriggerBegin == "" {item.TriggerList[i].TriggerBegin = nhutil.FixEmptyTimestamp()}
				}


				item.IsWhitelisted = false

				// Checking StackDatabase 
				// _rm> 2017-06-07
				verifiedStatus,verifiedVerdict := analyze.TaskIsVerified(item)
				if verifiedStatus {
					item.IsGoodTask = "true"
					item.NhComment.Date = nhutil.CurrentTimestamp()
					item.NhComment.Analyst = "nighthawk"
					item.NhComment.Comment = verifiedVerdict
				}
								
				// Set TlnTime
				item.TlnTime = item.CreationDate

				var rlrec nhs.RlRecord
				rlrec.ComputerName = caseinfo.ComputerName
				rlrec.CaseInfo = caseinfo 
				rlrec.AuditType = auditinfo 
				rlrec.Record = item 

				esrecords = append(esrecords, rlrec)
				count++
				total++
			}
		}
	}

	cmsg := fmt.Sprintf("Total TaskItem %d", total)
	nhlog.LogMessage("ParseTasks", "INFO", cmsg)

}