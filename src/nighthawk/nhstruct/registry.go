/* nighthawk.nhstruct.registry.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Registry Entry
 */

package nhstruct 


type RegistryItem struct {
    JobCreated      string `xml:"created,attr"`
    TlnTime         string `json:"TlnTime"`
    KeyPath         string
    Type            string
    Modified        string
    ValueName       string
    Username        string
    Path            string
    Text            string
    ReportedLengthInBytes   int
    Hive            string
    SecurityID      string
    IsKnownKey      string
    IsWhitelisted   bool
    IsBlacklisted   bool
    NHScore         int
    Tag             string 
    NhComment       NHComment `json:"Comment"`
}


