/*
 * @package     nightHawk
 * @file        audit.go
 * @author      roshan maskey <roshanmaskey@gmail.com>
 *
 * @description 
 *              
 */

 package nightHawk

 import (
    "encoding/json"
    "encoding/xml"
    "io/ioutil"
    "path/filepath"
 )


type RlJsonRecord []byte 


func LoadAuditData(ret OutputDataType, computerName string, caseInfo CaseInformation, targetDir string, auditFile string) (interface{}) {
    var rlRecord []RlJsonRecord

    xmlData,_ := ioutil.ReadFile(filepath.Join(targetDir, auditFile))

    var auditInfo RlAuditType
    xml.Unmarshal(xmlData, &auditInfo)

    var intf interface{}

    switch auditInfo.Generator {
    // __start_of_hx_triage_generator__
    
    case RL_AGENTSTATE:
        var rl RlAgentStateInspector
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        intf = rl

        for _,ml := range rl.EventList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }
    
    case RL_PERSISTENCE:
        var rl RlPersistence
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)

        intf = rl

        for _,ml := range rl.PersistenceList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_SERVICES:
        var rl RlService 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)

        intf = rl

        for _,ml := range rl.ServiceList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_PORTS:
        var rl RlPort 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)

        for _,ml := range rl.PortList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

        intf = rl

    case RL_USERACCOUNTS:
        var rl RlUserAccount 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)

        intf = rl

        for _,ml := range rl.UserList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_PREFETCH:
        var rl RlPrefetch
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
 
        intf = rl

        for _,ml := range rl.PrefetchList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,"", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_TASKS:
        var rl RlTask
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        intf = rl

        for _,ml := range rl.TaskList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_PROCESSMEMORY:
        var rl RlProcessMemory
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)

        intf = rl

        for _,ml := range rl.ProcessList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_REGRAW:
        var rl RlRegistryRaw
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl

        for _,ml := range rl.RegistryList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_SYSTEM:
        var rl RlSystemInfo
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        intf = rl
        jd,_ := json.Marshal(rl)
        rlRecord = append(rlRecord, jd)

    case RL_DISKS:
        var rl RlDisk
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl

        for _,ml := range rl.DiskList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_VOLUMES:
        var rl RlVolume 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl

        for _,ml := range rl.VolumeList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_URLHISTORY:
        var rl RlUrlHistory 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl

        for _,ml := range rl.UrlHistoryList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr,""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_FILEDLHISTORY:
        var rl RlFileDownloadHistory 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl

        for _,ml := range rl.FileDownloadHistoryList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_NETWORKDNS:
        var rl RlNetworkDns
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl
        
        for _,ml := range rl.DnsEntryList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }
        
    case RL_NETWORKROUTE:
        var rl RlNetworkRoute
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl

        for _,ml := range rl.RouteEntryList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    case RL_NETWORKARP:
        var rl RlNetworkArp
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        
        intf = rl

        for _,ml := range rl.ArpList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml

            //jd,_ := json.MarshalIndent(rr, ""," ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    // __end_of_hx_triage_generator__ 

    // Additional Redline modules
    case RL_APIFILES:
        var rl RlApiFile 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        intf = rl
        for _,ml := range rl.FileList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml 

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

     case RL_RAWFILES:
        var rl RlRawFile 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        intf = rl
        for _,ml := range rl.FileList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml 

            //jd,_ := json.MarshalIndent(rr, "", " ")
            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

     case RL_EVENTLOGS:
        var rl RlEventLog 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        intf = rl

        for _,ml := range rl.EventList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml 

            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }

    /*
    case RL_KERNELHOOK:
        var rl RlKernelHook 
        rl.ParseAuditData(computerName, caseInfo, auditInfo, xmlData)
        intf = rl
        for _,ml := range rl.HookList {
            var rr RlRecord
            rr.ComputerName = rl.ComputerName
            rr.CaseInfo = rl.CaseInfo
            rr.AuditType = rl.AuditType
            rr.Record = ml 

            jd,_ := json.Marshal(rr)
            rlRecord = append(rlRecord, jd)
        }
    */

    } // __end_of_swtich__  

    switch ret {
    case MOD_JSON:
        intf = rlRecord
    case MOD_OBJ:
        break
    }
    return intf
 }
