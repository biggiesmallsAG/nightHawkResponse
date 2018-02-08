/* nighthawk.nhstruct.system.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for System Information
 */

package nhstruct

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
