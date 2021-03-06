/* parser.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Entrypoint to parse all audit files
 */

package audit

import (
	"fmt"
	"nighthawkresponse/audit/audittype"
	"nighthawkresponse/audit/parser"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
	"strings"
)

func ParseAuditFile(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	switch strings.ToLower(auditinfo.Generator) {
	case audittype.RL_AGENTSTATE:
		parser.ParseAgentStateInspection(caseinfo, auditinfo, auditfile)

	case audittype.RL_PERSISTENCE:
		parser.ParsePersistence(caseinfo, auditinfo, auditfile)

	case audittype.RL_SERVICES:
		parser.ParseServices(caseinfo, auditinfo, auditfile)

	case audittype.RL_EVENTLOGS:
		parser.ParseEventLogs(caseinfo, auditinfo, auditfile)

	case audittype.RL_PORTS:
		parser.ParsePorts(caseinfo, auditinfo, auditfile)

	case audittype.RL_USERACCOUNTS:
		parser.ParseUserAccounts(caseinfo, auditinfo, auditfile)

	case audittype.RL_PREFETCH:
		parser.ParsePrefetch(caseinfo, auditinfo, auditfile)

	case audittype.RL_TASKS:
		parser.ParseTasks(caseinfo, auditinfo, auditfile)

	case audittype.RL_PROCESSMEMORY, audittype.RL_PROCESS_API:
		parser.ParseProcessMemory(caseinfo, auditinfo, auditfile)

	case audittype.RL_REGRAW:
		parser.ParseRegistry(caseinfo, auditinfo, auditfile)

	case audittype.RL_SYSTEM, audittype.RL_SYSINFO:
		parser.ParseSystemInfo(caseinfo, auditinfo, auditfile)

	case audittype.RL_DISKS:
		parser.ParseDisk(caseinfo, auditinfo, auditfile)

	case audittype.RL_VOLUMES:
		parser.ParseVolumes(caseinfo, auditinfo, auditfile)

	case audittype.RL_RAWFILES:
		parser.ParseRawFiles(caseinfo, auditinfo, auditfile)

	case audittype.RL_APIFILES:
		parser.ParseApiFiles(caseinfo, auditinfo, auditfile)

	case audittype.RL_URLHISTORY:
		parser.ParseUrlHistory(caseinfo, auditinfo, auditfile)

	case audittype.RL_FILEDLHISTORY:
		parser.ParseFileDownloadHistory(caseinfo, auditinfo, auditfile)

	case audittype.RL_NETWORKDNS:
		parser.ParseNetworkDns(caseinfo, auditinfo, auditfile)

	case audittype.RL_NETWORKROUTE:
		parser.ParseNetworkRoute(caseinfo, auditinfo, auditfile)

	case audittype.RL_NETWORKARP:
		parser.ParseNetworkArp(caseinfo, auditinfo, auditfile)

	case audittype.RL_DRIVER_MODULES:
		parser.ParseDriverModules(caseinfo, auditinfo, auditfile)

	default:
		// _rm_note> Do not exit on this ERROR
		// This error helps to identify unparsed audit type
		nhlog.LogMessage("ParseAuditFile", "ERROR", fmt.Sprintf("Failed to parse unsupported AuditType %s", auditinfo.Generator))
	}
}
