/* nighthawk.nhstruct.task.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Task
 */

package nhstruct

type TaskAction struct {
    ActionType          string
    Path                string `xml:"ExecProgramPath"`
    ExecProgramMd5sum   string
    ExecArguments       string 
    ExecWorkingDirectory string 
    COMClassId          string
    COMData             string
}

type TaskTrigger struct {
    TriggerEnabled      bool
    TriggerBegin        string
    TriggerFrequency    string
    TriggerDelay        string
    TriggerSubscription string
}

type TaskItem struct {
    JobCreated          string `xml:"created,attr"`
    TlnTime             string `json:"TlnTime"`
    Name                string 
    Path                string 
    Arguments           string
    VirtualPath         string
    ExitCode            string
    TaskComment         string `xml:"Comment"`
    CreationDate        string `xml:"CreationDate"`
    Creator             string
    MaxRunTime          string 
    Flag                string
    AccountName         string
    AccountRunLevel     string
    AccountLogonType    string
    MostRecentRunTime   string
    NextRunTime         string
    Status              string
    ActionList          []TaskAction `xml:"ActionList>Action"`
    TriggerList         []TaskTrigger `xml:"TriggerList>Trigger"`
    IsGoodTask          string 
    IsWhitelisted       bool 
    IsBlacklisted       bool
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

