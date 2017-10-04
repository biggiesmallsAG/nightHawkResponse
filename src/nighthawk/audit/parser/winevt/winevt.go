/*
 * This module is used by AuditParser parsing Windows Event Logs
 *
 * author: roshan maskey <0xredskull>
 */

package winevt

import (
	"strings"
)

const (
	EVENT_LOGON_SUCCESS     = 4624
	EVENT_LOGON_FAILED      = 4625
	EVENT_EXPLICIT_CRED     = 4648
	EVENT_NEW_PROCESS       = 4688
	EVENT_SERVICE_INSTALLED = 4697
)

type EventNull struct {
	EventDescription string
}

func ProcessEventItem(LogSource string, EventId int, Message string) interface{} {

	if strings.ToLower(LogSource) == "security" {
		switch EventId {
		case EVENT_LOGON_SUCCESS:
			var ev EventID4624
			ev.ParseEventMessage(Message)
			return ev

		case EVENT_LOGON_FAILED:
			var ev EventID4625
			ev.ParseEventMessage(Message)
			return ev

		case EVENT_EXPLICIT_CRED:
			var ev EventID4648
			ev.ParseEventMessage(Message)
			return ev

		case EVENT_NEW_PROCESS:
			var ev EventID4688
			ev.ParseEventMessage(Message)
			return ev

		case EVENT_SERVICE_INSTALLED:
			var ev EventID4697
			ev.ParseEventMessage(Message)
			return ev
		}

	}

	var ev EventNull
	return ev
}
