/* nighthawkresponse.audit.parser.systeminfo.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for SystemInformation
 */

package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"
	"nighthawkresponse/elastic"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
	nhutil "nighthawkresponse/util"
)

func ParseSystemInfo(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseSystemInfo", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
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

			if inElement == "SystemInfoItem" {
				var item nhs.SystemInfoItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.BiosDate == "" {
					item.BiosDate = nhutil.FixEmptyTimestamp()
				} else {
					fixedDate := FixBiosDate(item.BiosDate)
					item.BiosDate = fixedDate
				}

				if item.Date == "" {
					item.Date = nhutil.FixEmptyTimestamp()
				}
				if item.InstallDate == "" {
					item.InstallDate = nhutil.FixEmptyTimestamp()
				}
				if item.AppCreated == "" {
					item.AppCreated = nhutil.FixEmptyTimestamp()
				}

				for i := range item.NetworkList {
					if item.NetworkList[i].DhcpLeaseObtained == "" {
						item.NetworkList[i].DhcpLeaseObtained = nhutil.FixEmptyTimestamp()
					}

					if item.NetworkList[i].DhcpLeaseExpires == "" {
						item.NetworkList[i].DhcpLeaseExpires = nhutil.FixEmptyTimestamp()
					}
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

	cmsg := fmt.Sprintf("Total SystemInfoItem %d", total)
	nhlog.LogMessage("ParseSystemInfo", "INFO", cmsg)
}

// BiosDate is represented as MM/DD/YYYY
// Changing the date to ISO format
func FixBiosDate(biosdate string) string {
	s := strings.SplitN(biosdate, "/", 3)
	if len(s) == 3 {
		newBiosDate := fmt.Sprintf("%s-%s-%sT00:00:00Z", s[2], s[1], s[0])
		return newBiosDate
	} else {
		return nhutil.FixEmptyTimestamp()
	}

}
