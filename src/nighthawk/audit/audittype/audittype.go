/* audittype.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Contains AuditType Information
 */

package audittype

const (
	RL_AGENTSTATE    = "stateagentinspector" // This is HX triage specific not collected as a part of normal Redline Triage
	RL_PERSISTENCE   = "w32scripting-persistence"
	RL_SERVICES      = "w32services"
	RL_PORTS         = "w32ports"
	RL_USERACCOUNTS  = "w32useraccounts"
	RL_TASKS         = "w32tasks"
	RL_PROCESSMEMORY = "w32processes-memory"
	RL_PREFETCH      = "w32prefetch"
	RL_REGRAW        = "w32registryraw"
	RL_SYSTEM        = "w32system"
	RL_SYSINFO       = "sysinfo"
	RL_DISKS         = "w32disks"
	RL_VOLUMES       = "w32volumes"
	RL_URLHISTORY    = "urlhistory"
	RL_FILEDLHISTORY = "filedownloadhistory"
	RL_NETWORKDNS    = "w32network-dns"
	RL_NETWORKROUTE  = "w32network-route"
	RL_NETWORKARP    = "w32network-arp"
	RL_APIFILES      = "w32apifiles"
	RL_RAWFILES      = "w32rawfiles"
	RL_HIVELIST      = "w32hivelist"
	RL_SYSTEMRESTORE = "w32systemrestore"
	RL_KERNELHOOK    = "w32kernel-hookdetection"
	RL_EVENTLOGS     = "w32eventlogs"
)
