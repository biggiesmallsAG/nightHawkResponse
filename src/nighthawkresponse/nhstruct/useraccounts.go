/* nighthawk.nhstruct.accounts.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for User accounts
 */

package nhstruct

type UserItem struct {
	JobCreated       string `xml:"created,attr"`
	TlnTime          string `json:"TlnTime"`
	Username         string `xml:"Username"`
	SecurityID       string `xml:"SecurityID"`
	SecurityType     string `xml:"SecurityType"`
	Fullname         string `xml:"fullname"`
	Description      string `xml:"description"`
	HomeDirectory    string `xml:"homedirectory"`
	ScriptPath       string `xml:"scriptpath"`
	LastLogin        string `xml:"lastlogin"`
	Disabled         bool   `xml:"disabled"`
	LockedOut        bool   `xml:"lockedout"`
	PasswordRequired bool   `xml:"passwordrequired"`
	UserPasswordAge  string `xml:"userpasswordage"`
	IsGoodEntry      string
	IsWhitelisted    bool
	NHScore          int
	Tag              string
	NhComment        NHComment `json:"Comment"`
}
