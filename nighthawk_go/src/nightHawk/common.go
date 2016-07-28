/*
 *@package  nightHawk
 *@file     common.go
 *@author   roshan maskey <roshanmaskey@gmail.com>
 *
 *@description  nightHawk Response common
 */


 package nightHawk

 import (
    "fmt"
    "strings"
    "os"
    "time"
    "regexp"
 )

 const (
    RL_AGENTSTATE       = "stateagentinspector"         // This is HX triage specific not collected as a part of normal Redline Triage
    RL_PERSISTENCE      = "w32scripting-persistence"
    RL_SERVICES         = "w32services"
    RL_PORTS            = "w32ports"
    RL_USERACCOUNTS     = "w32useraccounts"
    RL_TASKS            = "w32tasks"
    RL_PROCESSMEMORY    = "w32processes-memory"
    RL_PREFETCH         = "w32prefetch"
    RL_REGRAW           = "w32registryraw"
    RL_SYSTEM           = "w32system" 
    RL_DISKS            = "w32disks"
    RL_VOLUMES          = "w32volumes"
    RL_URLHISTORY       = "urlhistory"
    RL_FILEDLHISTORY    = "filedownloadhistory"
    RL_NETWORKDNS       = "w32network-dns"
    RL_NETWORKROUTE     = "w32network-route"
    RL_NETWORKARP       = "w32network-arp"
    RL_APIFILES         = "w32apifiles"
    RL_RAWFILES         = "w32rawfiles"
    RL_HIVELIST         = "w32hivelist"
    RL_SYSTEMRESTORE    = "w32systemrestore"
    RL_KERNELHOOK       = "w32kernel-hookdetection"
    RL_EVENTLOGS        = "w32eventlogs"
 )

 func ShowVersion () {
    fmt.Printf("\tnightHawk Response ver %s\n", VERSION)
    fmt.Printf(">> Triage processor for Mandiant Redline\n")
    fmt.Printf(">> by Roshan Maskey and Daniel Eden\n")
 }


 func ShowAuditGenerators() {
    fmt.Println("nightHawk Response - Development in progress....")
 }


 func ConsoleMessage(level string, message string, verbose bool) {
    if verbose {
        fmt.Printf("%s - %s - %s\n", time.Now().UTC().Format(Layout), level, message)
    }
 }

func ExitOnError(errmsg string, errcode int) {
    ConsoleMessage("ERROR", errmsg, true)
    os.Exit(errcode)
}

 func GenerateCaseName() string {
    part_a := strings.ToUpper(NewSessionDir(5))
    part_b := strings.ToUpper(NewSessionDir(3))
    casename := fmt.Sprintf("%s-%s", part_a, part_b)
    return casename
 }


 func SourceDataFileType(filename string) DataFileType {
    if strings.HasSuffix(filename, ".xml") {
        return MOD_XML
    } 

    if strings.HasSuffix(filename, ".zip") {
        return MOD_ZIP
    }

    if strings.HasSuffix(filename, ".mans") {
        return MOD_MANS
    }

    /// Checking if the supplied path is directory
    fd,err := os.Open(filename)
    if err != nil {
        ConsoleMessage("ERROR", "Failed to read "+filename, true)
        os.Exit(ERROR_READING_TRIAGE_FILE)
    }
    defer fd.Close()
    finfo,_ := fd.Stat()

    if finfo.Mode().IsDir() {
        return MOD_REDDIR
    }

    // Default if nothing matches assume it is single XML audit file
    return MOD_XML
 }

 func UrlToHostname(Url string) (string, string) {
    var Hostname string = ""
    var Domain string = ""

    re,_ := regexp.Compile("(http|https|ftp)://[^/]+")

    if !re.MatchString(Url) {
        return "",""
    }

    r0,_ := regexp.Compile(":\\d+")
    baseUrl := r0.ReplaceAllString(Url, "")


    baseUrl = re.FindString(baseUrl)

    prefixList := []string{"http", "https", "ftp"}

    for _, urlProto := range prefixList {
            if strings.HasPrefix(baseUrl, urlProto) {
                    baseUrl = strings.Replace(baseUrl, urlProto+"://", "", 2)
            }
    }

    ipre,_ := regexp.Compile("(\\d+\\.){3}\\d+")
    if ipre.MatchString(baseUrl) {
            return baseUrl, baseUrl
    }

    urlPart := strings.Split(baseUrl, ".")
    numPart := len(urlPart)

    if numPart <= 2 {
        return baseUrl, baseUrl
    } else if numPart == 4{
        Hostname = baseUrl
        Domain = urlPart[1]+"."+urlPart[2]+"."+urlPart[3]
    }  else {
        if len(urlPart[numPart-1]) == 2 {
                Hostname = baseUrl
                Domain = urlPart[numPart-3]+"."+urlPart[numPart-2]+"."+urlPart[numPart-1]
        } else {
                Hostname = baseUrl
                Domain = urlPart[numPart-2]+"."+urlPart[numPart-1]
        }
    }

    return Hostname, Domain
 }



