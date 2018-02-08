package winevt

import (
	"fmt"
	"strings"
)

// https://docs.microsoft.com/en-us/windows/device-security/auditing/event-4688
type EventID4688 struct {
	Version            int    // Message version Win2008: 0, Win2012: 1, Win2016: 2
	SubjectUserSid     string `json:"CreatorSecurityID"`
	SubjectUserName    string `json:"CreatorAccountName"`
	SubjectDomainName  string `json:"CreatorAccountDomain"`
	SubjectLogonId     string `json:"CreatorLogonID"`
	TargetUserSid      string `json:"TargetSecurityID"`
	TargetUserName     string `json:"TargetAccountName"`
	TargetDomainName   string `json:"TargetAccountDomain"`
	TargetLogonId      string `json:"TargetLogonID"`
	NewProcessId       string `json:"NewProcessID"`
	NewProcessName     string `json:"NewProcessName"`
	TokenElevationType string `json:"TokenElevationType"`
	MandatoryLabel     string `json:"MandatoryLabel"`
	ProcessId          string `json:"CreatorProcessID"`
	ParentProcessName  string `json:"CreatorProcessName"`
	CommandLine        string `json:"ProcessCommandLine"`
}

func (evtlog *EventID4688) ParseEventMessage(message string) {
	var logbuffer []string
	lg := strings.Split(message, "\n")
	for _, line := range lg {
		matched := strings.Contains(line, ":")
		if matched {
			logbuffer = append(logbuffer, line)
		}
	}

	var data []string
	for _, line := range logbuffer {
		m := strings.SplitN(line, ":", 2)
		data = append(data, strings.TrimSpace(m[1]))
	}

	datasize := len(data)
	evtlog.Version = getEventID4688Version(message, datasize)

	// Check if ":" count matches one of the templates
	// Return null if no match is found
	if evtlog.Version == -1 {
		fmt.Println("EventID4688 - IndexSize does not match")
		return
	}

	switch evtlog.Version {
	case 0:
	case 1:
		// Subject: 0
		evtlog.SubjectUserSid = data[1]
		evtlog.SubjectUserName = data[2]
		evtlog.SubjectDomainName = data[3]
		evtlog.SubjectLogonId = data[4]
		// Process Information: 5
		evtlog.NewProcessId = data[6]
		evtlog.NewProcessName = data[7]
		evtlog.TokenElevationType = data[8]
		evtlog.MandatoryLabel = data[9]
		evtlog.ProcessId = data[10]
		evtlog.ParentProcessName = data[11]
		evtlog.CommandLine = data[12]
	case 2:
		// Subject: 0
		evtlog.SubjectUserSid = data[1]
		evtlog.SubjectUserName = data[2]
		evtlog.SubjectDomainName = data[3]
		evtlog.SubjectLogonId = data[4]
		// Target Subject: 5
		evtlog.TargetUserSid = data[6]
		evtlog.TargetUserName = data[7]
		evtlog.TargetDomainName = data[8]
		evtlog.TargetLogonId = data[9]
		// Process Information: 10
		evtlog.NewProcessId = data[11]
		evtlog.NewProcessName = data[12]
		evtlog.TokenElevationType = data[13]
		evtlog.MandatoryLabel = data[14]
		evtlog.ProcessId = data[15]
		evtlog.ParentProcessName = data[16]
		evtlog.CommandLine = data[17]
	}
}

func getEventID4688Version(message string, datasize int) int {
	if strings.Contains(message, "Creator Subject") && datasize == 18 {
		return 2 // Windows 2016
	}

	if strings.Contains(message, "Process Command Line") && datasize == 13 {
		return 1 // Windows 2012. Process Command Line was added in Windows 2012 event log
	}

	// This helps to cover non-English characters in
	// message detail
	switch datasize {
	case 18:
		return 2
	case 13:
		return 1
	}

	return -1 // Indicate IndexSize does not match with template
}
