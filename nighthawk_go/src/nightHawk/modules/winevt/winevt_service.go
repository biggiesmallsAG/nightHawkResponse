/* 
 * Parse Windows Event ID 4697 (Windows 2008): A service was installed in the system
 * Parse Windows Event ID 601 (Windows 2003): Attempt to install service
 * 
 * reference: https://www.ultimatewindowssecurity.com/securitylog/encyclopedia/event.aspx?eventid=601
 * 			  https://www.ultimatewindowssecurity.com/securitylog/encyclopedia/event.aspx?eventID=4697
 * 
 * credit: 	Phil Kealy for suggesting to include it.
 *
 * author:  roshan maskey <0xredskull>
 */

package winevt


import (
	"strings"
)

func (ev *EventServiceInstall)ParseEventMessage(message string) {
	message = strings.Replace(message, "\t", "", -1)
	data := strings.Split(message, "\r\n")
	dataLen := len(data)

	fields := []string{"Subject", "Service Information"}
	indexMap := GetSectionIndex(fields, data)

	// Windows 2008 and Windows 2003 service install message format is different

	if indexMap["Service Information"] > 0 {
		// Reference for Windows 2008 => https://www.ultimatewindowssecurity.com/securitylog/encyclopedia/event.aspx?eventID=4697
		// Parsing Windows 2008 Service Install log
		startIndex := indexMap["Subject"]
		stopIndex := indexMap["Service Information"]

		var acc AccountInfo
		acc.ReadAccountInfo(startIndex, stopIndex, data)
		ev.Subject = acc 

		// Parsing Service Information 
		startIndex = indexMap["Service Information"]
		stopIndex = dataLen 

		for i:= startIndex; i < stopIndex; i++ {
			keyval := strings.SplitN(data[i], ":", 2)
			if len(keyval) == 2 {
				switch keyval[0] {
				case "Service Name":
					ev.ServiceName = strings.TrimSpace(keyval[1])

				case "Service File Name":
					ev.ServiceFileName = strings.TrimSpace(keyval[1])

				case "Service Type":
					ev.ServiceType = strings.TrimSpace(keyval[1])

				case "Service Start Type":
					ev.ServiceStartType = strings.TrimSpace(keyval[1])

				case "Service Account":
					ev.ServiceAccount = strings.TrimSpace(keyval[1])
				}
			}
		}	
	} else {
		// Reference for Windows 2003 => https://www.ultimatewindowssecurity.com/securitylog/encyclopedia/event.aspx?eventid=601
		
		startIndex := 0
		stopIndex := dataLen
		
		var acc AccountInfo 

		for i:= startIndex; i < stopIndex; i++ {
			keyval := strings.SplitN(data[i], ":", 2)
			if len(keyval) == 2 {
				switch keyval[0] {
				case "Service Name":
					ev.ServiceName = strings.TrimSpace(keyval[1])

				case "Service File Name":
					ev.ServiceFileName = strings.TrimSpace(keyval[1])

				case "Service Type":
					ev.ServiceType = strings.TrimSpace(keyval[1])

				case "Service Start Type":
					ev.ServiceStartType = strings.TrimSpace(keyval[1])

				case "User Name":
					acc.AccountName = strings.TrimSpace(keyval[1])

				case "Domain":
					acc.AccountDomain = strings.TrimSpace(keyval[1])

				case "Logon ID":
					acc.LogonId = strings.TrimSpace(keyval[1])
				}
			}
		}

		ev.Subject = acc 
	}
	
}