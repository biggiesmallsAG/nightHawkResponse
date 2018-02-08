/* nighthawkresponse.audit.parser.urlhistory.go
 * author: 0xredskull
 *
 * Redline URLHistory XML file parser
 */

package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"

	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"
	"nighthawkresponse/elastic"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
	nhutil "nighthawkresponse/util"
)

func ParseUrlHistory(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {

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

			if inElement == "UrlHistoryItem" {
				var item nhs.UrlHistoryItem
				decoder.DecodeElement(&item, &se)

				// Fix empty timestamp
				if item.LastVisitDate == "" {
					item.LastVisitDate = nhutil.FixEmptyTimestamp()
				}
				item.TlnTime = item.LastVisitDate

				item.IsWhitelisted = false

				//Extracting domain and subdomain for stacking
				h, d := UrlToHostname(item.Url)
				item.UrlHostname = h
				item.UrlDomain = d

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

	cmsg := fmt.Sprintf("Total UrlHistoryItem %d", total)
	nhlog.LogMessage("ParseUrlHistory", "INFO", cmsg)
}

// Extract Hostname and Domain from url
func UrlToHostname(Url string) (string, string) {
	var Hostname string = ""
	var Domain string = ""

	re, _ := regexp.Compile("(http|https|ftp)://[^/]+")

	if !re.MatchString(Url) {
		return "", ""
	}

	r0, _ := regexp.Compile(":\\d+")
	baseUrl := r0.ReplaceAllString(Url, "")

	baseUrl = re.FindString(baseUrl)

	prefixList := []string{"http", "https", "ftp"}

	for _, urlProto := range prefixList {
		if strings.HasPrefix(baseUrl, urlProto) {
			baseUrl = strings.Replace(baseUrl, urlProto+"://", "", 2)
		}
	}

	ipre, _ := regexp.Compile("(\\d+\\.){3}\\d+")
	if ipre.MatchString(baseUrl) {
		return baseUrl, baseUrl
	}

	urlPart := strings.Split(baseUrl, ".")
	numPart := len(urlPart)

	if numPart <= 2 {
		return baseUrl, baseUrl
	} else if numPart == 4 {
		Hostname = baseUrl
		Domain = urlPart[1] + "." + urlPart[2] + "." + urlPart[3]
	} else {
		if len(urlPart[numPart-1]) == 2 {
			Hostname = baseUrl
			Domain = urlPart[numPart-3] + "." + urlPart[numPart-2] + "." + urlPart[numPart-1]
		} else {
			Hostname = baseUrl
			Domain = urlPart[numPart-2] + "." + urlPart[numPart-1]
		}
	}

	return Hostname, Domain
}
