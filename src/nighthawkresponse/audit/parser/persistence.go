/* audit.parser.persistence.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser to parse w32scripting-persistence
 */

package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"nighthawkresponse/analyze"
	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"
	"nighthawkresponse/elastic"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
	nhutil "nighthawkresponse/util"
)

func ParsePersistence(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParsePersistence", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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
			cmsg := fmt.Sprintf("Uploading %s::%s %d events records", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParsePersistence", "DEBUG", cmsg)

			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0
		}

		t, _ := decoder.Token()
		if t == nil {
			cmsg := fmt.Sprintf("Uploading %s::%s %d events records", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParsePersistence", "DEBUG", cmsg)

			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "PersistenceItem" {
				var item nhs.PersistenceItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamps
				if item.RegModified == "" {
					item.RegModified = nhutil.FixEmptyTimestamp()
				}
				if item.FileModified == "" {
					item.FileModified = nhutil.FixEmptyTimestamp()
				}
				if item.FileCreated == "" {
					item.FileCreated = nhutil.FixEmptyTimestamp()
				}
				if item.FileAccessed == "" {
					item.FileAccessed = nhutil.FixEmptyTimestamp()
				}
				if item.FileChanged == "" {
					item.FileChanged = nhutil.FixEmptyTimestamp()
				}
				if item.File.JobCreated == "" {
					item.File.JobCreated = nhutil.FixEmptyTimestamp()
				}
				if item.File.Created == "" {
					item.File.Created = nhutil.FixEmptyTimestamp()
				}
				if item.File.Modified == "" {
					item.File.Modified = nhutil.FixEmptyTimestamp()
				}
				if item.File.Accessed == "" {
					item.File.Accessed = nhutil.FixEmptyTimestamp()
				}
				if item.File.Changed == "" {
					item.File.Changed = nhutil.FixEmptyTimestamp()
				}
				if item.Registry.JobCreated == "" {
					item.Registry.JobCreated = nhutil.FixEmptyTimestamp()
				}
				if item.Registry.Modified == "" {
					item.Registry.Modified = nhutil.FixEmptyTimestamp()
				}
				if item.File.PeInfo.PETimeStamp == "" {
					item.File.PeInfo.PETimeStamp = nhutil.FixEmptyTimestamp()
				}

				// Fix for timeline search in Elasticsearch
				// These attributes were added for timeline and stacking
				item.TlnTime = item.FileCreated
				item.File.TlnTime = item.File.Created
				item.Registry.TlnTime = item.Registry.Modified
				item.StackPath = item.Registry.KeyPath + item.Registry.ValueName

				item.IsWhitelisted = analyze.PersistenceIsWhitelisted(&item)
				item.IsBlacklisted = analyze.PersistenceIsBlacklisted(&item)

				fStatus, fVerdict := analyze.PersistenceIsVerified(&item)
				if fStatus {
					item.File.IsGoodHash = "true"
					item.File.NhComment.Date = nhutil.CurrentTimestamp()
					item.File.NhComment.Analyst = "nighthawkresponse"
					item.File.NhComment.Comment = fVerdict
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
		default:
		}
	}
}
