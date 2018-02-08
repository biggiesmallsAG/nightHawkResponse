/* nighthawk.nhstruct.urlhistory.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Browser URL History
 */

package nhstruct

type UrlHistoryItem struct {
	JobCreated     string `xml:"created,attr"`
	TlnTime        string `json:"TlnTime"`
	Profile        string `xml:"Profile"`
	BrowserName    string `xml:"BrowseName"`
	BrowserVersion string `xml:"BrowserVersion"`
	Username       string `xml:"Username"`
	Url            string `xml:"URL"`
	UrlHostname    string
	UrlDomain      string
	PageTitle      string `xml:"PageTitle"`
	Hidden         bool   `xml:"Hidden"`
	LastVisitDate  string `xml:"LastVisitDate"`
	VisitFrom      string `xml:"VisitFrom"`
	VisitType      string `xml"VisitType"`
	VisitCount     int    `xml:"VisitCount"`
	IsGoodDomain   string
	IsNewDomain    string
	IsWhitelisted  bool
	NHScore        int
	Tag            string
	NhComment      NHComment `json:"Comment"`
}
