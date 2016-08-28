/*
 *@package  nightHawk
 *@file     manifest.go
 *@author   roshan maskey <roshanmaskey@gmail.com>
 *
 *@description  This file contains structure and functions to process Redline audit manifest file
 */

package nightHawk 

import (
    "io/ioutil"
    "encoding/json"
    "encoding/xml"
    "strings"
    "path/filepath"
    "errors"
    "os"
    "regexp"
)

type AuditResult struct {
    Payload             string `json:"payload"`
    PayloadType         string `json:"type"` 
}

type AuditGenerator struct {
    Generator           string `json:"generator"`
    GeneratorVersion    string `json:"generatorVersion"`
    AuditResults        []AuditResult `json:"results"`
}

type RlManifest struct {
    SysInfo             RlSystemInfo `json:"sysinfo"`
    Type                string `json:"type"`
    Version             string `json:"version"`
    Audits              []AuditGenerator `json:"audits"`
}

type RlAudit struct {
    AuditGenerator      string
    AuditFile           string 
}

func (rlman *RlManifest) ParseAuditManifest(filename string) error {
    manifestData, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    json.Unmarshal(manifestData, &rlman)
    return nil 
}



/// This function returns the filename for the give audit generator
/// It does not return full file path.
/// This function ignore any file containing audit issues
func (rlman *RlManifest) Payload(generator string) []string {
    var payload []string

    for _,auditgenerator := range rlman.Audits {
        if auditgenerator.Generator == generator {
            for _,p := range auditgenerator.AuditResults {
                // Only append the files without issue. We are going to ignore issue file
                if !strings.Contains(p.Payload, "issue") {
                    payload = append(payload, p.Payload)
                }
            }
        }
    }
    return payload
}


/// This function returns all the payloads file from manifest
func (rlman *RlManifest) Payloads(session_dir string) []string {
    var payload []string 
    var w32system_file string 

    for _,ag := range rlman.Audits {
        for _,p := range ag.AuditResults {
            if !strings.Contains(p.Payload, "issue") {
                payload = append(payload, p.Payload)

                if ag.Generator == "w32system" {
                    w32system_file = p.Payload
                }
            }
        }
    }
    
    xmlData,_ := ioutil.ReadFile(filepath.Join(session_dir, w32system_file))
    xml.Unmarshal(xmlData, &rlman.SysInfo)

    return payload
}


func (rlman *RlManifest)Payloads2(session_dir string) []RlAudit {
    var rlaudits []RlAudit
    var w32system_file string 

    for _,ag := range rlman.Audits {
        for _,p := range ag.AuditResults {
            if !strings.Contains(p.Payload, "issue") && !strings.Contains(p.PayloadType, "issue"){
                var r = RlAudit{AuditGenerator: ag.Generator, AuditFile: p.Payload}
                rlaudits = append(rlaudits, r)

                if ag.Generator == "w32system" {
                    w32system_file = p.Payload
                }
            }
        }
    }

    xmlData,_ := ioutil.ReadFile(filepath.Join(session_dir, w32system_file))
    xml.Unmarshal(xmlData, &rlman.SysInfo)

    return rlaudits
}


/// Load SystemInformation


/// This function returns manifest file in given session directory
func GetAuditManifestFile(session_dir string) (string,error) {
    filelist,err := filepath.Glob(filepath.Join(session_dir, "*"))
    if err != nil {
        return "",err
    } 
 
    var manifest_file string = ""
    for _,file := range filelist {
        manifest_file = ""
        _,filename := filepath.Split(file)

        lwrFilename := strings.ToLower(filename)
        // Checking keyword manifest in filename and file extension .json
        if strings.Contains(lwrFilename, "manifest") && strings.HasSuffix(lwrFilename,"json") {
            manifest_file = filename
            break
        } else if strings.HasSuffix(lwrFilename, "json") {
            // This is loose checking for manifest file. Just checking for json file if keyword is not found
            manifest_file = filename
        }
    }

    /// Something to do in future. Validate that json file is correct auditmanifest file by unmarshaling
    if len(manifest_file) > 2 {
        return manifest_file, nil    
    }

    return manifest_file, errors.New("Error no manifest file")
}



/// This function generates manifest file. This function is 
/// typically used for audits generated using MIR (Mandiant Intelligence Response)
/// The manifest file created by this function is different that Redline/HX manifest

func GenerateAuditManifestFile(session_dir string) string { 
    var manfilename string = "nh_manifest.json"

    var rlman RlManifest
    rlman.Type = "audit_manifest"
    rlman.Version = "1.0"


    filelist, err := filepath.Glob(filepath.Join(session_dir,"*"))
    if err != nil {panic(err.Error())}

    for _,file := range filelist {
        if IsRegularFile(file) {

            fh, err := os.Open(file)
            if err != nil {panic(err.Error())}
            defer fh.Close()

            buf := make([]byte, 500)
            fh.Read(buf)

            strbuf := string(buf)

            // No processing file containig Redline or MIR issues
            if strings.Contains(strbuf, "issue.xsd") {
                continue   
            }

            /// Begin creating manifest audits
            var ar AuditResult
            var ag AuditGenerator

            ar.Payload = filepath.Base(file)
            ar.PayloadType = "application/xml"
        
            re := regexp.MustCompile("generator=\"(.*)\" generatorVersion=\"([0-9.]+)\" ")
            match := re.FindStringSubmatch(strbuf)

            if len(match) > 2 {
                ag.Generator = match[1]
                ag.GeneratorVersion = match[2]

                // HX audit failed audit contains "FireEye Agent" as generator.
                // Ignoring audits with issue.
                if ag.Generator != "FireEye Agent" {
                    ag.AuditResults = append(ag.AuditResults, ar)
                    rlman.Audits = append(rlman.Audits, ag)    
                }
                    
            }
            

        }
    }

    manJsonData,_ := json.MarshalIndent(&rlman,"", " ")

    ConsoleMessage("INFO", "Generating manifest file " + manfilename, VERBOSE)
    
    err = ioutil.WriteFile(filepath.Join(session_dir,manfilename), manJsonData, 0644)
    if err != nil {
        ConsoleMessage("ERROR", "Error writing " + manfilename + " to session directory " + session_dir, VERBOSE)
        return ""
    }

    return manfilename
}


func IsRegularFile(file string) bool {
    fh, err := os.Open(file)
    if err != nil {
        ConsoleMessage("ERROR", "Error opening file " + file, true)
        return false
    }

    defer fh.Close()

    fi,_ := fh.Stat()

    if fi.Mode().IsRegular() {
        return true
    }
    return false
}

