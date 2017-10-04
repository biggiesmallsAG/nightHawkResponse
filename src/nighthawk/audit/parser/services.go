/* nighthawk.audit.parser.services.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for Windows Services
 */

package parser 

import (
 	"fmt"
 	"os"
 	"encoding/xml"

 	nhconfig "nighthawk/config"
 	nhs "nighthawk/nhstruct"
 	"nighthawk/elastic"
 	"nighthawk/analyze"
 	nhutil "nighthawk/util"
 	nhlog "nighthawk/log"
 	nhc "nighthawk/common"
)



func ParseServices(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile,err := os.Open(auditfile) 
	if err != nil {
		nhlog.LogMessage("ParseServices", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
		os.Exit(nhc.ERROR_READING_AUDIT_FILE)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	count := 0 // counter for XMLItem
	total := 0 // counter for total xmlItem

	var inElement string 
	var esrecords []nhs.RlRecord

	for {
		// Check for upload condition
		if count == MAX_RECORD {
			cmsg := fmt.Sprintf("Uploading %s:%s %d service items", caseinfo.ComputerName, auditinfo.Generator, count +1)
			nhlog.LogMessage("ParseServices", "DEBUG", cmsg)

			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0
		}	

		t,_ := decoder.Token()
		if t == nil {
			cmsg := fmt.Sprintf("Uploading %s:%s %d service items", caseinfo.ComputerName, auditinfo.Generator, count +1)
			nhlog.LogMessage("ParseServices", "DEBUG", cmsg)

			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "ServiceItem" {
				var item nhs.ServiceItem 
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				// no-timestamp data
				item.IsWhitelisted = false

				// Check ServiceItem in CommonStackDatabase
				// _rm> 2017-06-07
				verifiedStatus,verifiedVerdict := analyze.ServiceIsVerified(item)
				if verifiedStatus {
					item.IsGoodService = "true"
					item.NhComment.Date = nhutil.CurrentTimestamp()
					item.NhComment.Analyst = "nighthawk"
					item.NhComment.Comment = verifiedVerdict
				}

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

	cmsg := fmt.Sprintf("Total service event processed %d", total)
	nhlog.LogMessage("ParseServices", "INFO", cmsg)

}