/* nighthawk.audit.parser.arp.go
 * author: 0xredskull
 *
 * Parser for NetworkArp
*/

package parser 

import (
 	"fmt"
 	"os"
 	"encoding/xml"

 	"nighthawk/elastic"
 	nhconfig "nighthawk/config"
 	nhs "nighthawk/nhstruct"
 	nhc "nighthawk/common"
 	nhlog "nighthawk/log"
 	nhutil "nighthawk/util"
)


func ParseNetworkArp(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile,err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseNetworkArp", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

			if inElement == "ArpEntryItem" {
				var item nhs.ArpEntryItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.JobCreated == "" {item.JobCreated = nhutil.FixEmptyTimestamp()}

				item.IsWhitelisted = false	
				// Add timeline time
				// item.TlnTime = item.CreationTime

				//jitem,_ := json.Marshal(item)

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

	cmsg := fmt.Sprintf("Total ArpEntryItem %d", total)
	nhlog.LogMessage("ParseNetworkArp", "INFO", cmsg)
}



