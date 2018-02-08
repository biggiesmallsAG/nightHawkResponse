/* nighthawkresponse.audit.parser.files.go
 * author: 0xredskull
 *
 * Parser for RAWFiles and APIFiles
 */

package parser

import (
	"encoding/xml"
	"fmt"
	"nighthawkresponse/analyze"
	"os"

	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"
	"nighthawkresponse/elastic"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
	nhutil "nighthawkresponse/util"
)

func ParseRawFiles(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseRawFiles", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

			if inElement == "FileItem" {
				var item nhs.RawFileItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.JobCreated == "" {
					item.JobCreated = nhutil.FixEmptyTimestamp()
				}
				if item.Created == "" {
					item.Created = nhutil.FixEmptyTimestamp()
				}
				if item.Modified == "" {
					item.Modified = nhutil.FixEmptyTimestamp()
				}
				if item.Accessed == "" {
					item.Accessed = nhutil.FixEmptyTimestamp()
				}
				if item.Changed == "" {
					item.Changed = nhutil.FixEmptyTimestamp()
				}
				if item.FilenameCreated == "" {
					item.FilenameCreated = nhutil.FixEmptyTimestamp()
				}
				if item.FilenameModified == "" {
					item.FilenameModified = nhutil.FixEmptyTimestamp()
				}
				if item.FilenameAccessed == "" {
					item.FilenameAccessed = nhutil.FixEmptyTimestamp()
				}
				if item.FilenameChanged == "" {
					item.FilenameChanged = nhutil.FixEmptyTimestamp()
				}

				item.IsWhitelisted = analyze.RawFileIsWhitelisted(&item)
				item.IsBlacklisted = analyze.RawFileIsBlacklisted(&item)

				// Add timeline time
				item.TlnTime = item.FilenameCreated

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

	cmsg := fmt.Sprintf("Total RawFiles %d", total)
	nhlog.LogMessage("ParseRawFiles", "INFO", cmsg)
}

// Parsing Windows APIFiles

func ParseApiFiles(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseApiFiles", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
		os.Exit(nhc.ERROR_READING_AUDIT_FILE)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	count := 0
	total := 0

	var inElement string
	var esrecords []nhs.RlRecord

	for {
		if count == MAX_RECORD-1 {
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

			if inElement == "FileItem" {
				var item nhs.FileItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.JobCreated == "" {
					item.JobCreated = nhutil.FixEmptyTimestamp()
				}
				if item.Created == "" {
					item.Created = nhutil.FixEmptyTimestamp()
				}
				if item.Modified == "" {
					item.Modified = nhutil.FixEmptyTimestamp()
				}
				if item.Accessed == "" {
					item.Accessed = nhutil.FixEmptyTimestamp()
				}
				if item.Changed == "" {
					item.Changed = nhutil.FixEmptyTimestamp()
				}

				item.IsBlacklisted = analyze.FileIsBlacklisted(&item)
				item.IsWhitelisted = analyze.FileIsWhitelisted(&item)

				// Add timeline time
				item.TlnTime = item.Created

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

	cmsg := fmt.Sprintf("Total ApiFiles %d", total)
	nhlog.LogMessage("ParseApiFiles", "INFO", cmsg)
}
