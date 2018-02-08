/* evenlog.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Contains structur for EventLogItem
 */

package nhstruct

type EventLogItem struct {
	Log                string `xml:"log"`
	Source             string `xml:"source"`
	Index              int    `xml:"index"`
	EID                string `xml:"EID"`
	Type               string `xml:"type"`
	GenTime            string `xml:"genTime"`
	WriteTime          string `xml:"writeTime"`
	Machine            string `xml:"machine"`
	ExecutionProcessId int    `xml:"ExecutionProcessId"`
	ExecutionThreadId  int    `xml:"ExecutionThreadId"`
	Message            string `xml:"message"`
	MessageDetail      interface{}
	Category           string `xml:"category"`
	IsGoodEntry        string
	IsWhitelisted      bool
	NHScore            int
	Tag                string
	NhComment          NHComment `json:"Comment"`
}
