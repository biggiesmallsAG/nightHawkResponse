/* nighthawk.audit.parser.prefetch.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for Windows Prefetch 
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
)


func ParsePrefetch(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile,err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParsePrefetch", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
		os.Exit(nhc.ERROR_READING_AUDIT_FILE)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	count := 0
	total := 0

	var inElement string
	var esrecords []nhs.RlRecord

	for {
		// Checking ES post condition
		if count == MAX_RECORD {
			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0
		}

		t,_ := decoder.Token()

		// End of element
		if t == nil {
			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0

			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local

			if inElement == "PrefetchItem" {
				var item nhs.PrefetchItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.Created == "" { item.Created = nhutil.FixEmptyTimestamp() }
				if item.LastRun == "" { item.LastRun = nhutil.FixEmptyTimestamp() }

				item.IsWhitelisted = false

				// Fix timeline time
				// Using Prefetch LastRun time as preferred time
				item.TlnTime = item.LastRun

				//jitem,_ := json.Marshal(item)

				var rlrec nhs.RlRecord
				rlrec.ComputerName = caseinfo.ComputerName
				rlrec.CaseInfo = caseinfo 
				rlrec.AuditType = auditinfo
				//rlrec.Record = string(jitem)
				rlrec.Record = item  

				esrecords = append(esrecords, rlrec)
				count++
				total++
			}
		}
	}
	cmsg := fmt.Sprintf("Total PrefetchItem %d", total)
	nhlog.LogMessage("ParsePrefetch", "INFO", cmsg)

}