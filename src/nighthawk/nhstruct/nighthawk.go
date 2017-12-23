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
	CommonName 			string 	`json:"common_name"`
	EmailAddress		string 	`json:"email_address,omitempty"`
	Organization		string 	`json:"organization,omitempty"`
	OrganizationalUnit	string 	`json:"organizational_unit,omitempty"`
	SerialNumber 		string	`json:"serial_number,omitempty"`
	Locality			string 	`json:"locality,omitempty"`
	State 				string 	`json:"state,omitempty"`
	Country				string 	`json:"country,omitempty"`
	FingerprintSHA1 	string	`json:"fingerprint_sha1,omitempty"`
	FingerprintSHA256	string	`json:"fingerprint_sha256,omitempty"`
}


// StackItem structure contains commonly seen audit entry
type StackItem struct {
	AuditType 				string `json:"audit_type"`
	Name 					string `json:"name,omitempty"`
	Path					string `json:"path,omitempty"`
	Md5sum 					string `json:"md5,omitempty"`
	Arguments				string	`json:"arguments,omitempty"`
	RegPath					string `json:"reg_path,omitempty"`
	PersistenceType			string `json:"persistence_type,omitempty"`
	ServiceDescriptiveName 	string `json:"service_descriptive_name,omitempty"`
	TaskCreator				string `json:"task_creator,omitempty"`
}

// WhitelistItem structure contains entry that are whitelisted
// Elasticsearch Index: nighthawk/whitelist
type WhitelistItem struct {
	AuditType 				string `json:"audit_type"`
	Name 					string `json:"name,omitempty"`
	Path					string `json:"path,omitempty"`
	Md5sum 					string `json:"md5,omitempty"`
	Arguments				string	`json:"arguments,omitempty"`
	RegPath					string `json:"reg_path,omitempty"`
	PersistenceType			string `json:"persistence_type,omitempty"`
	ServiceDescriptiveName 	string `json:"service_descriptive_name,omitempty"`
	TaskCreator				string `json:"task_creator,omitempty"`
}

// BlacklistItem structure
// Elasticsearch index: nighthawk/blacklist
type BlacklistItem struct {
	AuditType 				string `json:"audit_type"`
	Name 					string `json:"name,omitempty"`
	Path					string `json:"path,omitempty"`
	Md5sum 					string `json:"md5,omitempty"`
	Arguments				string	`json:"arguments,omitempty"`
	RegPath					string `json:"reg_path,omitempty"`
	PersistenceType			string `json:"persistence_type,omitempty"`
	ServiceDescriptiveName 	string `json:"service_descriptive_name,omitempty"`
	TaskCreator				string `json:"task_creator,omitempty"`
}



// Alerts Structure 
// Elasticsearch index: nighthawk/alert
type Alert struct {
	CaseName			string	`json:"casename"`
	ComputerName		string 	`json:"computer_name"`
	AuditType 			string 	`json:"audit_type"`
	RecordId			string 	`json:"record_id"`
	AlertName			string 	`json:"alert_name"`
	AlertDescription	string 	`json:"alert_description"`
	MatchCondition		string 	`json:"match_condition"`
}
