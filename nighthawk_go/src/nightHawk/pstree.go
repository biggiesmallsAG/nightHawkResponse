/*
 *@package  nightHawk
 *@file     pstree.go
 *@author   roshan maskey <roshanmaskey@gmail.com>
 *
 *@description  nightHawk Response ProcessTree
 */


package nightHawk

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "math/rand"
)

const (
    PTGenerator = "w32processes-tree"
    PTGeneratorVersion = "0.0.1"
)

 type ProcessTree struct {
    ParentPid       int `json:"parent"`
    Pid             int `json:"pid"`
    Path            string `json:"path"`
    Name            string `json:"name"`
    Arguments       string `json:"arguments"`
    StartTime       string `json:"startTime"`
    Username        string `json:"Username"`
    Children        []*ProcessTree `json:"_children"`
 }


// This function creates process tree
 func CreateProcessTree(caseinfo CaseInformation, computername string, filename string) []byte {
    xmlData,_ := ioutil.ReadFile(filename)
    var rl RlProcessMemory
    var auditinfo RlAuditType

    rl.ParseAuditData(computername, caseinfo, auditinfo, xmlData)   

    var root ProcessTree
    root.ParentPid = 0
    root.Pid = 0
    root.Name = "Start"

    var PpidList []int 
    var PidList []int
    var ptList []ProcessTree

    for _,ml := range rl.ProcessList {
        
            var pt ProcessTree

            pt.ParentPid = ml.ParentPid
            pt.Pid = ml.Pid
            pt.Path = ml.Path
            pt.Name = ml.Name
            pt.Arguments = ml.Arguments
            pt.StartTime = ml.StartTime
            pt.Username = ml.Username

            ptList = append(ptList, pt)
    }

    ptList = QSortPpid(ptList)

    for i,pt := range ptList {
        PpidList = append(PpidList, ptList[i].ParentPid)
        PidList = append(PidList, ptList[i].Pid)
    
            FlagParentExists := ParentExists(ptList[i], &root)

            if !FlagParentExists {
                var ParentName string = fmt.Sprintf("%d", pt.ParentPid)
                var node = ProcessTree {ParentPid: 0, Pid: pt.ParentPid, Name: ParentName}
                node.AddChildNode(&ptList[i])

                root.AddChildNode(&node)
            }   
    }   

    // Create RlRecord object compatible with ElasticSearch mapping
    var ptAuditType = RlAuditType{Generator: PTGenerator, GeneratorVersion: PTGeneratorVersion}
    var ptRecord RlRecord
    ptRecord.ComputerName = computername 
    ptRecord.CaseInfo = caseinfo 
    ptRecord.AuditType = ptAuditType
    ptRecord.Record = root

    jsonData,_ := json.Marshal(ptRecord)
    return jsonData
 }

 func ParentExists(pt ProcessTree, node *ProcessTree) bool {
    if node.Pid == pt.ParentPid {
        node.AddChildNode(&pt)
        return true
    } 

    for _,childnode := range node.Children {
        if ParentExists(pt, childnode) {
            return true 
        }
    }
    return false
 }

 func (node *ProcessTree)AddChildNode(child *ProcessTree) {
    node.Children = append(node.Children, child)
 }

 // This function returns parent index number 
 func ParentPidIndex(pid int, name string, ptree []ProcessTree) int {
    for i,p := range ptree {
        if p.Pid == pid {
            return i
        }
    }
    // If nothing found return -1
    return -1
 }


// This quicksort code is based on vderyagin quicksort.go
func QSortPpid(slice []ProcessTree) []ProcessTree {
    length := len(slice)

    if length <= 1 {
        sliceCopy := make([]ProcessTree, length)
        copy(sliceCopy, slice)
        return sliceCopy
    }

    ppid := slice[rand.Intn(length)].ParentPid

    less := make([]ProcessTree, 0, length)
    middle := make([]ProcessTree, 0, length)
    more := make([]ProcessTree, 0, length)

    for _,process := range slice {
        switch {
        case process.ParentPid < ppid:
            less = append(less, process)
        case process.ParentPid == ppid:
            middle = append(middle, process)
        case process.ParentPid > ppid:
            more = append(more, process)
        }
    }

    less, more = QSortPpid(less), QSortPpid(more)
    less = append(less, middle...)
    less = append(less, more...)
    return less
}

