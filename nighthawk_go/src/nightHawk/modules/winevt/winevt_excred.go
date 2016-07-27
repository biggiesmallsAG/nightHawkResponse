/*
 *  
 * This file parses Windows EventID 4648
 *
 */

package winevt

import (
    "strings"
    //"fmt"
)


func (el *EventExplictCred)ParseEventMessage(message string) {
    message = strings.Replace(message, "\t", "", -1)
    data := strings.Split(message, "\r\n")
    dataLen := len(data)

    //subOffset, accOffset := GetAccountOffsets(data)
    fields := []string{"Subject", "Account Whose", "Target Server", "Process Information", "Network Information"}
    indexMap := GetSectionIndex(fields, data)
    //fmt.Println(indexMap)

    // Pasrsing Subject Account
    var acc AccountInfo 
    startIndex := indexMap["Subject"] 
    stopIndex := startIndex + 6 
    acc.ReadAccountInfo(startIndex, stopIndex, data)
    el.Subject = acc 

    //Parsing Account Whose Credential was used
    //startIndex = accOffset 
    startIndex = indexMap["Account Whose"]
    stopIndex = startIndex + 4
    acc.ReadAccountInfo(startIndex, stopIndex, data)
    el.Account = acc 

    var procinfo ProcessInformation
    procinfo.ReadProcessInfo(indexMap["Process Information"], indexMap["Network Information"], data)
    el.ProcessInfo = procinfo 

    el.NetworkInfo.ReadNetworkInfo(indexMap["Network Information"], dataLen, data)

    // Parsing remaining 
    startIndex = indexMap["Account Whose"] + 4 
    stopIndex = dataLen

    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2) 
        if len(keyval) == 2 {
            switch keyval[0] {
            case "Target Server Name":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                el.TargetServerName = value 

            case "Additional Information":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                el.AdditionalInfo = value 

            case "Process ID":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                el.ProcessId = value 

            case "Network Address":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                el.NetworkAddress = value 

            case "Port":
                var value string = ""
                if keyval[1] != "-" {
                    value = keyval[1]
                }
                el.Port = value 
            } // _end_switch_
        }
    }
}

