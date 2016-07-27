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
    SysInfo             RlSystemInfo
    Type                string 
    Version             string 
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
        _,filename := filepath.Split(file)

        // Checking keyword manifest in filename and file extension .json
        if strings.Contains(filename, "manifest") && strings.HasSuffix(filename,"json") {
            manifest_file = filename
            break
        } else if strings.HasSuffix(filename, "json") {
            // This is loose checking for manifest file. Just checking for json file if keyword is not found
            manifest_file = filename
        }
    }

    /// Something to do in future. Validate that json file is correct auditmanifest file by unmarshaling
    return manifest_file, nil
}


