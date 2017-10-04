/* nighthawk.nhstruct.process_memory.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Process-Memory
 */

package nhstruct

type Handle struct {
    Index               int
    AccessMask          int
    ObjectAddress       string
    HandleCount         int 
    PointerCount        int
    Type                string
    Name                string
    IsGoodHandle        string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type MemorySection struct {
    Protection          string
    RegionStart         int
    RegionSize          int
    Mapped              bool
    RawFlags            string
    IsGoodSection       string
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type ProcessItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Pid                 int     `xml:"pid"`
    ParentPid           int     `xml:"parentpid"`
    Path                string  `xml:"path"`
    Name                string  `xml:"name"`
    Arguments           string  `xml:"arguments"`
    Username            string  `xml:"Username"`
    SecurityID          string  `xml:"SecurityID"`
    SecurityType        string  `xml:"SecurityType"`
    StartTime           string  `xml:"startTime"`
    KernelTime          string  `xml:"kernelTime"`
    UserTime            string  `xml:"userTime"`
    HandleList          []Handle `xml:"HandleList>Handle"`
    //SectionList       []MemorySection `xml:"SectionList>MemorySection"`
    IsGoodProcess       string
    IsWhitelisted       bool
    IsBlacklisted       bool
    NHScore             int
    Tag                 string 
    NhComment           NHComment `json:"Comment"`
}



