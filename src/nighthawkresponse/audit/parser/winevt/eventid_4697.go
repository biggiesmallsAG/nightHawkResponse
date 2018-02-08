package winevt

import (
	"fmt"
	"strconv"
	"strings"
)

// https://docs.microsoft.com/en-us/windows/device-security/auditing/event-4697
type EventID4697 struct {
	Version           int
	SubjectUserSid    string `json:"SubjectSecurityID"`
	SubjectUserName   string `json:"SubjectAccountName"`
	SubjectDomainName string `json:"SubjectAccountDomain"`
	SubjectLogonId    string `json:"SubjectLogonID"`
	ServiceName       string `json:"ServiceName"`
	ServiceFileName   string `json:"ServiceFileName"`
	ServiceType       string `json:"ServiceType"`
	ServiceStartType  int    `json:"ServiceStartType"`
	ServiceAccount    string `json:"ServiceAccount"`
}

func (evtlog *EventID4697) ParseEventMessage(message string) {
	var logdata []string
	lg := strings.Split(message, "\n")

	for _, line := range lg {
		matched := strings.Contains(line, ":")
		if matched {
			m := strings.SplitN(line, ":", 2)
			logdata = append(logdata, strings.TrimSpace(m[1]))
		}
	}

	logdata_size := len(logdata)
	evtlog.Version = getEventID4697Version(message, logdata_size)

	if evtlog.Version == -1 {
		fmt.Println("EventID4697 - Invalid version - index size does not match")
		return
	}

	switch evtlog.Version {
	case 0:
		// Subject: 0
		evtlog.SubjectUserSid = logdata[1]
		evtlog.SubjectUserName = logdata[2]
		evtlog.SubjectDomainName = logdata[3]
		evtlog.SubjectLogonId = logdata[4]
		// Service Information: 5
		evtlog.ServiceName = logdata[6]
		evtlog.ServiceFileName = logdata[7]
		evtlog.ServiceType = logdata[8]
		service_start_type, _ := strconv.Atoi(logdata[9])
		evtlog.ServiceStartType = service_start_type
		evtlog.ServiceAccount = logdata[10]
	}
}

func getEventID4697Version(message string, logdata_size int) int {
	if logdata_size == 11 {
		return 0 // Windows 2008 compatible
	}
	// Default return error
	return -1
}
