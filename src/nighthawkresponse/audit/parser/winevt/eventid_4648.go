package winevt

import (
	"fmt"
	"strconv"
	"strings"
)

// https://docs.microsoft.com/en-us/windows/device-security/auditing/event-4648
type EventID4648 struct {
	Version           int    // MessageVersion
	SubjectUserSid    string `json:"SubjectSecurityID"`
	SubjectUserName   string `json:"SubjectAccountName"`
	SubjectDomainName string `json:"SubjectAccountName"`
	SubjectLogonId    string `json:"SubjectLogonID"`
	LogonGuid         string `json:"SubjectLogonGUID"`
	TargetUserName    string `json:"LogonAccountName"`
	TargetDomainName  string `json:"LogonAccountDomain"`
	TargetLogonGuid   string `json:"LogonGUID"`
	TargetServerName  string `json:"TargetServerName"`
	TargetInfo        string `json:"AdditionalInformation"`
	ProcessId         string `json:"ProcessID"`
	ProcessName       string `json:"ProcessName"`
	IpAddress         string `json:"NetworkAddress"`
	IpPort            int    `json:"Port"`
}

func (evtlog *EventID4648) ParseEventMessage(message string) {
	var data []string
	lg := strings.Split(message, "\n")

	for _, line := range lg {
		matched := strings.Contains(line, ":")
		if matched {
			m := strings.SplitN(line, ":", 2)
			data = append(data, strings.TrimSpace(m[1]))
		}
	}

	datasize := len(data)
	evtlog.Version = getEventID4648Version(message, datasize)

	if evtlog.Version == -1 {
		fmt.Println("EventID4648 - Invalid version - index size mismatch")
		return
	}

	switch evtlog.Version {
	case 0:
		// Subject: 0
		evtlog.SubjectUserSid = data[1]
		evtlog.SubjectUserName = data[2]
		evtlog.SubjectDomainName = data[3]
		evtlog.SubjectLogonId = data[4]
		evtlog.LogonGuid = data[5]
		// Account Whose Credential Were Used: 6
		evtlog.TargetUserName = data[7]
		evtlog.TargetDomainName = data[8]
		evtlog.TargetLogonGuid = data[9]
		// Target Server: 10
		evtlog.TargetServerName = data[11]
		evtlog.TargetInfo = data[12]
		// Process Information: 13
		evtlog.ProcessId = data[14]
		evtlog.ProcessName = data[15]
		// Network Information: 16
		evtlog.IpAddress = data[17]
		ip_port, _ := strconv.Atoi(data[18])
		evtlog.IpPort = ip_port
	case 1:
		// not in use - placeholder for future implementation
	case 2:
		// not in use - placeholder for future implementation
	} // __end_switch_

}

func getEventID4648Version(message string, datasize int) int {

	if datasize == 19 {
		return 0 // Windows 2008 comptabile format
	}
	// Default returns error
	return -1
}
