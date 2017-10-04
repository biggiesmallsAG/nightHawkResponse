/* nighthawk.audit.parser.useraccounts.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for W32UserAccounts
 */

package parser

import (
 	"fmt"
 	"os"
 	"encoding/xml"

 	
 	nhconfig "nighthawk/config"
 	nhs "nighthawk/nhstruct"
 	"nighthawk/elastic"
 	nhutil "nighthawk/util"
 	nhlog "nighthawk/log"
 	nhc "nighthawk/common"
)


func ParseUserAccounts(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile,err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseUserAccounts", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

			if inElement == "UserItem" {
				var ui nhs.UserItem 
				decoder.DecodeElement(&ui, &se)

				// Fix empty timestamp
				if ui.LastLogin == "" { ui.LastLogin = nhutil.FixEmptyTimestamp() }

				// Timeline timestamp
				ui.TlnTime = ui.LastLogin
				ui.IsWhitelisted = false

				var rlrec nhs.RlRecord
				rlrec.ComputerName = caseinfo.ComputerName
				rlrec.CaseInfo = caseinfo 
				rlrec.AuditType = auditinfo
				rlrec.Record = ui 

				esrecords = append(esrecords, rlrec)
				count++
				total++
			}
		}
	}

	cmsg := fmt.Sprintf("Total UserAccountItem %d", total)
	nhlog.LogMessage("ParseUserAccounts", "INFO", cmsg)
}