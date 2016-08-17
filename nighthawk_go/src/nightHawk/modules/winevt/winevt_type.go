/*
 * This file contains structures that holds
 * information parsing Windows Event Message
 *
 * author: roshan maskey <0xredskull>
*/
package winevt


// Common Subject structure
type AccountInfo struct {
    SecurityId              string `json:"SecurityId",omitempty`
    AccountName             string `json:"AccountName",omitempty`
    AccountDomain           string `json:"AccountDomain",omitempty`
    LogonId                 string `json:"LogonId",omitempty`
    LogonGuid               string `json:"LogonGuid",omitempty`
    LinkedLogonId           string `json:"LinkedLogonId",omitempty`         // Windows 10 and Windows 2016
    NetworkAccountName      string `json:"NetworkAccountName",omitempty`    // Windows 10 and Windows 2016
    NetworkAccountDomain    string `json:"NetworkAccountDomain",omitempty`  // Windows 10 and Windows 2016
   
}

type LogonInformation struct {
    LogonType               int     `json:"LogonType",omitempty`
    LogonTypeDetail         string `json:"LogonTypeDetail",omitempty`
    RestrictedAdminMode     string `json:"RestrictedAdminMode",omitempty`   // Windows 10 and Windows 2016
    VirtualAccount          string `json:"VirtualAccount",omitempty`        // Windows 10 and Windows 2016
    ElevatedToken           string `json:"ElevatedToken",omitempty`         // Windows 10 and Windows 2016
}

type NetworkInformation struct {
    WorkstationName         string 
    SourceNetworkAddress    string 
    SourcePort              string 
}

type ProcessInformation struct {
    ProcessIdHex            string 
    ProcessId               int 
    ProcessName             string 
    ProcessCommandline      string `json:"ProcessCommandline",omitempty`
}


// Structure to hold attributes of Windows
// Security Event ID 4624
type EventNewAccount struct {
    Subject                 AccountInfo `json:"Subject"`
    LogonInfo               LogonInformation `json:"LogonInformation"`
    NewAccount              AccountInfo `json:"Account"`
    ProcessInfo             ProcessInformation `json:"ProcessInformation"`
    NetworkInfo             NetworkInformation `json:"NetworkInformation"`
    //ProcessId               string 
    //ProcessName             string 
    //WorkstationName         string 
    //SourceNetworkAddress    string 
    //SourcePort              string 
    LogonProcess            string 
    AuthenticationPackage   string 
    TransitedServices       string 
    PackageName             string 
    KeyLength               int 
}

type EventLogonFailed struct {
    Subject                 AccountInfo         `json:"Subject"`
    LogonInfo               LogonInformation    `json:"LogonInformation"`
    Account                 AccountInfo         `json:"Account"`
    FailureReason           string 
    FailureStatus           string 
    FailureSubStatus        string 
    FailureDescription      string 
    ProcessInfo             ProcessInformation  `json:"ProcessInformation"`
    NetworkInfo             NetworkInformation  `json:"NetworkInformation"`
    LogonProcess            string 
    AuthenticationPackage   string 
    TransitedServices       string
    PackageName             string 
    KeyLength               int 
}


// Structure to hold attributes of Windows 
// Security EventID 4648
type EventExplictCred struct {
    Subject                 AccountInfo `json:"Subject"`
    Account                 AccountInfo `json:"Account"`
    TargetServerName        string 
    AdditionalInfo          string 
    ProcessId               string 
    ProcessName             string 
    NetworkAddress          string 
    Port                    string 
    ProcessInfo             ProcessInformation `json:"ProcessInformation"`
    NetworkInfo             NetworkInformation `json:"NetworkInformation"`
}


/// Structure for Windows Security Event ID 4688
/// 4688: A new process has been created
type EventNewProcess struct {
    Subject                 AccountInfo     `json:"Subject"`
    NewProcessIdHex         string          
    NewProcessId            int     
    NewProcessName          string      
    TokenElevationType      string 
    TokenElevationDesc      string 
    MandatoryLabel          string 
    CreatorProcessIdHex     string 
    CreatorProcessId        int 
    CreatorProcessName      string 
    ProcessCommandline      string          
}


// Struct for Window EventID 4697
// A service was installed in the system
// Thanks for Phil Kealy for suggesting to include it
type EventServiceInstall struct {
    Subject                 AccountInfo     `json:"Subject"`
    ServiceName             string 
    ServiceFileName         string 
    ServiceType             string 
    ServiceTypeDesc         string 
    ServiceStartType        string 
    ServiceStartTypeDesc    string 
    ServiceAccount          string 
}

