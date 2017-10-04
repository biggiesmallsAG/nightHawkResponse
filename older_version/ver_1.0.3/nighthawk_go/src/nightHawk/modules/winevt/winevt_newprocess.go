/* 
 * Parse Windows Event ID 4688: A new process has been created
 * 
 * reference: https://www.ultimatewindowssecurity.com/securitylog/encyclopedia/event.aspx?eventID=4688
 * 
 * author:  roshan maskey <0xredskull>
 */

 package winevt

 import (
    //"fmt"
    "strings"
    "strconv"
 )

 func (ev *EventNewProcess)ParseEventMessage(message string) {
    message = strings.Replace(message, "\t", "", -1)
    data := strings.Split(message, "\r\n")
    dataLen := len(data)

    fields := []string{"Subject", "Process Information"}
    indexMap := GetSectionIndex(fields, data)

    // Parsing Windows Event Subject 
    startIndex := indexMap["Subject"]
    stopIndex := startIndex + 8
    var acc AccountInfo 
    acc.ReadAccountInfo(startIndex, stopIndex, data)
    ev.Subject = acc 

    // Parsing for Process Information 
    startIndex = indexMap["Process Information"]
    stopIndex = startIndex + 8

    if stopIndex > dataLen {
        stopIndex = dataLen 
    }

    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)
        if len(keyval) == 2 {
            switch keyval[0] {
            case "New Process ID":
                var value string = "0x00"
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                ev.NewProcessIdHex = value 
                value = strings.Replace(value, "0x", "", -1)
                pid,_ := strconv.ParseInt(value, 16, 64)
                ev.NewProcessId = int(pid)

            case "New Process Name":
                ev.NewProcessName = strings.TrimSpace(keyval[1])

            case "Token Elevation Type":
                ev.TokenElevationType = keyval[1]

                if strings.Contains(keyval[1], "(1)") {
                    ev.TokenElevationDesc = "UAC is disabled. Full administrator privilege token assinged"
                }

                if strings.Contains(keyval[1], "(2)") {
                    ev.TokenElevationDesc = "UAC is enabled. Elevated token assinged. Program started using Run as administrator"
                }

                if strings.Contains(keyval[1], "(3)") {
                    ev.TokenElevationDesc = "UAC is enabled. Limited tokent assinged. Program was not started using Run as administrator"
                }

            case "Mandatory Label":
                ev.MandatoryLabel = keyval[1]

            case "Creator Process ID":
                var value string = "0x00"
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                ev.CreatorProcessIdHex = value 
                value = strings.Replace(value, "0x", "", -1)
                pid,_ := strconv.ParseInt(value, 16, 64)
                ev.CreatorProcessId = int(pid)

            case "Creator Process Name":
                ev.CreatorProcessName = keyval[1]

            case "Process Command Line":
                ev.ProcessCommandline = keyval[1]
            } // _end_switch_
        }
    }

 }


