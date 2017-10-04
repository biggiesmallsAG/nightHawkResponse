/* nighthawk.nhstruct.apifile.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for APIFiles
 */

package nhstruct

type VirusTotal struct {
	CheckedDate string
	VTScore     int
}

type PEChecksum struct {
	PEFileRaw     int
	PEFileAPI     int
	PEComputedAPI int
}

type DigitalSignature struct {
	SignatureExists    bool
	SignatureVerified  bool
	Description        string
	CertificateSubject string
	CertificateIssuer  string
}

type PEInfo struct {
	Type             string
	Subsystem        string
	BaseAddress      string
	PETimeStamp      string
	PeChecksum       PEChecksum       `xml:"PEChecksum"`
	DigitalSignature DigitalSignature `xml:"DigitalSignature"`
}

type FileItem struct {
	JobCreated     string `xml:"created,attr"`
	TlnTime        string `json:"TlnTime"`
	DevicePath     string
	Path           string `xml:"FullPath"`
	Drive          string
	FilePath       string
	FileName       string
	FileExtension  string
	SizeInBytes    int
	Created        string
	Modified       string
	Accessed       string
	Changed        string
	FileAttributes string
	Username       string
	SecurityID     string
	SecurityType   string
	Md5sum         string
	PeInfo         PEInfo `xml:"PEInfo"`
	IsGoodHash     string
	IsWhitelisted  bool
	IsBlacklisted  bool
	NHScore        int
	VTResults      []VirusTotal
	Tag            string
	NhComment      NHComment `json:"Comment"`
}

type RawFileItem struct {
	JobCreated       string `xml:"created,attr"`
	TlnTime          string `json:"TlnTime"`
	DevicePath       string
	Path             string `xml:"FullPath"`
	Drive            string
	FilePath         string
	FileName         string
	FileExtension    string
	SizeInBytes      int
	Created          string
	Modified         string
	Accessed         string
	Changed          string
	FilenameCreated  string
	FilenameModified string
	FilenameAccessed string
	FilenameChanged  string
	FileAttributes   string
	Username         string
	SecurityID       string
	SecurityType     string
	Md5sum           string
	PeInfo           PEInfo `xml:"PEInfo"`
	IsGoodHash       string
	IsWhitelisted    bool
	IsBlacklisted    bool
	NHScore          int
	VTResults        []VirusTotal
	Tag              string
	NhComment        NHComment `json:"Comment"`
}
