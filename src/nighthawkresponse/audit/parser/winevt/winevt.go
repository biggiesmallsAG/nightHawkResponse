/*
 * This module is used by AuditParser parsing Windows Event Logs
 *
 * author: roshan maskey <0xredskull>
 */

package winevt

import (
	"strings"
)

// Windows EventID to be parsed
const (
	EventLogonSuccess     = "4624"
	EventLogonFailed      = "4625"
	EventExplicitCred     = "4648"
	EventNewProcess       = "4688"
	EventServiceInstalled = "4697"
)

type EventNull struct {
	EventDescription string
}

// ProcessEventItem extracts key-value from EventLog Message field.
// It supports EventID 4624, 4625, 4648, 4688, 4697
func ProcessEventItem(logSource string, eventId string, message string) interface{} {

	if strings.ToLower(logSource) == "security" {
		switch eventId {
		case EventLogonSuccess:
			var ev EventID4624
			ev.ParseEventMessage(message)
			return ev

		case EventLogonFailed:
			var ev EventID4625
			ev.ParseEventMessage(message)
			return ev

		case EventExplicitCred:
			var ev EventID4648
			ev.ParseEventMessage(message)
			return ev

		case EventNewProcess:
			var ev EventID4688
			ev.ParseEventMessage(message)
			return ev

		case EventServiceInstalled:
			var ev EventID4697
			ev.ParseEventMessage(message)
			return ev
		}

	}

	var ev EventNull
	return ev
}
