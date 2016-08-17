package winevt

import (
    "strings"
    "strconv"

)


func (el *EventNewAccount)ParseEventMessage(message string) {
    message = strings.Replace(message, "\t", "", -1)
    data := strings.Split(message, "\r\n")
    dataLen := len(data)

    
    fields := []string{"Subject", "Logon Information", "New Logon", "Process Information", "Network Information", "Detailed Authentication Information"}
    indexMap := GetSectionIndex(fields, data)

    // Initializing Index variables
    startIndex := 0
    stopIndex := 0

    // Parsing subject section
    startIndex = indexMap["Subject"]
    if indexMap["Logon Information"] > 0 {
        stopIndex = indexMap["Logon Information"]
    } else {
        stopIndex = startIndex + 5  
    }

    var acc AccountInfo 
    acc.ReadAccountInfo(startIndex, stopIndex, data)
    el.Subject = acc 

    // Parsing Logon Information section
    // This section can have multiple attributes, however
    // this section is between "Subject:"" and "New Logon:"
    if indexMap["Logon Information"] > 0 {
        startIndex = indexMap["Logon Information"]
    } else {
        startIndex = stopIndex  // Addresss where Subject Information Ended             
    }
    stopIndex = indexMap["New Logon"]
    var logoninfo LogonInformation 
    logoninfo.ReadLogonInfo(startIndex, stopIndex, data)
    el.LogonInfo = logoninfo 

    /// Parsing New Logon section
    startIndex = indexMap["New Logon"]
    stopIndex = indexMap["Process Information"]
    acc.ReadAccountInfo(startIndex, stopIndex, data)
    el.NewAccount = acc 

    /// Parsing Process Information
    startIndex = indexMap["Process Information"]
    stopIndex = indexMap["Network Information"]
    var procinfo ProcessInformation
    procinfo.ReadProcessInfo(startIndex, stopIndex, data)
    el.ProcessInfo = procinfo 

    /// Parsing Network Information 
    startIndex = indexMap["Network Information"]
    stopIndex = indexMap["Detailed Authentication Information"]
    var netinfo NetworkInformation
    netinfo.ReadNetworkInfo(startIndex, stopIndex, data)
    el.NetworkInfo = netinfo 


    /// Parsing Detailed Authentication Information 
    startIndex = indexMap["Detailed Authentication Information"]
    stopIndex = dataLen 

    for i := startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)
        if len(keyval) == 2 {
            switch keyval[0] {
            case "Logon Process": 
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                el.LogonProcess = value 

            case "Authentication Package":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                el.AuthenticationPackage = value

            case "Transited Services":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                el.TransitedServices = value

            case "Package Name (NTLM only)":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                } 
                el.PackageName = value

            case "Key Length":
                el.KeyLength,_ = strconv.Atoi(keyval[1])
            }
        }
    }

}


func (ev *EventLogonFailed)ParseEventMessage(message string) {
    message = strings.Replace(message, "\t", "", -1)
    data := strings.Split(message, "\r\n")
    dataLen := len(data)

    fields := []string{"Subject", "Logon Type", "Account For", "Failure Information", "Process Information", "Network Information", "Detailed Authentication Information"}
    indexMap := GetSectionIndex(fields, data)

    // Initialize index variables 
    startIndex := 0
    stopIndex := 0

    /// Section Subject
    startIndex = indexMap["Subject"]
    stopIndex = indexMap["Logon Type"]
    var acc AccountInfo 
    acc.ReadAccountInfo(startIndex, stopIndex, data)
    ev.Subject = acc

    /// Logon Type 
    startIndex = indexMap["Logon Type"]
    stopIndex = indexMap["Account For"]
    var logoninfo LogonInformation
    logoninfo.ReadLogonInfo(startIndex, stopIndex, data)
    ev.LogonInfo = logoninfo 

    /// Target Account
    startIndex = indexMap["Account For"]
    stopIndex = indexMap["Failure Information"]
    acc.ReadAccountInfo(startIndex, stopIndex, data)
    ev.Account = acc 

    // Failure Information
    startIndex = indexMap["Failure Information"]
    stopIndex = indexMap["Process Information"]
    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)
        if len(keyval) == 2 {
            switch keyval[0] {
            case "Failure Reason":
                ev.FailureReason = strings.TrimSpace(keyval[1])

            case "Status":
                ev.FailureStatus = strings.TrimSpace(keyval[1])

            case "Sub Status":
                ev.FailureSubStatus = strings.TrimSpace(keyval[1])

            }
        }
    }

    // Process Information 
    startIndex = indexMap["Process Information"]
    stopIndex = indexMap["Network Information"]
    var procinfo ProcessInformation
    procinfo.ReadProcessInfo(startIndex, stopIndex, data)
    ev.ProcessInfo = procinfo 

    // Network Information
    startIndex = indexMap["Network Information"]
    stopIndex = indexMap["Detailed Authentication Information"]
    var netinfo NetworkInformation
    netinfo.ReadNetworkInfo(startIndex, stopIndex, data)
    ev.NetworkInfo = netinfo 

    // Detailed Authentication Information 
    startIndex = indexMap["Detailed Authentication Information"]
    stopIndex = startIndex + 6 // Generally not more than 5 attributes
    if stopIndex > dataLen {
        stopIndex = dataLen 
    }

    for i:= startIndex; i < stopIndex; i++ {
        keyval := strings.SplitN(data[i], ":", 2)
        if len(keyval) == 2 {
            switch keyval[0] {
            case "Logon Process":
                ev.LogonProcess = strings.TrimSpace(keyval[1])

            case "Authentication Package":
                ev.AuthenticationPackage = strings.TrimSpace(keyval[1])

            case "Transited Services":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                ev.TransitedServices = value 

            case "Package Name (NTML Only)":
                var value string = ""
                if keyval[1] != "-" {
                    value = strings.TrimSpace(keyval[1])
                }
                ev.PackageName = value 

            case "Key Length":
                ev.KeyLength,_ = strconv.Atoi(keyval[1])
            }
        }
    }
}

/// This function will check and return offset of Subject and NewLogon Section
func GetAccountOffsets(data []string) (subOffset int, accOffset int) {
    subOffset = 0
    accOffset = 0

    for i := range data {
        if strings.Contains(data[i], "Subject:") {
            subOffset = i 
        }

        if strings.Contains(data[i], "New Logon:") || strings.Contains(data[i], "Account Whose Credentials Were Used:") {
            accOffset = i
            break
        }
    }
    return subOffset, accOffset
}

