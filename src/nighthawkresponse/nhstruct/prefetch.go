/* nighthawk.nhstruct.prefetch.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Prefetch 
 */

package nhstruct

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
    IsWhitelisted       bool
    NHScore             int
    Tag                 string 
    NhComment           NHComment `json:"Comment"`
}

