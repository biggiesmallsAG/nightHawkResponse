/* nhstruct::common.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Contains common structures
 */

package nhstruct

type AuditType struct {
	Generator        string `xml:"generator,attr"`
	GeneratorVersion string `xml:"generatorVersion,attr"`
}

type RlAuditType struct {
	Generator        string `xml:"generator,attr"`
	GeneratorVersion string `xml:"generatorVersion,attr"`
}

type CaseInformation struct {
	CaseName     string
	CaseDate     string
	CaseAnalyst  string
	ComputerName string
}

type RlRecord struct {
	ComputerName string
	CaseInfo     CaseInformation
	AuditType    AuditType
	Record       interface{}
}

type NHComment struct {
	Date    string
	Analyst string
	Comment string
}
