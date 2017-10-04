/* nighthawk.nhstruct.service.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Service
 */

package nhstruct


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
    IsWhitelisted               bool
    IsBlacklisted               bool
    NHScore                     int
    VTResults                   []VirusTotal
    Tag                         string 
    NhComment                   NHComment `json:"Comment"`
}

