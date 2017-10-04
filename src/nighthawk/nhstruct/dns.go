/* nighthawk.nhstruct.dns.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Network DNS
 */

package nhstruct

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
    IsWhitelisted       bool
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}
