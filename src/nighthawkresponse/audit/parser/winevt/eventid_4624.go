package winevt

import (
	"fmt"
	"strconv"
	"strings"
)

// https://docs.microsoft.com/en-us/windows/device-security/auditing/event-4624
type EventID4624 struct {
	Version                   int
	SubjectUserSid            string `json:"SubjectSecurityID"`
	SubjectUserName           string `json:"SubjectAccountName"`
	SubjectDomainName         string `json:"SubjectAccountDomain"`
	SubjectLogonID            string `json:"SubjectLogonID"`
	LogonType                 string `json:"LogonType"`
	RestrictedAdminMode       string `json:"RestrictedAdminMode"`
	VirtualAccountNo          string `json:"VirtualAccountNo"`
	ElevatedToken             string `json:"ElevatedToken"`
	ImpersonationLevel        string `json:"ImpersonationLevel"`
	TargetUserSid             string `json:"LogonSecurityID"`
	TargetUserName            string `json:"LogonAccountName"`
	TargetDomainName          string `json:"LogonAccountDomain"`
	TargetLogonId             string `json:"LogonID"`
	TargetLinkedLogonId       string `json:"LinkedLogonID"`
	TargetOutboundUserName    string `json:"NetworkAccountName"`
	TargetOutboundDomainName  string `json:"NetworkAccountDomain"`
	LogonGuid                 string `json:"LogonGUID"`
	ProcessId                 string `json:"ProcessID"`
	ProcessName               string `json:"ProcessName"`
	WorkstationName           string `json:"WorkstationName"`
	IpAddress                 string `json:"SourceNetworkAddress"`
	IpPort                    int    `json:"SourcePort"`
	LogonProcessName          string `json:"LogonProcess"`
	AuthenticationPackageName string `json:"AuthenticationPackage"`
	TransitedServices         string `json:"TransitedServices"`
	LmPackageName             string `json:"PackageName"`
	KeyLength                 int    `json:"KeyLength"`
}

func (evtlog *EventID4624) ParseEventMessage(message string) {
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
	evtlog.Version = getEventID4624Version(message, datasize)

	// Return if version number check returns error
	if evtlog.Version == -1 {
		fmt.Printf("EventID4624 - Invalid version - %d index size does not match\n", datasize)
		return
	}

	switch evtlog.Version {
	case 2:
		// Subject
		evtlog.SubjectUserSid = data[1]
		evtlog.SubjectUserName = data[2]
		evtlog.SubjectDomainName = data[3]
		evtlog.SubjectLogonID = data[4]
		// Logon Information
		//_logon_type, _ := strconv.Atoi(data[6])
		//evtlog.LogonType = _logon_type
		evtlog.LogonType = data[6]
		evtlog.RestrictedAdminMode = data[7]
		evtlog.VirtualAccountNo = data[8]
		evtlog.ElevatedToken = data[9]
		evtlog.ImpersonationLevel = data[10]
		// New Logon
		evtlog.TargetUserSid = data[12]
		evtlog.TargetUserName = data[13]
		evtlog.TargetDomainName = data[14]
		evtlog.TargetLogonId = data[15]
		evtlog.TargetLinkedLogonId = data[16]
		evtlog.TargetOutboundUserName = data[17]
		evtlog.TargetOutboundDomainName = data[18]
		evtlog.LogonGuid = data[19]
		// Process Information
		evtlog.ProcessId = data[21]
		evtlog.ProcessName = data[22]
		// Network Information
		evtlog.WorkstationName = data[24]
		evtlog.IpAddress = data[25]
		_source_port, _ := strconv.Atoi(data[26])
		evtlog.IpPort = _source_port
		// Detailed Authentication Information
		evtlog.LogonProcessName = data[28]
		evtlog.AuthenticationPackageName = data[29]
		evtlog.TransitedServices = data[30]
		evtlog.LmPackageName = data[31]
		_key_length, _ := strconv.Atoi(data[32])
		evtlog.KeyLength = _key_length

	case 1:
		// Subject
		evtlog.SubjectUserSid = data[1]
		evtlog.SubjectUserName = data[2]
		evtlog.SubjectDomainName = data[3]
		evtlog.SubjectLogonID = data[4]
		//_logon_type, _ := strconv.Atoi(data[5])
		//evtlog.LogonType = _logon_type
		evtlog.LogonType = data[5]
		evtlog.ImpersonationLevel = data[6]
		// New Logon
		evtlog.TargetUserSid = data[8]
		evtlog.TargetUserName = data[9]
		evtlog.TargetDomainName = data[10]
		evtlog.TargetLogonId = data[11]
		evtlog.LogonGuid = data[12]
		// Process Information
		evtlog.ProcessId = data[14]
		evtlog.ProcessName = data[15]
		// Network Information
		evtlog.WorkstationName = data[17]
		evtlog.IpAddress = data[18]
		_source_port, _ := strconv.Atoi(data[19])
		evtlog.IpPort = _source_port
		// Detailed Authentication Information
		evtlog.LogonProcessName = data[21]
		evtlog.AuthenticationPackageName = data[22]
		evtlog.TransitedServices = data[23]
		evtlog.LmPackageName = data[24]
		_key_length, _ := strconv.Atoi(data[25])
		evtlog.KeyLength = _key_length

	case 0:
		// Subject: 0
		evtlog.SubjectUserSid = data[1]
		evtlog.SubjectUserName = data[2]
		evtlog.SubjectDomainName = data[3]
		evtlog.SubjectLogonID = data[4]
		//_logon_type, _ := strconv.Atoi(data[5])
		//evtlog.LogonType = _logon_type
		evtlog.LogonType = data[5]
		//New Logon
		evtlog.TargetUserSid = data[7]
		evtlog.TargetUserName = data[8]
		evtlog.TargetDomainName = data[9]
		evtlog.TargetLogonId = data[10]
		evtlog.LogonGuid = data[11]
		// Process Informa tion
		evtlog.ProcessId = data[13]
		evtlog.ProcessName = data[14]
		// Network Information
		evtlog.WorkstationName = data[16]
		evtlog.IpAddress = data[17]
		_source_port, _ := strconv.Atoi(data[18])
		evtlog.IpPort = _source_port
		// Detailed Authentication Information
		evtlog.LogonProcessName = data[20]
		evtlog.AuthenticationPackageName = data[21]
		evtlog.TransitedServices = data[22]
		evtlog.LmPackageName = data[23]
		_key_length, _ := strconv.Atoi(data[24])
		evtlog.KeyLength = _key_length
	} //__end_switch_

}

func getEventID4624Version(message string, datasize int) int {

	if strings.Contains(message, "Linked Logon ID:") && datasize == 33 {
		return 2 // Windows 2016 compatible logon event
	}

	if strings.Contains(message, "Impersonation Level:") && datasize == 26 {
		return 1 // Windows 2012 compatible logon event
	}

	// This covers non-english charset
	// Note: This is purely counting ":" in message detail
	// may result in unexpected mapping in some cases
	switch datasize {
	case 33:
		return 2 // Guessed windows 2016 logon event
	case 26:
		return 1 // Guessed windows 2012 logon event
	case 25:
		return 0 // Guessed windows 2008 logon event
	}

	// Default error response
	return -1
}

func (evtlog *EventID4624) MessageString() string {
	var event string = ""
	event += fmt.Sprintln("Subject SecurityID: ", evtlog.SubjectUserSid)
	event += fmt.Sprintln("Subject Account Name: ", evtlog.SubjectUserName)
	event += fmt.Sprintln("Subject Account Domain: ", evtlog.SubjectDomainName)
	event += fmt.Sprintln("Subject LogonID: ", evtlog.SubjectLogonID)
	event += fmt.Sprintln("Logon Type: ", evtlog.LogonType)
	event += fmt.Sprintln("Logon Restricted Admin Mode: ", evtlog.RestrictedAdminMode)
	event += fmt.Sprintln("Virtual Account No: ", evtlog.VirtualAccountNo)
	event += fmt.Sprintln("Elevated Token: ", evtlog.ElevatedToken)
	event += fmt.Sprintln("Impersonation Level: ", evtlog.ImpersonationLevel)
	event += fmt.Sprintln("Logon SecurityID: ", evtlog.TargetUserSid)
	event += fmt.Sprintln("Logon Account Name: ", evtlog.TargetUserName)
	event += fmt.Sprintln("Logon Account Domain: ", evtlog.TargetDomainName)
	event += fmt.Sprintln("LogonID: ", evtlog.TargetLogonId)
	event += fmt.Sprintln("Linked LogonID: ", evtlog.TargetLinkedLogonId)
	event += fmt.Sprintln("Logon Network Account Name: ", evtlog.TargetOutboundUserName)
	event += fmt.Sprintln("Logon Network Account Domain: ", evtlog.TargetOutboundDomainName)
	event += fmt.Sprintln("Logon GUID: ", evtlog.LogonGuid)
	event += fmt.Sprintln("ProcessID: ", evtlog.ProcessId)
	event += fmt.Sprintln("ProcessName: ", evtlog.ProcessName)
	event += fmt.Sprintln("ProcessName: ", evtlog.ProcessName)
	event += fmt.Sprintln("Workstation Name: ", evtlog.WorkstationName)
	event += fmt.Sprintln("Source Network Address: ", evtlog.IpAddress)
	event += fmt.Sprintln("Source Port: ", evtlog.IpPort)
	event += fmt.Sprintln("Logon Process: ", evtlog.LogonProcessName)
	event += fmt.Sprintln("Authentication Package: ", evtlog.AuthenticationPackageName)
	event += fmt.Sprintln("Transited Services: ", evtlog.TransitedServices)
	event += fmt.Sprintln("Package Name: ", evtlog.LmPackageName)
	event += fmt.Sprintln("Key Length: ", evtlog.KeyLength)

	return event
}
