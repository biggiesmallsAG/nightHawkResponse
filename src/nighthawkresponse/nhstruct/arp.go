/* nighthawk.nhstruct.arp.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Network ARP
 */
 
package nhstruct

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
    IsGoodEntry         string
    IsWhitelisted       bool
    NHScore             int 
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}


