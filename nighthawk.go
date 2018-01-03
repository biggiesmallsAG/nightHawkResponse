/*
 * nighthawk.go
 * authors:
 *	- roshan maskey <roshanmaskey@gmail.com>,
 * 	- Daniel Eden <danieleden@gmail.com>
 *
 * description: nighthawk main binary
 *
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"nighthawk"
	"nighthawk/audit"
	"nighthawk/audit/audittype"
	nhc "nighthawk/common"
	nhconfig "nighthawk/config"
	nhlog "nighthawk/log"
	nhs "nighthawk/nhstruct"
	"nighthawk/rabbitmq"
	nhutil "nighthawk/util"
	"nighthawkapi/api/handlers/upload"
)

type RuntimeOptions struct {
	CaseName     	string
	CaseDate     	string
	CaseAnalyst  	string
	ComputerName 	string
	ConfigFile   	string
	TriageFile   	string
	Version      	bool
	Verbose      	bool
	VerboseLevel 	int 
	Daemon       	bool
	PidFile      	string
	OutputType		int 		// bitmask 
	Standalone		bool		// run standalone nighthawk without elasticsearch
}

func main() {
	runtime.GOMAXPROCS(nhconfig.MaxProcs())
	//nhconfig.SetEnvironment()

	// Setting commandline parser
	var runopt RuntimeOptions

	flag.StringVar(&runopt.CaseName, "casename", "", "Case name for collected triage. Default: System generated")
	flag.StringVar(&runopt.CaseDate, "casedate", "", "Case date for collected triage. Default: Today")
	flag.StringVar(&runopt.CaseAnalyst, "analyst", "", "Case analyst working on collected triage")
	flag.StringVar(&runopt.ComputerName, "computername", "", "Computername of collected triage")
	flag.StringVar(&runopt.ConfigFile, "config", nhconfig.CONFIG_FILE, "nightHawk Response configuration file.")
	flag.StringVar(&runopt.TriageFile, "file", "", "Collected triage file")
	flag.BoolVar(&runopt.Version, "version", false, "Display Version Information")
	flag.BoolVar(&runopt.Verbose, "verbose", false, "Show verbose message")
	flag.BoolVar(&runopt.Daemon, "daemon", false, "Daemonize the process")
	flag.StringVar(&runopt.PidFile, "pid", "", "PID file")
	flag.IntVar(&runopt.VerboseLevel, "verbose-level",4,"Log verbose level. Default is ERROR")
	flag.IntVar(&runopt.OutputType,"output-type", 2, "Default is Elasticsearch")
	flag.BoolVar(&runopt.Standalone,"standalone", false, "Run nighthawk standalone without elasticsearch. Default is NO")

	flag.Parse()

	if runopt.Version {
		ShowVersion()
		os.Exit(0)
	}

	nhconfig.SetConfigFile(runopt.ConfigFile)
	err := nhconfig.LoadNighthawkConfig()
	if err != nil {
		nhlog.LogMessage("main", "ERROR", err.Error())
		os.Exit(nhc.ERROR_CONFIG_FILE_READ)
	}

	if runopt.Verbose {
		nhconfig.SetVerbose()
	}

	// Setting Log VerboseLevel and Output-Type (OpControl)
	nhconfig.SetVerboseLevel(runopt.VerboseLevel)
	nhconfig.SetOutputType(runopt.OutputType)
	
	if runopt.Standalone {
		nhconfig.SetStandalone()
	}


	/*
	 * Daemonize the process as nighthawk_worker process
	 */
	if runopt.Daemon && runopt.PidFile != "" {
		//nighthawk.ConsoleMessage("INFO", "Running in Daemon mode", nhconfig.VERBOSE)
		nhlog.LogMessage("main", "INFO", "Running in Daemon mode")

		// RabbtiMQ is always used in worker mode
		rconfig := rabbitmq.LoadRabbitMQConfig(rabbitmq.RABBITMQ_CONFIG_FILE)
		conn := rabbitmq.Connect(rconfig.Server)
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			//panic(err.Error())
			nhlog.LogMessage("main", "ERROR", err.Error())
			os.Exit(nhc.ERROR_CHANNEL_CONNECT)
		}

		// Init log exchange
		_ = rabbitmq.RabbitExchangeDeclare(ch, rconfig.LogTopicExchange)
		// Init work exchange
		_ = rabbitmq.RabbitExchangeDeclare(ch, rconfig.WorkTopicExchange)
		// Init worker queue
		_ = rabbitmq.RabbitQueueDeclare(ch, rconfig.Worker)
		// Bind to Log exchange
		_ = rabbitmq.RabbitBindQueue(ch, rconfig.Logger)
		// Bind to Work exchange
		_ = rabbitmq.RabbitBindQueue(ch, rconfig.Worker)

		messages := rabbitmq.RabbitQueueConsumer(ch, rconfig.Worker)

		forever := make(chan bool)
		go func() {
			for data := range messages {

				var jobmsg upload.Job
				json.Unmarshal(data.Body, &jobmsg)

				jobmsg.InProg = true
				//nhutil.LoggerPublish(conn, "INFO", fmt.Sprintf("New job received: %s, Total Audit files: %d", jobmsg.UID, len(jobmsg.Audits)))
				//nhutil.LoggerPublish(conn, "JOB", &jobmsg)

				nhlog.LogMessage("main", "INFO", fmt.Sprintf("New job received: %s, Total Audit files: %d", jobmsg.UID, len(jobmsg.Audits)))
				nhlog.CreateJobMessage(&jobmsg)

				if jobmsg.CaseID == "" {
					jobmsg.CaseID = nhc.GenerateCaseName()
				}

				var caseinfo = nhs.CaseInformation{jobmsg.CaseID, jobmsg.TS, jobmsg.UserID, jobmsg.ComputerName}

				for _, _audit := range jobmsg.Audits {
					sourcetype := nhutil.SourceDataFileType(_audit)
					if sourcetype == nhutil.MOD_XML {
						ParseSingleAuditFile(caseinfo, _audit)
					} else if sourcetype == nhutil.MOD_ZIP {
						ParseTriageFile(caseinfo, _audit)
					} else {
						//nhutil.LoggerPublish(conn, "ERROR", "Unsupported triage file.")
						nhlog.LogMessage("main", "ERROR", "Unsupported triage file")
					}
				}

				jobmsg.InProg = false
				jobmsg.Complete = true

				//nhutil.LoggerPublish(conn, "INFO", fmt.Sprintf("Finished Job UID: %s", jobmsg.UID))
				//nhutil.LoggerPublish(conn, "JOB", &jobmsg)
				nhlog.LogMessage("main", "INFO", fmt.Sprintf("Finished Job UID: %s", jobmsg.UID))
				nhlog.CreateJobMessage(&jobmsg)
			}
		}()

		//nighthawk.ConsoleMessage("INFO", "Waiting for Jobs.. Hit Ctrl-C to exit..", nhconfig.VERBOSE)
		nhlog.LogMessage("main", "INFO", "Waiting for Jobs... Hit Ctrl-C to exit")
		<-forever

		return
	}

	/**************************************************************
	 *                  NON_DAEMON_MODE
	 * Section of code processes individual triage without running
	 * in daemon mode without receiving instructions from RabbitMQ
	 * server.
	 *
	 * Single triage file can be processed using following command
	 * line syntax:
	 *       go run nighthawk.go -N "CASE-001" -case-date "2017-01-01"  -C "computername" -a "analyst_1" -f triagefile.zip
	 *
	 *
	 **************************************************************/

	if runopt.CaseName == "" {
		runopt.CaseName = nhc.GenerateCaseName()
	}

	if runopt.CaseDate == "" {
		runopt.CaseDate = fmt.Sprintf("%s", time.Now().UTC().Format(nhc.Layout))
	}

	if runopt.TriageFile == "" {
		//ExitOnError("Triage file required", nighthawk.ERROR_NO_TRIAGE_FILE)
		nhlog.LogMessage("main", "ERROR", "Triage file required")
		os.Exit(nhc.ERROR_NO_TRIAGE_FILE)
	}

	// Setting absolute path for input file
	absTriageFile, err := filepath.Abs(runopt.TriageFile)
	if err == nil {
		runopt.TriageFile = absTriageFile
	}

	//if runopt.Redis {
	//    nhconfig.REDIS_PUB = true
	//}

	if runopt.CaseName == runopt.ComputerName {
		//nighthawk.RedisPublish("ERROR", runopt.CaseName + " is same as ComputerName. This is not allowed.", nhconfig.REDIS_PUB)
		//ExitOnError("CaseName and ComputerName can not be same", nighthawk.ERROR_SAME_CASE_AND_COMPUTERNAME)
		nhlog.LogMessage("main", "ERROR", "CaseName and ComputerName cannot be same")
		os.Exit(nhc.ERROR_SAME_CASE_AND_COMPUTERNAME)
	}

	// __end_of_commandline_parsing

	var caseinfo = nhs.CaseInformation{CaseName: runopt.CaseName, CaseDate: runopt.CaseDate, CaseAnalyst: runopt.CaseAnalyst, ComputerName: runopt.ComputerName}

	sourcetype := nhutil.SourceDataFileType(runopt.TriageFile)
	if sourcetype == nhutil.MOD_XML {
		ParseSingleAuditFile(caseinfo, runopt.TriageFile)
	} else if sourcetype == nhutil.MOD_ZIP {
		ParseTriageFile(caseinfo, runopt.TriageFile)
	} else {
		nhlog.LogMessage("main", "ERROR", "Unsupported triage file")
		os.Exit(nhc.ERROR_UNSUPPORTED_TRIAGE_FILE)

	}

} // __end_of_main__

func ParseSingleAuditFile(caseinfo nhs.CaseInformation, auditfile string) {
	if caseinfo.ComputerName == "" {
		//ExitOnError("Computer name must be supplied. Use -C switch", nighthawk.ERROR_AUDIT_COMPUTERNAME_REQUIRED)
		nhlog.LogMessage("ParseSingleAuditFile", "ERROR", "ComputerName must be supplied. User -C switch")
		os.Exit(nhc.ERROR_AUDIT_COMPUTERNAME_REQUIRED)
	}

	nhlog.LogMessage("ParseSingleAuditFile", "INFO", fmt.Sprintf("Parsing single audit file for %s", caseinfo.ComputerName))
	_, filename := filepath.Split(auditfile)

	auditinfo, err := audittype.GetAuditInfoFromFile(auditfile)
	if err != nil {
		nhlog.LogMessage("ParseSingleAuditFile", "ERROR", "Error encountered while extracting audit type information")
		nhlog.LogMessage("ParseSingleAuditFile", "ERROR", err.Error())
		os.Exit(nhc.ERROR_AUDITTYPE_INFO_PARSE)
	}

	nhlog.LogMessage("ParseSingleAuditFile", "INFO", fmt.Sprintf("Parsing %s::%s from audit file %s", caseinfo.ComputerName, auditinfo.Generator, filename))
	audit.ParseAuditFile(caseinfo, auditinfo, auditfile)

}

// Process triage package
// caseinfo: CaseInformation structure defined in nhstruct
// triagefile: Triage package file (user supplied). It can be either HX mans, or Redline audit output
func ParseTriageFile(caseinfo nhs.CaseInformation, triagefile string) {
	// Enable RabbitMQ if configuration file is available
	// If RabbitMQ config file is not eanble, ignore messaging
	// specially useful as standalone binary
	if nhutil.FileExists(rabbitmq.RABBITMQ_CONFIG_FILE) {
		rconfig := rabbitmq.LoadRabbitMQConfig(rabbitmq.RABBITMQ_CONFIG_FILE)
		conn := rabbitmq.Connect(rconfig.Server)
		defer conn.Close()
	}

	nhlog.LogMessage("ParseTriageFile", "INFO", fmt.Sprintf("Processing triage package %s", triagefile))

	targetDir := CreateSessionDirectory(triagefile)
	nhlog.LogMessage("ParseTriageFile", "INFO", fmt.Sprintf("Session directory %s created for %s", targetDir, triagefile))

	// Fix for multi-level folder for Redline audits
	if !IsRedlineAuditDirectory(targetDir) {
		nhlog.LogMessage("ParseTriageFile", "INFO", fmt.Sprintf("%s is not Redline audit directory", targetDir))

		dirList, _ := filepath.Glob(filepath.Join(targetDir, "*"))

		for _, d := range dirList {
			if IsRedlineAuditDirectory(d) {
				targetDir = d
				nhlog.LogMessage("ParseTriageFile", "INFO", fmt.Sprintf("Session directory updated to %s", targetDir))

				break
			}
		}
	}

	// Checking for manifest file in Redline audit folder
	// Manifiest file is automatically generated by HX, Redline
	manifest, err := audit.GetAuditManifestFile(targetDir)
	if err != nil {
		nhlog.LogMessage("ParseTriageFile", "INFO", "Audit manifest file not found in session directory")

		manifest = audit.GenerateAuditManifestFile(targetDir)
	}

	var rlman audit.RlManifest
	rlman.ParseAuditManifest(filepath.Join(targetDir, manifest))
	stAudits := rlman.Payloads2(targetDir)

	// _rm> 2017-08-25
	// Only extract computername if it has not been provided
	// in commandline argument
	if caseinfo.ComputerName == "" {
		// _rm> 2017-08-25
		// TODO: Implement function to extract computername from xml file.
		// Either HX3.5 or HXTool2 collection method does not provide metadata.json file
		// code needs to be updated to extract computer name from sysinfo OR w32system
		computername := rlman.SysInfo.SystemInfo.Machine
		caseinfo.ComputerName = computername
	}

	if caseinfo.ComputerName == "" {
		nhlog.LogMessage("ParseTriageFile", "ERROR", "Failed to get ComputerName from audit")
		os.Exit(nhc.ERROR_READING_COMPUTERNAME)
	}

	// Check if Casename is same as computername
	if caseinfo.CaseName == caseinfo.ComputerName {
		nhlog.LogMessage("ParseTriageFile", "ERROR", "CaseName and ComputerName can not be same")
		os.Exit(nhc.ERROR_SAME_CASE_AND_COMPUTERNAME)
	}

	nhlog.LogMessage("ParseTriageFile", "INFO", fmt.Sprintf("Creating new case %s", caseinfo.CaseName))
	nhlog.LogMessage("ParseTriageFile", "INFO", fmt.Sprintf("Processing Redline audits for %s", caseinfo.ComputerName))

	var wg sync.WaitGroup

	for _, stAudit := range stAudits {

		auditfile := filepath.Join(targetDir, stAudit.AuditFile)

		auditinfo, err := audittype.GetAuditInfoFromFile(auditfile)
		if err != nil {
			nhlog.LogMessage("ParseTriageFile", "ERROR", "Error encountered while extracting audit type information")
			os.Exit(nhc.ERROR_AUDITTYPE_INFO_PARSE)
		}

		nhlog.LogMessage("ParseTriageFile", "INFO", fmt.Sprintf("Parsing %s::%s from audit file %s", caseinfo.ComputerName, auditinfo.Generator, stAudit.AuditFile))

		wg.Add(1)
		GoParseAudit(&wg, caseinfo, auditinfo, auditfile)
	}
	wg.Wait()
	nhlog.LogMessage("ParseTriageFile", "DEBUG", "Waiting for GOROUTINES to complete")
	os.RemoveAll(targetDir)
	nhlog.LogMessage("ParseTriageFile", "DEBUG", fmt.Sprintf("Session directory %s removed", targetDir))
}

func GoParseAudit(wg *sync.WaitGroup, caseinfo nhs.CaseInformation, auditinfo nhs.AuditType, auditfile string) {
	defer wg.Done()
	audit.ParseAuditFile(caseinfo, auditinfo, auditfile)
}

func ShowVersion() {
	fmt.Println("\t nightHawk Response")
	fmt.Printf("\t Version: %s, Build: %s\n", nighthawk.Version, nighthawk.Build)
	fmt.Println("\t By Roshan Maskey and Daniel Eden")
}

/// This function will create and extract supplied archive audit file
/// and returns the full path of the file
func CreateSessionDirectory(filename string) string {
	sessionDir := nhc.NewSessionDir(nhconfig.SessionDirSize())
	targetDir := filepath.Join(nhconfig.WORKSPACE, sessionDir)

	os.MkdirAll(targetDir, 0755)
	err := nhutil.Unzip(filename, targetDir)

	if err != nil {
		nhlog.LogMessage("CreateSessionDirectory", "ERROR", "Error encountered extracting Redline archive")
		nhlog.LogMessage("CreateSessionDirectory", "ERROR", err.Error())
		os.Exit(nhc.ERROR_EXTRACTING_REDLINE_ARCHIVE)
	}

	return targetDir
}

/// Check if the given folder contains Redline Audits
func IsRedlineAuditDirectory(dirPath string) bool {
	nhlog.LogMessage("IsRedlineAuditDirectory", "DEBUG", fmt.Sprintf("Checking if %s is Redline directory", dirPath))

	fList, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		nhlog.LogMessage("IsRedlineAuditDirectory", "DEBUG", fmt.Sprintf("Error listing files in %s", dirPath))
		nhlog.LogMessage("IsRedlineAuditDirectory", "ERROR", err.Error())
	}

	if len(fList) <= 5 {
		return false
	}

	// Read all the files header. w32system file must be available.
	// Continue checking until w32system audit is found.
	for _, f := range fList {
		// Session directory may contain subfolder. Ignore the subfolders
		fh, err := os.Open(f)
		if err != nil {
			nhlog.LogMessage("IsRedlineAuditDirectory", "DEBUG", err.Error())
		}
		defer fh.Close()

		fi, _ := fh.Stat()

		if fi.Mode().IsRegular() {
			// Read the audit file
			buf := make([]byte, 500)
			_, err := fh.Read(buf)
			if err != nil {
				panic(err.Error())
			}

			if strings.Contains(string(buf), "audit_manifest") || strings.Contains(string(buf), "w32system") {
				return true
			}
		}
	}
	return false
}
