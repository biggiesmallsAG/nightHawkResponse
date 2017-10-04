/* nhstruct::common.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Contains common structures  
 */

package nhstruct


type AuditType struct {
    Generator           string `xml:"generator,attr"`
    GeneratorVersion    string `xml:"generatorVersion,attr"` 
}


type RlAuditType struct {
    Generator           string `xml:"generator,attr"`
    GeneratorVersion    string `xml:"generatorVersion,attr"` 
}


type CaseInformation struct {
    CaseName            string `json:"case_name"`
    CaseDate            string `json:"case_date"`
    CaseAnalyst         string `json:"case_analyst"`
    ComputerName        string `json:"computer_name"`
}


type RlRecord struct {
    ComputerName        string 
    CaseInfo            CaseInformation
    AuditType           AuditType
    Record              interface{}
}


type NHComment struct {
    Date            string `json:"Date"`
    Analyst         string `json:"Analyst"`
    Comment         string `json:"Comment"`
}
