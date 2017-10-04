package search

import (
    "encoding/json"
    "fmt"
    "strings"

    nhs "nighthawk/nhstruct"
    nat "nighthawk/audit/audittype"

    elastic "gopkg.in/olivere/elastic.v5"
)

/*
 * Common structure to be used for rendering
 * Timeline and Global search output
 */

type SearchOutput struct {
    Id          string `json:"id"`
    Time        string `json:"time"`
    CaseName    string `json:"case_name"`
    ComputerName string `json:"computer_name"`
    Generator   string `json:"audit_generator"`
    Name        string `json:"name"`
    Path        string `json:"path"`
    Summary     string `json:"summary,omitempty"`
    SiMTime     string `json:"si_modified_time,omitempty"`
    SiATime     string `json:"si_accessed_time,omitempty"`
    SiCTime     string `json:"si_changed_time,omitempty"`
    SiBTime     string `json:"si_birth_time,omitempty"`
    FnMTime     string `json:"fn_modified_time,omitempty"`
    FnATime     string `json:"fn_accessed_time,omitempty"`
    FnCTime     string `json:"fn_changed_time,omitempty"`
    FnBTime     string `json:"fn_birth_time,omitempty"`
}



func (so *SearchOutput) Unmarshal(hit *elastic.SearchHit, is_globalsearch bool) {
    var t nhs.RlRecord
    json.Unmarshal(*hit.Source, &t)

    so.Id = hit.Id
    so.CaseName = t.CaseInfo.CaseName
    so.ComputerName = t.ComputerName
    so.Generator = t.AuditType.Generator

    data, _ := json.Marshal(t.Record)

    switch(t.AuditType.Generator) {
    case nat.RL_PERSISTENCE:
        so.unmarshalPersistence(data)
    case nat.RL_SERVICES:
        if is_globalsearch {so.unmarshalServices(data)}
    case nat.RL_PORTS:
        so.unmarshalPorts(data)
    case nat.RL_USERACCOUNTS:
        if is_globalsearch {so.unmarshalUserAccounts(data)}
    case nat.RL_TASKS:
        so.unmarshalTasks(data)
    case nat.RL_PROCESSMEMORY:
        so.unmarshalProcessMemory(data)
    case nat.RL_PREFETCH:
        so.unmarshalPrefetch(data)
    case nat.RL_REGRAW:
        so.unmarshalRegistryRaw(data)
    case nat.RL_DISKS:
        if is_globalsearch {so.unmarshalDisks(data)}
    case nat.RL_VOLUMES:
        so.unmarshalVolumes(data)
    case nat.RL_URLHISTORY:
        so.unmarshalUrlHistory(data)
    case nat.RL_FILEDLHISTORY:
        so.unmarshalFileDownloadHistory(data)
    case nat.RL_NETWORKDNS:
        if is_globalsearch {so.unmarshalNetworkDns(data)}
    case nat.RL_NETWORKROUTE:
        so.unmarshalNetworkRoute(data)
    case nat.RL_NETWORKARP:
        if is_globalsearch {so.unmarshalNetworkArp(data)}
    case nat.RL_APIFILES:
        so.unmarshalApiFiles(data)
    case nat.RL_RAWFILES:
        so.unmarshalRawFiles(data)
    case nat.RL_EVENTLOGS:
        so.unmarshalEventLogs(data)

    }

    // Truncating large path
    if len(so.Path) > 256 {
        so.Path = so.Path[:255]
    }

    // Escape " in path
    so.Path = strings.Replace(so.Path, "\"","\\\"",-1)

}

func (so *SearchOutput) unmarshalPersistence(data []byte) {
    var record nhs.PersistenceItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.RegPath
    so.Path = record.Path
    so.Summary = record.FileOwner
    so.SiMTime = record.FileModified
    so.SiATime = record.FileAccessed
    so.SiCTime = record.FileChanged
    so.SiBTime = record.FileCreated
}

func (so *SearchOutput) unmarshalServices(data []byte) {
    var record nhs.ServiceItem
    json.Unmarshal(data, &record)

    so.Name = record.Name
    so.Path = record.Path
    so.Summary = record.DescriptiveName
}

func (so *SearchOutput) unmarshalPorts(data []byte) {
    var record nhs.PortItem
    json.Unmarshal(data, &record)

    so.Time = record.JobCreated
    so.Name = record.Process
    so.Path = record.Path
    so.Summary = fmt.Sprintf("Protocol: %s, Local: %s:%d, Remote: %s:%d", record.Protocol, record.LocalIp, record.LocalPort, record.RemoteIp, record.RemotePort)
}

func (so *SearchOutput) unmarshalUserAccounts(data []byte) {
    var record nhs.UserItem
    json.Unmarshal(data, &record)

    so.Name = record.Username
    so.Path = record.HomeDirectory
    so.Summary = record.Description
}

func (so *SearchOutput) unmarshalTasks(data []byte) {
    var record nhs.TaskItem
    json.Unmarshal(data, &record)

    so.Name = record.Name
    so.Path = record.Path
    so.Summary = fmt.Sprintf("LastRun: %s, NextRun: %s", record.MostRecentRunTime, record.NextRunTime)
}

func (so *SearchOutput) unmarshalProcessMemory(data []byte) {
    var record nhs.ProcessItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.Name
    so.Path = record.Path
    so.Summary = record.Arguments
}

func (so *SearchOutput) unmarshalPrefetch(data []byte) {
    var record nhs.PrefetchItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.ApplicationFileName
    so.Path = record.ApplicationFullPath
    so.Summary = record.LastRun
}
func (so *SearchOutput) unmarshalRegistryRaw(data []byte) {
    var record nhs.RegistryItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.ValueName
    so.Path = record.Path
    so.Summary = record.Text
}

func (so *SearchOutput) unmarshalDisks(data []byte) {
    var record nhs.DiskItem
    json.Unmarshal(data, &record)

    so.Name = record.DiskName
    so.Summary = fmt.Sprintf("DiskSize: %s", record.DiskSize)
}

func (so *SearchOutput) unmarshalVolumes(data []byte) {
    var record nhs.VolumeItem
    json.Unmarshal(data, &record)

    so.Time = record.CreationTime
    so.Name = record.VolumeName
    so.Path = record.DevicePath
    so.Summary = fmt.Sprintf("SerialNumber: %s, FileSystem: %s", record.SerialNumber, record.FileSystemName)
}


func (so *SearchOutput) unmarshalUrlHistory(data []byte) {
    var record nhs.UrlHistoryItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.UrlDomain
    so.Path = record.Url
    so.Summary = record.BrowserName
}

func (so *SearchOutput) unmarshalFileDownloadHistory(data []byte) {
    var record nhs.FileDownloadHistoryItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.Filename
    so.Path = record.SourceUrl
    so.Summary = record.TargetDirectory
    so.SiMTime = record.LastModifiedDate
    so.SiATime = record.LastCheckedDate
}

func (so *SearchOutput) unmarshalNetworkDns(data []byte) {
    var record nhs.DnsEntryItem
    json.Unmarshal(data, &record)

    so.Name = record.RecordName
    summary := "IP Addresses: "
    for _,rd := range record.RecordDataList {
        summary = fmt.Sprintf("%s %s", summary, rd.Ipv4Address)
    }
    so.Summary = summary
}

func (so *SearchOutput) unmarshalNetworkRoute(data []byte) {
    var record nhs.RouteEntryItem
    json.Unmarshal(data, &record)

    so.Time = record.JobCreated
    so.Name = record.Interface
    so.Summary = fmt.Sprintf("Destination: %s, Netmask: %s, Gateway: %s", record.Destination, record.Netmask, record.Gateway)
}


func (so *SearchOutput) unmarshalNetworkArp(data []byte) {
    var record nhs.ArpEntryItem
    json.Unmarshal(data, &record)

    so.Name = record.Interface
    so.Summary = fmt.Sprintf("MAC: %s, IP: %s", record.PhysicalAddress, record.Ipv4Address)
}


func (so *SearchOutput) unmarshalApiFiles(data []byte) {
    var record nhs.FileItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.FileName
    so.Path = record.FilePath
    so.Summary = record.FileExtension
    so.SiMTime = record.Modified
    so.SiATime = record.Accessed
    so.SiCTime = record.Changed
    so.SiBTime = record.Created
}

func (so *SearchOutput) unmarshalRawFiles(data []byte) {
    var record nhs.RawFileItem
    json.Unmarshal(data, &record)

    so.Time = record.TlnTime
    so.Name = record.FileName
    so.Path = record.FilePath
    so.Summary = record.FileExtension
    so.SiMTime = record.Modified
    so.SiATime = record.Accessed
    so.SiCTime = record.Changed
    so.SiBTime = record.Created
    so.FnMTime = record.FilenameModified
    so.FnATime = record.FilenameAccessed
    so.FnCTime = record.FilenameChanged
    so.FnBTime = record.FilenameCreated
}
func (so *SearchOutput) unmarshalEventLogs(data []byte) {
    var record nhs.EventLogItem
    json.Unmarshal(data, &record)

    so.Time = record.GenTime
    so.Name = fmt.Sprintf("%s", record.EID)
    so.Path = record.Source
    so.Summary = record.Message
}


func PrintQueryObject(query elastic.Query) {

	jq, _ := query.Source()
	jquery, _ := json.Marshal(jq)
	fmt.Println(string(jquery))
}
