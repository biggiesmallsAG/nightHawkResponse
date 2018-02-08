/* nighthawkresponse.elastic.opcontrol.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Control Output flow
 */

package elastic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	nhconfig "nighthawkresponse/config"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
)

// Output Control constants
const (
	OP_CTRL_CONSOLE = 0x0001
	OP_CTRL_ELASTIC = 0x0002
	OP_CTRL_FILE    = 0x0004
	OP_CTRL_SPLUNK  = 0x0008
	OP_CTRL_CSV     = 0x0010
)

func ProcessOutput(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, rlrecords []nhs.RlRecord) {

	OPCONTROL := nhconfig.OpControl()

	//cmsg := fmt.Sprintf("Uploading %s:%s %d items", caseinfo.ComputerName, auditinfo.Generator, len(rlrecords))
	//nhlog.ConsoleMessage("INFO", cmsg, nhconfig.VERBOSE)

	nhlog.LogMessage("ProcessOutput", "INFO", fmt.Sprintf("Uploading %s: %s %d items", caseinfo.ComputerName, auditinfo.Generator, len(rlrecords)))

	if OPCONTROL&OP_CTRL_CONSOLE == OP_CTRL_CONSOLE {
		fmt.Println(rlrecords)
	}

	if OPCONTROL&OP_CTRL_ELASTIC == OP_CTRL_ELASTIC {
		ExportToElasticsearch(caseinfo.ComputerName, auditinfo.Generator, rlrecords)
	}

	if OPCONTROL&OP_CTRL_FILE == OP_CTRL_FILE {

		nhlog.LogMessage("ProcessOutput", "INFO", fmt.Sprintf("Starting writing %s to file", auditinfo.Generator))
		//opfilename := fmt.Sprintf("%s_%s_%s.json", caseinfo.CaseName, caseinfo.ComputerName, auditinfo.Generator)
		opfilename := fmt.Sprintf("%s.log", caseinfo.ComputerName)
		opfile := filepath.Join(nhconfig.DBDIR, opfilename)

		fd, err := os.OpenFile(opfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			/*
				if VERBOSE_LEVEL == 7 {
					fmt.Println(err.Error())
				}
				nighthawkresponse.ExitOnError("Error writing output file", nighthawkresponse.ERROR_WRITING_OUTPUT_FILE)
			*/
			nhlog.LogMessage("ProcessOutput", "DEBUG", err.Error())

		}
		defer fd.Close()

		for _, record := range rlrecords {
			j, _ := json.Marshal(record)
			fd.WriteString(string(j))
			fd.WriteString("\n")
		}

		//nhlog.ConsoleMessage("INFO", "Complete writing " + auditinfo.Generator +" to file", nhconfig.VERBOSE)
		nhlog.LogMessage("ProcessOutput", "INFO", fmt.Sprintf("Complete writing %s to file", auditinfo.Generator))
		fd.Close()
	}

	if OPCONTROL&OP_CTRL_SPLUNK == OP_CTRL_SPLUNK {
		UploadToSplunk(caseinfo.ComputerName, auditinfo.Generator, rlrecords)
	}

	if OPCONTROL&OP_CTRL_CSV == OP_CTRL_CSV {
		GenerateAuditCsv(caseinfo.ComputerName, auditinfo.Generator, rlrecords)
	}
}
