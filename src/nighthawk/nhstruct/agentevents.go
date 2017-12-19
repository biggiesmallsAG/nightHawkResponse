/* nighthawk.nhstruct.agentevents.go
 *
 * DataStructure for Agent Events
 */
package nhstruct

type ItemDetail struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type AgentEventItem struct {
	Timestamp string       `xml:"timestamp"`
	EventType string       `xml:"eventType"`
	Details   []ItemDetail `xml:"details>detail"`
}

type AddressNotificationEvent struct {
	Timestamp string
	EventType string
	Address   string
}

type FileWriteEvent struct {
	Timestamp    string `xml:"timestamp"`
	EventType    string `xml:"eventType"`
	FullPath     string `xml:"fullPath"`
	FilePath     string `xml:"filePath"`
	FileName     string `xml:"fileName"`
	FileExt      string `xml:"fileExtension"`
	Drive        string `xml:"drive"`
	DevicePath   string `xml:"devicePath"`
	ProcessID    int    `xml:"pid"`
	ProcessName  string `xml:"process"`
	WriteCount   int    `xml:"writes"`
	BytesWritten int    `xml:"numBytesSeenWritten"`
	DataOffset   int    `xml:"lowestFileOffsetSeen"`
	Data         string `xml:"dataAtLowestOffset"`
	TextData     string `xml:"textAtLowerOffset"`
	IsClosed     bool   `xml:"closed"`
	FileSize     int    `xml:"size"`
	MD5          string `xml:"md5"`
}

type ImageLoadEvent struct {
	Timestamp   string `xml:"timestamp"`
	EventType   string `xml:"eventType"`
	FullPath    string `xml:"fullPath"`
	FilePath    string `xml:"filePath"`
	Drive       string `xml:"drive"`
	FileName    string `xml:"fileName"`
	FileExt     string `xml:"fileExtension"`
	ProcessID   int    `xml:"pid"`
	ProcessName string `xml:"process"`
}

type NetworkEvent struct {
	Timestamp   string `xml:"timestamp"`
	EventType   string `xml:"eventType"`
	RemoteIP    string `xml:"remoteIP"`
	RemotePort  int    `xml:"remotePort"`
	LocalIP     string `xml:"localIP"`
	LocalPort   int    `xml:"localPort"`
	Protocol    string `xml:"protocol"`
	ProcessID   int    `xml:"pid"`
	ProcessName string `xml:"process"`
}

type DnsLookupEvent struct {
	Timestamp   string
	EventType   string
	Hostname    string
	ProcessID   int
	ProcessName string
}

type RegKeyEvent struct {
	Timestamp        string
	EventType        string
	Hive             string
	KeyPath          string
	Path             string
	NotificationType int // 1: key change; 2: ;3: key created; 4: key deleted
	ProcessID        int
	ProcessName      string
	ValueName        string
	ValueType        string
	Value            string
	Text             string
}
