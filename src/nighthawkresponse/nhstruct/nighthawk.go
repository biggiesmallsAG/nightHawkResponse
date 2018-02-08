/*
 * nighthawk.nhstruct.nighthawk.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * This file contains structures to enrich collected data
 * analyzer modules (nighthawk.analyze) will utilize information
 * stored by these structures
 */
package nhstruct

// IssuingCA structure contains DigitalCertificate
// attributes for Issuing CA
type IssuingCA struct {
	CommonName         string `json:"CommonName"`
	EmailAddress       string `json:"EmailAddress,omitempty"`
	Organization       string `json:"Organization,omitempty"`
	OrganizationalUnit string `json:"OrganizationalUnit,omitempty"`
	SerialNumber       string `json:"SerialNumber,omitempty"`
	Locality           string `json:"Locality,omitempty"`
	State              string `json:"State,omitempty"`
	Country            string `json:"Country,omitempty"`
	FingerprintSHA1    string `json:"FingerprintSha1,omitempty"`
	FingerprintSHA256  string `json:"FingerprintSha256,omitempty"`
}

// StackItem structure contains commonly seen audit entry
type StackItem struct {
	AuditType              string `json:"AuditType"`
	Name                   string `json:"Name,omitempty"`
	Path                   string `json:"Path,omitempty"`
	Md5sum                 string `json:"Md5,omitempty"`
	Arguments              string `json:"Arguments,omitempty"`
	RegPath                string `json:"RegPath,omitempty"`
	PersistenceType        string `json:"PersistenceType,omitempty"`
	ServiceDescriptiveName string `json:"ServiceDescriptiveName,omitempty"`
	TaskCreator            string `json:"TaskCreator,omitempty"`
}

// WhitelistItem structure contains entry that are whitelisted
// Elasticsearch Index: nighthawk/whitelist
type WhitelistItem struct {
	AuditType              string `json:"AuditType"`
	Name                   string `json:"Name,omitempty"`
	Path                   string `json:"Path,omitempty"`
	Md5sum                 string `json:"Md5,omitempty"`
	Arguments              string `json:"Arguments,omitempty"`
	RegPath                string `json:"RegPath,omitempty"`
	PersistenceType        string `json:"PersistenceType,omitempty"`
	ServiceDescriptiveName string `json:"ServiceDescriptiveName,omitempty"`
	TaskCreator            string `json:"TaskCreator,omitempty"`
}

// BlacklistItem structure
// Elasticsearch index: nighthawk/blacklist
type BlacklistItem struct {
	AuditType              string `json:"AuditType"`
	Name                   string `json:"Name,omitempty"`
	Path                   string `json:"Path,omitempty"`
	Md5sum                 string `json:"Md5,omitempty"`
	Arguments              string `json:"Arguments,omitempty"`
	RegPath                string `json:"RegPath,omitempty"`
	PersistenceType        string `json:"PersistenceType,omitempty"`
	ServiceDescriptiveName string `json:"ServiceDescriptiveName,omitempty"`
	TaskCreator            string `json:"TaskCreator,omitempty"`
}

// Alert strucutre
type Alert struct {
	CaseName         string
	ComputerName     string
	AuditType        string
	DocID            string
	AlertName        string
	AlertDescription string
	MatchCondition   string
}
