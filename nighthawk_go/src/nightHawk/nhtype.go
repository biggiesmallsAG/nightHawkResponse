/*
 *@package  nightHawk
 *@file     nhtype.go
 *@author   roshan maskey <roshanmaskey@gmail.com>
 *
 *@description  nightHawk Response Reline DataType in Gopher
 */


 package nightHawk


type DataFileType int 
const (
    MOD_XML DataFileType = 1 + iota
    MOD_ZIP
    MOD_MANS
    MOD_REDDIR
)

type OutputDataType int 
const (
    MOD_JSON OutputDataType = 1 + iota
    MOD_CSV
    MOD_OBJ
)


/// General support structure
type RlAuditType struct {
    Generator           string `xml:"generator,attr"`
    GeneratorVersion    string `xml:"generatorVersion,attr"` 
}


type CaseInformation struct {
    CaseName            string `json:"case_name"`
    CaseDate            string `json:"case_date"`
    CaseAnalyst         string `json:"case_analyst"`
}


type RlRecord struct {
    ComputerName        string 
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    Record              interface{}
}


type NHComment struct {
    Date            string `json:"Date"`
    Analyst         string `json:"Analyst"`
    Comment         string `json:"Comment"`
}


/* __start_of_stageagentinspector__ */
type EventDetail struct {
    Name            string `xml:"name"`
    Value           string `xml:"value"`
}

type AgentEventItem struct {
    Timestamp           string `xml:"timestamp"`
    EventType           string `xml:"eventType"`
    Details             []EventDetail `xml:"details>detail"`
}

type RlAgentStateInspector struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    EventList           []AgentEventItem `xml:"eventItem"`
}

/* __start_of_w32scripting-persistence */
type VirusTotal struct {
    CheckedDate     string
    VTScore         int
}

type PEChecksum struct {
    PEFileRaw           int
    PEFileAPI           int
    PEComputedAPI       int
}

type DigitalSignature struct {
    SignatureExists     bool
    SignatureVerified   bool
    Description         string 
    CertificateSubject  string
    CertificateIssuer   string
}

type PEInfo struct {
    Type                string
    Subsystem           string
    BaseAddress         string
    PETimeStamp         string
    PeChecksum          PEChecksum `xml:"PEChecksum"`
    DigitalSignature    DigitalSignature `xml:"DigitalSignature"`
}

type FileItem struct {
    JobCreated      string `xml:"created,attr"`
    TlnTime         string `json:"TlnTime"`
    DevicePath      string
    Path            string  `xml:"FullPath"`
    Drive           string
    FilePath        string
    FileName        string
    FileExtension   string
    SizeInBytes     int 
    Created         string
    Modified        string
    Accessed        string
    Changed         string
    FileAttributes  string
    Username        string
    SecurityID      string
    SecurityType    string
    Md5sum          string
    PeInfo          PEInfo `xml:"PEInfo"`
    IsGoodHash      string
    NHScore         int
    VTResults       []VirusTotal 
    Tag             string 
    NhComment       NHComment `json:"Comment"`
}

type RegistryItem struct {
    JobCreated      string `xml:"created,attr"`
    TlnTime         string `json:"TlnTime"`
    KeyPath         string
    Type            string
    Modified        string
    ValueName       string
    Username        string
    Path            string
    Text            string
    ReportedLengthInBytes   int
    Hive            string
    SecurityID      string
    IsKnownKey      string
    NHScore         int
    Tag             string 
    NhComment       NHComment `json:"Comment"`
}

type PersistenceItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    PersistenceType     string
    RegPath             string
    RegOwner            string
    RegModified         string
    Path                string `xml:"FilePath"`
    FileOwner           string
    FileCreated         string
    FileModified        string
    FileAccessed        string
    FileChanged         string
    Md5sum              string `xml:"md5sum"`
    File                FileItem `xml:"FileItem"`
    Registry            RegistryItem `xml:"RegistryItem"`
    StackPath           string
    IsGoodPersistence   string 
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlPersistence struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    PersistenceList     []PersistenceItem `xml:"PersistenceItem"`
}

/* __start_of_w32services__ */
type ServiceItem struct {
    JobCreated                  string `xml:"created,attr"`
    Name                        string `xml:"name"`
    DescriptiveName             string `xml:"descriptiveName"`
    Description                 string `xml:"description"`
    Mode                        string `xml:"mode"`
    StartedAs                   string `xml:"startedAs"`
    Path                        string `xml:"path"`
    PathMd5sum                  string `xml:"pathmd5sum"`
    PathSignatureExists         bool    `xml:"pathSignatureExists"`
    PathSignatureVerified       bool    `xml:"pathSignatureVerified"`
    PathSignatureDescription    string `xml:"pathSignatureDescription"`
    PathCertificateSubject      string `xml:"pathCertificateSubject"`
    PathCertificateIssuer       string `xml:"pathCertificateIssuer"`
    Arguments                   string `xml:"arguments"`
    Status                      string `xml:"status"`
    Pid                         int `xml:"pid"`
    Type                        string `xml:"type"`
    IsGoodService               string
    NHScore                     int
    VTResults                   []VirusTotal
    Tag                         string 
    NhComment                   NHComment `json:"Comment"`
}

type RlService struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    ServiceList         []ServiceItem `xml:"ServiceItem"`
}

/* __start_of_w32ports__ */
type PortItem struct {
    JobCreated          string `xml:"created,attr"`
    Pid                 int `xml:"pid"`
    Process             string `xml:"process"`
    Path                string `xml:"path"`
    State               string `xml:"state"`
    LocalIp             string `xml:"localIP"`
    RemoteIp            string `xml:"remoteIP"`
    LocalPort           int `xml:"localPort"`
    RemotePort          int `xml:"remotePort"`
    Protocol            string `xml:"protocol"`
    IsGoodPort          string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlPort struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    PortList            []PortItem `xml:"PortItem"`
}


/* __start_of_w32useraccounts__ */

type UserItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Username            string `xml:"Username"`
    SecurityID          string `xml:"SecurityID"`
    SecurityType        string `xml:"SecurityType"`
    Fullname            string `xml:"fullname"`
    Description         string `xml:"description"`
    HomeDirectory       string `xml:"homedirectory"`
    ScriptPath          string `xml:"scriptpath"`
    LastLogin           string `xml:"lastlogin"`
    Disabled            bool `xml:"disabled"`
    LockedOut           bool `xml:"lockedout"`
    PasswordRequired    bool `xml:"passwordrequired"`
    UserPasswordAge     string `xml:"userpasswordage"`
    Tag                 string 
    NhComment           NHComment `json:"Comment"`
}

type RlUserAccount struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    UserList            []UserItem `xml:"UserItem"`
}


/* __start_of_w32tasks__ */
type TaskAction struct {
    ActionType          string
    Path                string `xml:"ExecProgramPath"`
    ExecProgramMd5sum   string
    ExecArguments       string 
    ExecWorkingDirectory string 
    COMClassId          string
    COMData             string
}

type TaskTrigger struct {
    TriggerEnabled      bool
    TriggerBegin        string
    TriggerFrequency    string
    TriggerDelay        string
    TriggerSubscription string
}

type TaskItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Name                string 
    VirtualPath         string
    ExitCode            string
    TaskComment         string `xml:"Comment"`
    CreationDate        string `xml:"CreationDate"`
    Creator             string
    MaxRunTime          string 
    Flag                string
    AccountName         string
    AccountRunLevel     string
    AccountLogonType    string
    MostRecentRunTime   string
    NextRunTime         string
    Status              string
    ActionList          []TaskAction `xml:"Action"`
    TriggerList         []TaskTrigger `xml:"Trigger"`
    IsGoodTask          string 
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlTask struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    TaskList            []TaskItem `xml:"TaskItem"`
}

/* __start_of_w32process-memory */
type Handle struct {
    Index               int
    AccessMask          int
    ObjectAddress       string
    HandleCount         int 
    PointerCount        int
    Type                string
    Name                string
    IsGoodHandle        string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type MemorySection struct {
    Protection          string
    RegionStart         int
    RegionSize          int
    Mapped              bool
    RawFlags            string
    IsGoodSection       string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type ProcessItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Pid                 int     `xml:"pid"`
    ParentPid           int     `xml:"parentpid"`
    Path                string  `xml:"path"`
    Name                string  `xml:"name"`
    Arguments           string  `xml:"arguments"`
    Username            string  `xml:"Username"`
    SecurityID          string  `xml:"SecurityID"`
    SecurityType        string  `xml:"SecurityType"`
    StartTime           string  `xml:"startTime"`
    KernelTime          string  `xml:"kernelTime"`
    UserTime            string  `xml:"userTime"`
    HandleList          []Handle `xml:"HandleList>Handle"`
    //SectionList       []MemorySection `xml:"SectionList>MemorySection"`
    IsGoodProcess       string
    NHScore             int
    Tag                 string 
    NhComment           NHComment `json:"Comment"`
}

type RlProcessMemory struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    ProcessList         []ProcessItem `xml:"ProcessItem"`
}


/* __start_of_w32prefetch__ */
type PrefetchItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Path                string `xml:"FullPath"`
    Created             string
    SizeInBytes         int
    ReportedSizeInBytes int
    ApplicationFileName string
    LastRun             string
    TimesExecuted       int
    AccessedFileList    []string `xml:"AccessedFileList>AccessedFile"`
    ApplicationFullPath string
    VolumeDevicePath    string `xml:"DevicePath"`
    VolumeCreationTime  string `xml:"CreationTime"`
    VolumeSerialNumber  string `xml:"SerialNumber"`
    IsGoodPrefetch      string
    NHScore             int
    Tag                 string 
    NhComment           NHComment `json:"Comment"`
}

type RlPrefetch struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    PrefetchList        []PrefetchItem `xml:"PrefetchItem"`
}


/* __start_of_w32registryraw__ */
type RlRegistryRaw struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    RegistryList        []RegistryItem `xml:"RegistryItem"`
}


/* __start_of_w32sytem__ */
type IpInfo struct {
    IpAddress           string `xml:"ipAddress"`
    SubnetMask          string `xml:"subnetMask"`
    Ipv6Address         string `xml:"ipv6Address"`
}
type NetworkInfo struct {
    Adapter             string `xml:"adapter"`
    Description         string `xml:"description"`
    Mac                 string `xml:"MAC"`
    IpInfoList          []IpInfo `xml:"ipInfo"`
    IpGateway           []string    `xml:"ipGatewayArray>ipGateway"`
    DhcpServer          []string    `xml:"dhcpServerArray>dhcpServer"`
    DhcpLeaseObtained   string  `xml:"dhcpLeaseObtained"`
    DhcpLeaseExpires    string  `xml:"dhcpLeaseExpires"`    
}

type SystemInfoItem struct {
    JobCreated          string `xml:"created,attr"`
    Machine             string `xml:"machine"`
    Uptime              string `xml:"uptime"`
    ContainmentState    string `xml:"containmentState"`
    BiosDate            string `xml:"biosInfo>biosDate"`
    BiosVersion         string `xml:"biosInfo>biosVersion"`
    BiosType            string `xml:"biosInfo>biosType"`
    Directory           string `xml:"directory"`
    Drives              string `xml:"drives"`
    ProcType            string `xml:"procType"`
    RegOrg              string `xml:"regOrg"`
    RegOwner            string `xml:"regOwner"`
    Processor           string `xml:"processor"`
    ProcVmGuest         string `xml:"vmGuest"`
    ProcVirtualization  string `xml:"virtualization"`
    ProcLpcDevice       string `xml:"lpcDevice"`
    OS                  string `xml:"OS"`
    ProductName         string `xml:"productName"`
    BuildNumber         int     `xml:"buildNumber"`
    ProductID           string `xml:"productID"`
    InstallDate         string `xml:"installDate"`
    OsBitness           string `xml:"OSbitness"`
    TimezoneDst         string `xml:"timezoneDST"`
    TimezoneStandard    string `xml:"timezoneStandard"`
    Timezone            string `xml:"timezone"`
    GmtOffset           string `xml:"gmtoffset"`
    ClockSkew           string `xml:"clockSkew"`
    Date                string `xml:"date"`
    StateAgentStatus    string `xml:"stateAgentStatus"`
    Hostname            string `xml:"hostname"`
    Domain              string `xml:"domain"`
    PrimaryIpv4Address  string `xml:"primaryIpv4Address"`
    PrimaryIpAddress    string `xml:"primaryIpAddress"`
    Mac                 string `xml:"MAC"`
    TotalPhysical       int     `xml:"totalphysical"`
    AvailPhysical       int     `xml:"availphysical"`
    User                string  `xml:"user"`
    LoggedOnUser        string  `xml:"loggedOnUser"`
    AppVersion          string  `xml:"appVersion"`
    Platform            string  `xml:"platform"`
    AppCreated          string  `xml:"appCreated"`
    NetworkList         []NetworkInfo `xml:"networkArray>networkInfo"`
}

type RlSystemInfo struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    SystemInfo          SystemInfoItem `xml:"SystemInfoItem"`
}

/* __start_of_w32disks__ */
type PartitionItem struct {
    PartitionNumber     int 
    PartitionOffset     int
    PartitionLength     int
    PartitionType       string
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type DiskItem struct {
    JobCreated          string `xml:"created,attr"`
    DiskName            string
    DiskSize            string
    PartitionList       []PartitionItem `xml:"PartitionList>Partition"`
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlDisk struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    DiskList            []DiskItem `xml:"DiskItem"`
}

/* __start_of_w32volumes__ */
type VolumeItem struct {
    JobCreated                      string `xml:"created,attr"`
    VolumeName                      string 
    DevicePath                      string
    Type                            string
    Name                            string
    SerialNumber                    string
    FileSystemFlags                 string
    FileSystemName                  string
    ActualAvailableAllocationUnits  int
    TotalAllocationUnits            int
    BytesPerSector                  int
    SectorsPerAllocationUnit        int
    CreationTime                    string
    IsMounted                       bool
    Tag                             string
    NhComment                       NHComment `json:"Comment"`
}

type RlVolume struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    VolumeList          []VolumeItem `xml:"VolumeItem"`
}

/* __start_of_w32 */
type UrlHistoryItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Profile             string `xml:"Profile"`
    BrowserName         string `xml:"BrowseName"`
    BrowserVersion      string `xml:"BrowserVersion"`
    Username            string `xml:"Username"`
    Url                 string `xml:"URL"`
    UrlHostname         string 
    UrlDomain           string 
    PageTitle           string `xml:"PageTitle"`
    Hidden              bool    `xml:"Hidden"`
    LastVisitDate       string  `xml:"LastVisitDate"`
    VisitFrom           string `xml:"VisitFrom"`
    VisitType           string `xml"VisitType"`
    VisitCount          int `xml:"VisitCount"`
    IsGoodDomain        string
    IsNewDomain         string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlUrlHistory struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    UrlHistoryList      []UrlHistoryItem `xml:"UrlHistoryItem"`
}


/* __start_of_filedownloadhistory */
type FileDownloadHistoryItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Profile             string
    BrowserName         string
    BrowserVersion      string
    Username            string
    DownloadType        string 
    Filename            string `xml:"FileName"`
    SourceUrl           string `xml:"SourceURL"`
    UrlHostname         string
    UrlDomain           string 
    TargetDirectory     string `xml:"TargetDirectory"`
    FullHttpHeader      string
    LastModifiedDate    string
    BytesDownloaded     int 
    MaxBytes            int
    CacheFlags          string
    CacheHitCount       int
    LastCheckedDate     string
    IsGoodFile          string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlFileDownloadHistory struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    FileDownloadHistoryList []FileDownloadHistoryItem `xml:"FileDownloadHistoryItem"`
}


/* __start_of_w32network-dns */
type RecordData struct {
    Ipv4Address         string  `xml:"IPv4Address"`
}

type DnsEntryItem struct {
    //JobCreated          string `xml:"created,attr"`
    Host                string 
    RecordName          string 
    RecordType          string 
    TimeToLive          string       
    Flags               string
    DataLength          string
    RecordDataList      []RecordData `xml:"RecordData"`
    IsGoodEntry         string
    NHScore             int
    Tag                 string 
    NhComment           NHComment `json:"Comment"`
}

type RlNetworkDns struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    DnsEntryList        []DnsEntryItem `xml:"DnsEntryItem"`
}


/* __start_of_32network-route */
type RouteIpv6Item struct {
    ValidLifeTime       string
    PreferredLifeTime   string
    IsLoopback          bool
    IsAutoconfiguredAddress bool
    IsPublished         bool
    IsImortal           bool
    Origin              string
}

type RouteEntryItem struct {
    JobCreated          string `xml:"created, attr"`
    IsIpv6              bool `xml:"IsIPv6"`
    //Ipv4Route
    Interface           string 
    Destination         string
    Netmask             string
    RouteType           string
    //Common for Ipv4 and Ipv6 route
    Gateway             string
    Protocol            string
    RouteAge            string
    Metric              int
    // IPv6 Route
    ValidLifeTime       string
    PreferredLifeTime   string
    IsLoopback          bool
    IsAutoconfiguredAddress bool
    IsPublished         bool
    IsImortal           bool
    Origin              string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlNetworkRoute struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    RouteEntryList      []RouteEntryItem `xml:"RouteEntryItem"`
}

/* __start_of_w32network-arp */
type ArpEntryItem struct {
    JobCreated          string `xml:"created,attr"`
    Interface           string
    // IPv4 Structure
    PhysicalAddress     string `xml:"PhysicalAddress,omitempty"`
    Ipv4Address         string `xml:"IPv4Address,omitempty"`
    CacheType           string `xml:"CacheType,omitempty"`
    // IPv6 Structure
    Ipv6Address         string `xml:"IPv6Address,omitempty"`
    InterfaceType       string `xml:"InterfaceType,omitempty"`
    State               string `xml:"State,omitempty"`
    IsRouter            bool `xml:"IsRouter,omitempty"`
    LastReachable       string `xml:"LastReachable,omitempty"`
    LastUnreachable     string `xml:"LastUnreachable,omitempty"`
    NHScore             int 
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type RlNetworkArp struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    ArpList             []ArpEntryItem `xml:"ArpEntryItem"`
}


type RlApiFile struct {
    ComputerName        string
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    FileList            []FileItem `xml:"FileItem"`
}


/* __start_of_32rawfiles__ */

type RawFileItem struct {
    JobCreated      string `xml:"created,attr"`
    TlnTime         string `json:"TlnTime"`
    DevicePath      string
    Path            string  `xml:"FullPath"`
    Drive           string
    FilePath        string
    FileName        string
    FileExtension   string
    SizeInBytes     int 
    Created         string
    Modified        string
    Accessed        string
    Changed         string
    FilenameCreated string 
    FilenameModified string
    FilenameAccessed string
    FilenameChanged string
    FileAttributes  string
    Username        string
    SecurityID      string
    SecurityType    string
    Md5sum          string
    PeInfo          PEInfo `xml:"PEInfo"`
    IsGoodHash      string
    NHScore         int
    VTResults       []VirusTotal 
    Tag             string 
    NhComment       NHComment `json:"Comment"`
}

type RlRawFile struct {
    ComputerName        string 
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    FileList            []RawFileItem `xml:"FileItem"`
}


/* __start_of_w32eventlogs__ */


type EventLogItem struct {
    Log             string  `xml:"log"`
    Source          string  `xml:"source"`
    Index           int     `xml:"index"`
    EID             int     `xml:"EID"`
    Type            string  `xml:"type"`
    GenTime         string  `xml:"genTime"`   
    WriteTime       string  `xml:"writeTime"`
    Machine         string  `xml:"machine"`
    ExecutionProcessId int  `xml:"ExecutionProcessId"`
    ExecutionThreadId int   `xml:"ExecutionThreadId"`
    Message         string  `xml:"message"`
    MessageDetail   interface{}
    Category        string  `xml:"category"` 
    NHScore         int 
    Tag             string 
    NhComment       NHComment `json:"Comment"`
}

type RlEventLog struct {
    ComputerName        string 
    CaseInfo            CaseInformation
    AuditType           RlAuditType
    EventList           []EventLogItem  `xml:"EventLogItem"`
}