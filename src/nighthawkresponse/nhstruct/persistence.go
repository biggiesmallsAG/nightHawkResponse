/* nighthawk.nhstruct.persistence.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for persistence
 */

package nhstruct

type PersistenceItem struct {
	JobCreated        string `xml:"created,attr"`
	TlnTime           string
	PersistenceType   string
	RegPath           string
	RegOwner          string
	RegModified       string
	Path              string `xml:"FilePath"`
	FileOwner         string
	FileCreated       string
	FileModified      string
	FileAccessed      string
	FileChanged       string
	Md5sum            string       `xml:"md5sum"`
	File              FileItem     `xml:"FileItem"`
	Registry          RegistryItem `xml:"RegistryItem"`
	StackPath         string
	IsGoodPersistence string
	IsWhitelisted     bool
	IsBlacklisted     bool
	NHScore           int
	Tag               string
	NhComment         NHComment `json:"Comment"`
}
