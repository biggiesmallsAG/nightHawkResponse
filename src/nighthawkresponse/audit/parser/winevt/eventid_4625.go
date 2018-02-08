package winevt

import (
	"fmt"
	"strconv"
	"strings"
)

// https://docs.microsoft.com/en-us/windows/device-security/auditing/event-4625
type EventID4625 struct {
	Version                   int    // Message Version
	SubjectUserSid            string `json:"SubjectSecurityID"`
	SubjectUserName           string `json:"SubjectAccountName"`
	SubjectDomainName         string `json:"SubjectAccountDomain"`
	SubjectLogonId            string `json:"SubjectLogonID"`
	LogonType                 string `json:"LogonType"`
	TargetUserSid             string `json:"LogonSecurityID"`
	TargetUserName            string `json:"LogonAccountName"`
	TargetDomainName          string `json:"LogonAccountDomain"`
	FailureReason             string `json:"FailureReason"`
	Status                    string `json:"Status"`
	SubStatus                 string `json:"SubStatus"`
	ProcessId                 string `json:"CallerProcessID"`
	ProcessName               string `json:"CallerProcessName"`
	WorkstationName           string `json:"WorkstationName"`
	IpAddress                 string `json:"SourceNetworkAddress"`
	IpPort                    int    `json:"SourcePort"`
	LogonProcessName          string `json:"LogonProcess"`
	AuthenticationPackageName string `json:"AuthenticationPackage"`
	TransitedServices         string `json:"TransitedServices"`
	LmPackageName             string `json:"PackageName"`
	KeyLength                 int    `json:"KeyLength"`
}

func (evtlog *EventID4625) ParseEventMessage(message string) {
	var logbuffer []string
	lg := strings.Split(message, "\n")

	for _, line := range lg {
		matched := strings.Contains(line, ":")
		if matched {
			logbuffer = append(logbuffer, line)
		}
	}

	var logdata []string
	for _, line := range logbuffer {
		m := strings.SplitN(line, ":", 2)
		logdata = append(logdata, strings.TrimSpace(m[1]))
	}

	logdata_size := len(logdata)
	evtlog.Version = getEventID4625Version(message, logdata_size)

	// If evtlog.Version is not 0 then error out
	if evtlog.Version == -1 {
		fmt.Println("EventID2625 - Invalid version - IndexSize does not match")
		return
	}

	switch evtlog.Version {
	case 0:
		// Subject: 0
		evtlog.SubjectUserSid = logdata[1]
		evtlog.SubjectUserName = logdata[2]
		evtlog.SubjectDomainName = logdata[3]
		evtlog.SubjectLogonId = logdata[4]
		//logon_type, _ := strconv.Atoi(logdata[5])
		//evtlog.LogonType = logon_type
		evtlog.LogonType = logdata[5]
		// Account For Which Logon Failed: 6
		evtlog.TargetUserSid = logdata[7]
		evtlog.TargetUserName = logdata[8]
		evtlog.TargetDomainName = logdata[9]
		// Failure Information: 10
		evtlog.FailureReason = logdata[11]
		evtlog.Status = logdata[12]
		evtlog.SubStatus = logdata[13]
		// Process Information: 14
		evtlog.ProcessId = logdata[15]
		evtlog.ProcessName = logdata[16]
		// Network Information: 17
		evtlog.WorkstationName = logdata[18]
		evtlog.IpAddress = logdata[19]
		_sport, _ := strconv.Atoi(logdata[20])
		evtlog.IpPort = _sport
		// Detailed Authentication Information: 21
		evtlog.LogonProcessName = logdata[22]
		evtlog.AuthenticationPackageName = logdata[23]
		evtlog.TransitedServices = logdata[24]
		evtlog.LmPackageName = logdata[25]
		_keylength, _ := strconv.Atoi(logdata[26])
		evtlog.KeyLength = _keylength
	case 1:
		// not in use - placeholder for futuer implementation
	case 2:
		// not in use - placeholder for future implementation
	} //__end_switch_
}

func getEventID4625Version(message string, logdata_size int) int {
	// This helps to cover non-English characters
	// in message body
	if logdata_size == 27 {
		return 0 // Windows 2008 compatible. This message structure hasn't changed
	}
	// Default returns error
	return -1
}

func (evtlog *EventID4625) MessageString() string {
	var event string = ""
	event += fmt.Sprintln("Subject SecurityID: ", evtlog.SubjectUserSid)
	event += fmt.Sprintln("Subject Account Name: ", evtlog.SubjectUserName)
	event += fmt.Sprintln("Subject Account Domain: ", evtlog.SubjectDomainName)
	event += fmt.Sprintln("Subject LogonID: ", evtlog.SubjectLogonId)
	event += fmt.Sprintln("Logon Type: ", evtlog.LogonType)
	event += fmt.Sprintln("Logon SecurityID: ", evtlog.TargetUserSid)
	event += fmt.Sprintln("Logon Account Name: ", evtlog.TargetUserName)
	event += fmt.Sprintln("Logon Account Domain: ", evtlog.TargetDomainName)
	event += fmt.Sprintln("Failure Reason: ", evtlog.FailureReason)
	event += fmt.Sprintln("Failure Status: ", evtlog.Status)
	event += fmt.Sprintln("Failure Sub-Status: ", evtlog.SubStatus)
	event += fmt.Sprintln("Caller ProcessID: ", evtlog.ProcessId)
	event += fmt.Sprintln("Caller ProcessName: ", evtlog.ProcessName)
	event += fmt.Sprintln("Workstation Name: ", evtlog.WorkstationName)
	event += fmt.Sprintln("Source Network Address: ", evtlog.IpAddress)
	event += fmt.Sprintln("Source Port: ", evtlog.IpPort)
	event += fmt.Sprintln("Logon Process: ", evtlog.LogonProcessName)
	event += fmt.Sprintln("Authenticated Package: ", evtlog.AuthenticationPackageName)
	event += fmt.Sprintln("Transited Services: ", evtlog.TransitedServices)
	event += fmt.Sprintln("Package Name: ", evtlog.LmPackageName)
	event += fmt.Sprintln("Key Length: ", evtlog.KeyLength)
	return event
}
