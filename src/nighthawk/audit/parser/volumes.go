/* nighthawk.audit.parser.volumes.go
 * author: 0xredskull
 *
 * Windows Volume information parser
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
)


func ParseVolumes(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile,err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseVolumes", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

			if inElement == "VolumeItem" {
				var item nhs.VolumeItem 
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.CreationTime == "" {item.CreationTime = nhutil.FixEmptyTimestamp()}
				item.IsWhitelisted = false
				
				// Add timeline time
				//item.TlnTime = item.CreationTime

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

	cmsg := fmt.Sprintf("Total VolumeItem %d", total)
	nhlog.LogMessage("ParseVolumes", "INFO", cmsg)
}



