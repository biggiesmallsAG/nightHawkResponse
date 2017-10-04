/* nighthawk.elastic.opcontrol.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Control Output flow
 */

package elastic 

import (
 	"fmt"
 	"os"
 	"path/filepath"
 	"encoding/json"

 	nhlog "nighthawk/log"
 	nhs "nighthawk/nhstruct"
 	nhconfig "nighthawk/config"
)

func ProcessOutput(caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, rlrecords []nhs.RlRecord) {

	OPCONTROL := nhconfig.OpControl()

	//cmsg := fmt.Sprintf("Uploading %s:%s %d items", caseinfo.ComputerName, auditinfo.Generator, len(rlrecords))
	//nhlog.ConsoleMessage("INFO", cmsg, nhconfig.VERBOSE)

	nhlog.LogMessage("ProcessOutput", "INFO", fmt.Sprintf("Uploading %s: %s %d items", caseinfo.ComputerName, auditinfo.Generator, len(rlrecords)))

	if OPCONTROL == nhconfig.OP_CONSOLE_ONLY {
		fmt.Println(rlrecords)
	} else if OPCONTROL == nhconfig.OP_DATASTORE_ONLY {
		ExportToElasticsearch(caseinfo.ComputerName, auditinfo.Generator, rlrecords)
	} else if OPCONTROL == nhconfig.OP_CONSOLE_DATASTORE {
		fmt.Println(rlrecords)
		ExportToElasticsearch(caseinfo.ComputerName, auditinfo.Generator, rlrecords)
		
	} else if OPCONTROL == nhconfig.OP_SPLUNK_FILE {

		nhlog.LogMessage("ProcessOutput", "INFO", fmt.Sprintf("Starting writing %s to file", auditinfo.Generator))
		//opfilename := fmt.Sprintf("%s_%s_%s.json", caseinfo.CaseName, caseinfo.ComputerName, auditinfo.Generator)
		opfilename := fmt.Sprintf("%s.log", caseinfo.ComputerName)
		opfile := filepath.Join(nhconfig.DBDIR, opfilename)

		fd,err := os.OpenFile(opfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
		if err != nil {
			/*
			if VERBOSE_LEVEL == 7 {
				fmt.Println(err.Error())
			}
			nighthawk.ExitOnError("Error writing output file", nighthawk.ERROR_WRITING_OUTPUT_FILE)
			*/
			nhlog.LogMessage("ProcessOutput", "DEBUG", err.Error())

		}
		defer fd.Close()

		for _,record := range rlrecords {
			j,_ := json.Marshal(record)
			fd.WriteString(string(j))
			fd.WriteString("\n")
		}

		//nhlog.ConsoleMessage("INFO", "Complete writing " + auditinfo.Generator +" to file", nhconfig.VERBOSE)
		nhlog.LogMessage("ProcessOutput", "INFO", fmt.Sprintf("Complete writing %s to file", auditinfo.Generator))
		fd.Close()
	} else if OPCONTROL == nhconfig.OP_SPLUNK {
		UploadToSplunk(caseinfo.ComputerName, auditinfo.Generator, rlrecords)
	}
}