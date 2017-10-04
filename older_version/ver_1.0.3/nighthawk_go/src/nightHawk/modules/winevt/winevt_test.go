package winevt

import (
    "fmt"
    "testing"
    
    "io/ioutil"
    "encoding/xml"
    "encoding/json"
    
    //"encoding/hex"
    "strconv"

)

type EventLogItem struct {
    Log             string  `xml:"log"`
    Source          string  `xml:"source"`
    Index           int     `xml:"index"`
    EID             int     `xml:"EID"`
    Type            string  `xml:"type"`
    GenTime         string  `xml:"genTime"`   
    WriteTime       string  `xml:"writeTime"`
    Machine         string  `xml:"machine"`
    ExecutionProcessId int  `xml:"ExecutionProcessId"`
    ExecutionThreadId int   `xml:"ExecutionThreadId"`
    Message         string  `xml:"message"`
    Category        string  `xml:"category"`
    MessageDetail   interface{}
    NHScore         int 
    Tag             string 
    NhComment       string `json:"Comment"`
}

type RlEventLog struct {
    ComputerName        string 
    //CaseInfo            CaseInformation
    //AuditType           RlAuditType
    EventList           []EventLogItem  `xml:"EventLogItem"`
}

func (eli *EventLogItem)PrintEventLogItem() {
    fmt.Println("Log: ", eli.Log)
    fmt.Println("Source: ", eli.Source)
    fmt.Println("Index: ", eli.Index)
    fmt.Println("EID: ", eli.EID)
    fmt.Println("Type: ", eli.Type)
    fmt.Println("Gen Time: ", eli.GenTime)
    fmt.Println("Write Time: ", eli.WriteTime)
    fmt.Println("Machine: ", eli.Machine)
    fmt.Println("Execution Process ID: ", eli.ExecutionProcessId)
    fmt.Println("Execution Thread ID: ", eli.ExecutionThreadId)
    //fmt.Println("Message: ", eli.Message)
    fmt.Println("Category: ", eli.Category)
    fmt.Println("MessageDetail: ", eli.MessageDetail)
    fmt.Println("NHScore: ", eli.NHScore)
    fmt.Println("Tag: ", eli.Tag)
    fmt.Println("NhComment: ", eli.NhComment)
}


func TestHelloWorld(t *testing.T) {
    var pidHex string = "0020"
    //pid,_ := hex.DecodeString(pidHex)
    pid,_ := strconv.ParseInt(pidHex, 16, 64)
    fmt.Println(pidHex, pid)
}


func TestEventLogParsing(t *testing.T) {
    filename := "/Users/rmaskey/Downloads/RedlineAudit/w32eventlogs.5nvCp4kJTQadD3abegDuEJ"
    xmlData,_ := ioutil.ReadFile(filename)

    var rl RlEventLog
    xml.Unmarshal(xmlData, &rl)

    for _,el := range rl.EventList {
        el.MessageDetail = ProcessEventItem(el.Log, el.EID, el.Message)
        //if el.EID == 4624  || el.EID == 4648 || el.EID == 4688 {
        if el.EID == 4697 {
            //el.PrintEventLogItem()    
            j,_ := json.MarshalIndent(el,""," ")
            fmt.Println(string(j))
        }
        
    }
}