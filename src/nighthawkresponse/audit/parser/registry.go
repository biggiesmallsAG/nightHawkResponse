/* nighthawkresponse.audit.parser.registry.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for Windows registry
 */

package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"
	"nighthawkresponse/elastic"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
	nhutil "nighthawkresponse/util"
)

func ParseRegistry(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseRegistry", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

		t, _ := decoder.Token()
		if t == nil {
			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0

			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local

			if inElement == "RegistryItem" {
				var item nhs.RegistryItem
				decoder.DecodeElement(&item, &se)

				// Fix timestamp
				if item.Modified == "" {
					item.Modified = nhutil.FixEmptyTimestamp()
				}
				// Set timestamp
				item.TlnTime = item.Modified
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

	cmsg := fmt.Sprintf("Total RegistryItem %d", total)
	nhlog.LogMessage("ParseRegistry", "INFO", cmsg)
}
