/*
 * This module is used by AuditParser parsing Windows Event Logs 
 *
 * author: roshan maskey <0xredskull>
 */

package winevt

import (
    "strings"
    "strconv"
)

const (
    EVENT_LOGON_SUCCESS             = 4624
    EVENT_LOGON_FAILED              = 4625
    EVENT_EXPLICIT_CRED             = 4648
    EVENT_NEW_PROCESS               = 4688
    EVENT_SERVICE_INSTALLED         = 4697
)



func ProcessEventItem(LogSource string, EventId int, Message string) interface{} {

    if strings.ToLower(LogSource) == "security" {
        switch EventId {
        case EVENT_LOGON_SUCCESS: 
            var ev EventNewAccount
            ev.ParseEventMessage(Message)
            return ev 

        case EVENT_LOGON_FAILED:
            var ev EventLogonFailed 
            ev.ParseEventMessage(Message)
            return ev 
        
        case EVENT_EXPLICIT_CRED:
            var ev EventExplictCred 
            ev.ParseEventMessage(Message)
            return ev

        case EVENT_NEW_PROCESS:
            var ev EventNewProcess 
            ev.ParseEventMessage(Message)
            return ev 

        case EVENT_SERVICE_INSTALLED:
            var ev EventServiceInstall 
            ev.ParseEventMessage(Message)
            return ev 
        }

    }
    return ""  
}



// Parse AccountInfo defined in winevt_type.go 
// startIndex: Index in arrary data where account information starts
// stopIndex: Index in array data where account information stops
// data: array containing EventMessage splitted using line terminator
func (acc *AccountInfo)ReadAccountInfo(startIndex int, stopIndex int, data []string) {
    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)

        if len(keyval) == 2 {
            switch keyval[0] {
            case "Security ID":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.SecurityId = value 

            case "Account Name":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.AccountName = value 

            case "Account Domain":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.AccountDomain = value 

            case "Logon ID":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.LogonId = value 

            case "Logon GUID":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.LogonGuid = value 

            case "Linked Logon ID":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.LinkedLogonId = value 

            case "Network Account Name":
                var value string = "" 
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.NetworkAccountName = value 

            case "Network Account Domain":
                var value string = "" 
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                acc.NetworkAccountDomain = value 
            } // __end_switch_
        }
    } // _end_for_
}




// Parse AccountInfo defined in winevt_type.go 
// startIndex: Index in arrary data where Logon information starts
// stopIndex: Index in array data where Logon information stops
// data: array containing EventMessage splitted using line terminator
func (logon *LogonInformation)ReadLogonInfo(startIndex int, stopIndex int, data []string) {
    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)
        if len(keyval) == 2 {
            switch keyval[0] {
            case "Logon Type":
                var value int = 0
                if keyval[1] != "-" {
                    value,_ = strconv.Atoi(strings.TrimSpace(keyval[1]))
                }
                //value,_ = strconv.Atoi(value)
                logon.LogonType = value 
                logon.LogonTypeDetail = GetLogonTypeDetail(logon.LogonType)

            case "Restricted Admin Mode":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                logon.RestrictedAdminMode = value 

            case "Virtual Account":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                logon.VirtualAccount = value 

            case "Elevated Token":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                logon.ElevatedToken = value 
            } // _end_switch_
        }
    } //_end_for_
}

/// This function parses Network Information in Windows EventLog Message
func (netinfo *NetworkInformation)ReadNetworkInfo(startIndex int, stopIndex int, data []string) {
    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)
        if len(keyval) == 2 {
            
            if keyval[0] == "Workstation Name" || keyval[0] == "Source Workstation Name" {
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                netinfo.WorkstationName = value 
            }
                

            if keyval[0] == "Source Network Address" || keyval[0] == "Network Address" {
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                netinfo.SourceNetworkAddress = value 
            }
                
            if keyval[0] == "Source Port" || keyval[0] == "Source Port" {
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                netinfo.SourcePort = value 
            }
                
        } // _end_if_len_
    }
}

/// This function parses Process Information in Windows EventLog Message
func (procinfo *ProcessInformation)ReadProcessInfo(startIndex int, stopIndex int, data []string) {
    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)
        if len(keyval) == 2 {
            switch keyval[0] {
            case "Process ID":
                var value string = "0x00"
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                procinfo.ProcessIdHex = value 
                value = strings.Replace(value, "0x", "", -1)
                pid,_ := strconv.ParseInt(value, 16, 64)
                procinfo.ProcessId = int(pid)

            case "Process Name":
                procinfo.ProcessName = strings.TrimSpace(keyval[1])

            }
        }
    }
}

/// This function returns Logon Type
func GetLogonTypeDetail(logonType int) string {
    var logonTypeDesc string = "UNKNOWN"
    switch logonType {
    case 2:
        logonTypeDesc = "INTERACTIVE"
    case 3:
        logonTypeDesc = "NETWORK"
    case 4:
        logonTypeDesc = "BATCH_OR_SCHEDULE_TASK"
    case 5:
        logonTypeDesc = "SERVICE"
    case 7:
        logonTypeDesc = "UNLOCK"
    case 8:
        logonTypeDesc = "NETWORK_CLEAR_TEXT"
    case 9:
        logonTypeDesc = "NEW_CREDENTIAL"
    case 10:
        logonTypeDesc = "REMOTE_INTERACTIVE"
    case 11:
        logonTypeDesc = "CACHED_INTERACTIVE"
    } 
    return logonTypeDesc
}


 func GetSectionIndex(fields []string, data []string) map[string]int {
    sectionMap := make(map[string]int)

    for i := 0; i < len(data); i++ {
        if len(fields) == len(sectionMap) {
            break
        }

        for _,field := range fields {
            if strings.HasPrefix(data[i], field) {
                sectionMap[field] = i
            }
        }
    }

    return sectionMap
 }


