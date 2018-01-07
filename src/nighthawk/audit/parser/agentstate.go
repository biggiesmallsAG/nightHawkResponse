/* nighthawk.audit.parser.agentstate.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Parser for HX agentstateinspection
 */

package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	nhc "nighthawk/common"
	nhconfig "nighthawk/config"
	"nighthawk/elastic"
	nhlog "nighthawk/log"
	nhs "nighthawk/nhstruct"
)

const (
	AGENT_EVENT_ADDR_NOTIF   = "addressNotificationEvent"
	AGENT_EVENT_IPV4_NETWORK = "ipv4NetworkEvent"
	AGENT_EVENT_IMAGE_LOAD   = "imageLoadEvent"
	AGENT_EVENT_DNS_LOOKUP   = "dnsLookupEvent"
	AGENT_EVENT_FILE_WRITE   = "fileWriteEvent"
	AGENT_EVENT_REG_KEY      = "regKeyEvent"
)

func ParseAgentStateInspection(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	if !nhconfig.ParserConfigSetting("agentstate") {
		nhlog.LogMessage("ParseAgentStateInspection", "WARNING", "Parsing agentstateinspection not enabled")
		return
	}

	MAX_RECORD := nhconfig.BulkPostSize()

	xmlFile, err := os.Open(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseAgentStateInspection", "ERROR", fmt.Sprintf("Failed to read audit file. %s", err.Error()))
		os.Exit(nhc.ERROR_READING_AUDIT_FILE)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	count := 0 // counter for xmlItem
	total := 0 // counter for total xmlItem

	var inElement string
	var esrecords []nhs.RlRecord

	for {
		// Upload records to Elasticsearch
		if count == MAX_RECORD {
			cmsg := fmt.Sprintf("Uploading %s:%s %d stage agent events to elasticsearch", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParseAgentStateInspection", "DEBUG", cmsg)

			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0
		}

		t, _ := decoder.Token()

		// Upload data to Elasticsearch if not further xml record found
		if t == nil {
			cmsg := fmt.Sprintf("Uploading %s:%s %d stage agent events to elasticsearch", caseinfo.ComputerName, auditinfo.Generator, count+1)
			nhlog.LogMessage("ParseAgentStateInspection", "DEBUG", cmsg)

			elastic.ProcessOutput(caseinfo, auditinfo, esrecords)
			esrecords = esrecords[:0]
			count = 0

			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "eventItem" {
				var item nhs.AgentEventItem
				decoder.DecodeElement(&item, &se)

				var rlrec nhs.RlRecord
				rlrec.ComputerName = caseinfo.ComputerName
				rlrec.CaseInfo = caseinfo
				rlrec.AuditType = auditinfo

				// Switching for rlrec.Record
				switch item.EventType {

				case AGENT_EVENT_ADDR_NOTIF:
					var agentitem nhs.AddressNotificationEvent

					agentitem.TlnTime = item.Timestamp
					agentitem.Timestamp = item.Timestamp
					agentitem.EventType = item.EventType

					for _, eventitem := range item.Details {
						switch eventitem.Name {
						case "address":
							agentitem.Address = eventitem.Value
						}
					}

					rlrec.Record = agentitem

				case AGENT_EVENT_IPV4_NETWORK:
					var agentitem nhs.NetworkEvent

					agentitem.TlnTime = item.Timestamp
					agentitem.Timestamp = item.Timestamp
					agentitem.EventType = item.EventType

					for _, eventitem := range item.Details {
						switch eventitem.Name {
						case "remoteIP":
							agentitem.RemoteIP = eventitem.Value
						case "remotePort":
							numRemotePort, _ := strconv.Atoi(eventitem.Value)
							agentitem.RemotePort = numRemotePort
						case "localIP":
							agentitem.LocalIP = eventitem.Value
						case "localPort":
							numLocalPort, _ := strconv.Atoi(eventitem.Value)
							agentitem.LocalPort = numLocalPort
						case "protocol":
							agentitem.Protocol = eventitem.Value
						case "pid":
							numProcessID, _ := strconv.Atoi(eventitem.Value)
							agentitem.ProcessID = numProcessID
						case "process":
							agentitem.ProcessName = eventitem.Value
						}
					}
					rlrec.Record = agentitem

				case AGENT_EVENT_DNS_LOOKUP:
					var agentitem nhs.DnsLookupEvent

					agentitem.TlnTime = item.Timestamp
					agentitem.Timestamp = item.Timestamp
					agentitem.EventType = item.EventType

					for _, eventitem := range item.Details {
						switch eventitem.Name {
						case "hostname":
							agentitem.Hostname = eventitem.Value
						case "pid":
							numProcessID, _ := strconv.Atoi(eventitem.Value)
							agentitem.ProcessID = numProcessID
						case "process":
							agentitem.ProcessName = eventitem.Value
						}
					}

					rlrec.Record = agentitem

				case AGENT_EVENT_IMAGE_LOAD:
					var agentitem nhs.ImageLoadEvent

					agentitem.TlnTime = item.Timestamp
					agentitem.Timestamp = item.Timestamp
					agentitem.EventType = item.EventType

					for _, eventitem := range item.Details {
						switch eventitem.Name {
						case "fullPath":
							agentitem.FullPath = eventitem.Value
						case "filePath":
							agentitem.FilePath = eventitem.Value
						case "drive":
							agentitem.Drive = eventitem.Value
						case "fileExtension":
							agentitem.FileExt = eventitem.Value
						case "pid":
							numProcessID, _ := strconv.Atoi(eventitem.Value)
							agentitem.ProcessID = numProcessID
						case "process":
							agentitem.ProcessName = eventitem.Value
						}
					}
					rlrec.Record = agentitem

				case AGENT_EVENT_FILE_WRITE:
					var agentitem nhs.FileWriteEvent

					agentitem.TlnTime = item.Timestamp
					agentitem.Timestamp = item.Timestamp
					agentitem.EventType = item.EventType

					for _, eventitem := range item.Details {
						switch eventitem.Name {
						case "fullPath":
							agentitem.FullPath = eventitem.Value
						case "filePath":
							agentitem.FilePath = eventitem.Value
						case "drive":
							agentitem.Drive = eventitem.Value
						case "fileName":
							agentitem.FileName = eventitem.Value
						case "fileExtension":
							agentitem.FileExt = eventitem.Value
						case "devicePath":
							agentitem.DevicePath = eventitem.Value
						case "pid":
							numProcessID, _ := strconv.Atoi(eventitem.Value)
							agentitem.ProcessID = numProcessID
						case "process":
							agentitem.ProcessName = eventitem.Value
						case "writes":
							numWrites, _ := strconv.Atoi(eventitem.Value)
							agentitem.WriteCount = numWrites
						case "numBytesSeenWritten":
							numBytesWritten, _ := strconv.Atoi(eventitem.Value)
							agentitem.BytesWritten = numBytesWritten
						case "lowestFileOffsetSeen":
							numFileOffsetSeen, _ := strconv.Atoi(eventitem.Value)
							agentitem.DataOffset = numFileOffsetSeen
						case "dataAtLowerOffset":
							agentitem.Data = eventitem.Value
						case "textAtLowerOffset":
							agentitem.TextData = eventitem.Value
						case "closed":
							agentitem.IsClosed = true
							if strings.ToLower(eventitem.Value) == "false" {
								agentitem.IsClosed = false
							}
						case "size":
							numSize, _ := strconv.Atoi(eventitem.Value)
							agentitem.FileSize = numSize
						case "md5":
							agentitem.MD5 = eventitem.Value
						}
					}

					rlrec.Record = agentitem

				case AGENT_EVENT_REG_KEY:
					var agentitem nhs.RegKeyEvent

					agentitem.TlnTime = item.Timestamp
					agentitem.Timestamp = item.Timestamp
					agentitem.EventType = item.EventType

					for _, eventitem := range item.Details {
						switch eventitem.Name {
						case "hive":
							agentitem.Hive = eventitem.Value
						case "keyPath":
							agentitem.KeyPath = eventitem.Value
						case "path":
							agentitem.Path = eventitem.Value
						case "eventType":
							numNotifType, _ := strconv.Atoi(eventitem.Value)
							agentitem.NotificationType = numNotifType
						case "pid":
							numProcessID, _ := strconv.Atoi(eventitem.Value)
							agentitem.ProcessID = numProcessID
						case "process":
							agentitem.ProcessName = eventitem.Value
						case "valueName":
							agentitem.ValueName = eventitem.Value
						case "valueType":
							agentitem.ValueType = eventitem.Value
						case "value":
							agentitem.Value = eventitem.Value
						case "text":
							agentitem.Text = eventitem.Value
						}
					}

					rlrec.Record = agentitem

				} // __end_switch_item.EventType

				esrecords = append(esrecords, rlrec)
			}
			count++
			total++
		} // __end_of_switch_se__
	}

	cmsg := fmt.Sprintf("Total number of AgentStateInpsection events %d", total+1)
	nhlog.LogMessage("ParseAgentStateInspection", "INFO", cmsg)

} // __end_of_ParseAgentStateInspection__
