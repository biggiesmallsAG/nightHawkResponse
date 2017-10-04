/* nighthawk.audit.parser.process_memory.go
 * author: roshan maskey <roshanmaskey>
 *
 * Parser for Windows Memory
 */

package parser

import (
 	"fmt"
 	"os"
 	"encoding/xml"

 	nhconfig "nighthawk/config"
 	nhs "nighthawk/nhstruct"
 	nhutil "nighthawk/util"
 	nhlog "nighthawk/log"
 	nhc "nighthawk/common"
 	"nighthawk/elastic"
 	"nighthawk/analyze"	

)


func ParseProcessMemory(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile,err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseProcessMemory", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
		os.Exit(nhc.ERROR_READING_AUDIT_FILE)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	count := 0
	total := 0

	var inElement string
	var esrecords []nhs.RlRecord
	var processitems []nhs.ProcessItem 

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
			if inElement == "ProcessItem" {
				var item nhs.ProcessItem 
				decoder.DecodeElement(&item, &se)

				// Fix timestamp
				if item.StartTime == "" {item.StartTime = nhutil.FixEmptyTimestamp()}
				if item.KernelTime == "" {item.KernelTime = nhutil.FixEmptyTimestamp()}
				if item.UserTime == "" {item.UserTime = nhutil.FixEmptyTimestamp()}

				item.IsWhitelisted = false

				// Checking process properties in Stack Database
				// _rm> 2017-06-07
				verifiedStatus, verifiedVerdict := analyze.ProcessIsVerified(item)
				if verifiedStatus {
					item.IsGoodProcess = "true"
					item.NhComment.Date = nhutil.CurrentTimestamp()
					item.NhComment.Analyst = "nighthawk"
					item.NhComment.Comment = verifiedVerdict 
				}

				// Set timeline time
				item.TlnTime = item.StartTime

				var rlrec nhs.RlRecord
				rlrec.ComputerName = caseinfo.ComputerName
				rlrec.CaseInfo = caseinfo 
				rlrec.AuditType = auditinfo
				rlrec.Record = item 

				esrecords = append(esrecords, rlrec)
				processitems = append(processitems, item)
				count++
				total++
			}
		}
	}

	cmsg := fmt.Sprintf("Total ProcessItem %d", total)
	nhlog.LogMessage("ParseProcessMemory", "INFO", cmsg)


	// Create process tree
	CreateProcessTree(caseinfo, processitems)
}
