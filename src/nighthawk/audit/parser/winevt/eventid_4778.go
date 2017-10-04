package winevt

import (
	"fmt"
	"strings"
)

// https://docs.microsoft.com/en-us/windows/device-security/auditing/event-4778
type EventID4778 struct {
	Version       int
	AccountName   string `json:"AccountName"`
	AccountDomain string `json:"AccountDomain"`
	LogonID       string `json:"LogonID"`
	SessionName   string `json:"SessionName"`
	ClientName    string `json:"ClientName"`
	ClientAddress string `json:"ClientAddress"`
}

func (evtlog *EventID4778) ParseEventMessage(message string) {
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
	evtlog.Version = getEventID4778Version(message, logdata_size)

	if evtlog.Version == -1 {
		fmt.Println("EventID4778 - invalid version - index size does not match")
		return
	}

	switch evtlog.Version {
	case 0:
		// Subject: 0
		evtlog.AccountName = logdata[1]
		evtlog.AccountDomain = logdata[2]
		evtlog.LogonID = logdata[3]
		// Session: 4
		evtlog.SessionName = logdata[5]
		// Additional Information: 6
		evtlog.ClientName = logdata[7]
		evtlog.ClientAddress = logdata[8]
	}
}

func getEventID4778Version(message string, logdata_size int) int {
	if logdata_size == 9 {
		return 0
	}
	// Default return error
	return -1
}
