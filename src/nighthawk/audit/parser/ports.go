/* nighthawk.audit.parser.port.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * w32port parser
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
 	"nighthawk/elastic"
)


func ParsePorts(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

		MAX_RECORD := nhconfig.BulkPostSize()

		xmlFile, err := os.Open(auditfile) 
		if err != nil {
			nhlog.LogMessage("ParsePorts", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

				if inElement == "PortItem" {
					var item nhs.PortItem
					decoder.DecodeElement(&item, &se)

					// Fix empty timestamp
					// PortItem has no attribute with timestamp
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
			}
		}
		cmsg := fmt.Sprintf("Total PortItem processed %d", total)
		nhlog.LogMessage("ParsePorts", "INFO", cmsg)
}