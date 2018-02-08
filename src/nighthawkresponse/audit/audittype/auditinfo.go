/* auditinfo.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Contains AuditInfo structure
 */

package audittype

import (
	nhs "nighthawkresponse/nhstruct"
	"os"
	"regexp"
)

const (
	PTGenerator        = "w32processes-tree"
	PTGeneratorVersion = "0.0.1"
)

func GetAuditInfoFromFile(auditfile string) (nhs.AuditType, error) {
	var auditinfo nhs.AuditType

	fd, err := os.Open(auditfile)
	if err != nil {
		return auditinfo, err
	}
	defer fd.Close()

	buffer := make([]byte, 500)
	_, err = fd.Read(buffer)
	if err != nil {
		return auditinfo, err
	}

	s := string(buffer)
	re := regexp.MustCompile("generator=\"(.*)\" generatorVersion=\"([0-9.]+)\" ")
	match := re.FindStringSubmatch(s)

	auditinfo.Generator = match[1]
	auditinfo.GeneratorVersion = match[2]
	return auditinfo, nil
}
