/* nighthawkresponse.audit.parser.filedlhistory.go
 * author: 0xredskull
 *
 * Redline FileDownloadHistory xml to ES
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

func ParseFileDownloadHistory(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseFileDownloadHistory", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

			if inElement == "FileDownloadHistoryItem" {
				var item nhs.FileDownloadHistoryItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.LastModifiedDate == "" {
					item.LastModifiedDate = nhutil.FixEmptyTimestamp()
				}
				if item.LastCheckedDate == "" {
					item.LastCheckedDate = nhutil.FixEmptyTimestamp()
				}

				item.IsWhitelisted = false

				// Add timeline time
				item.TlnTime = item.LastModifiedDate
				host, domain := UrlToHostname(item.SourceUrl)
				item.UrlHostname = host
				item.UrlDomain = domain

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

	cmsg := fmt.Sprintf("Total FileDownloadHistoryItem %d", total)
	nhlog.LogMessage("ParseFileDownloadHistory", "INFO", cmsg)
}
