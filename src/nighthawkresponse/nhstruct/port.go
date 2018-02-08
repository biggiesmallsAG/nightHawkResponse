/* nighthawk.nhstruct.port.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Ports
 */

package nhstruct

type PortItem struct {
	JobCreated    string `xml:"created,attr"`
	Pid           int    `xml:"pid"`
	Process       string `xml:"process"`
	Path          string `xml:"path"`
	State         string `xml:"state"`
	LocalIP       string `xml:"localIP"`
	RemoteIP      string `xml:"remoteIP"`
	LocalPort     int    `xml:"localPort"`
	RemotePort    int    `xml:"remotePort"`
	Protocol      string `xml:"protocol"`
	IsGoodPort    string
	IsWhitelisted bool
	NHScore       int
	Tag           string
	NhComment     NHComment `json:"Comment"`
}
