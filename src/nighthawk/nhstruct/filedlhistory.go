/* nighthawk.nhsruct.downloadhistory.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for FileDownloadHistory
 */

package nhstruct

type FileDownloadHistoryItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Profile             string
    BrowserName         string
    BrowserVersion      string
    Username            string
    DownloadType        string 
    Filename            string `xml:"FileName"`
    SourceUrl           string `xml:"SourceURL"`
    UrlHostname         string
    UrlDomain           string 
    TargetDirectory     string `xml:"TargetDirectory"`
    FullHttpHeader      string
    LastModifiedDate    string
    BytesDownloaded     int 
    MaxBytes            int
    CacheFlags          string
    CacheHitCount       int
    LastCheckedDate     string
    IsGoodFile          string
    IsWhitelisted       bool
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

