/*
 *@package  nightHawk
 *@file     opcontrol.go
 *@author   roshan maskey <roshanmaskey@gmail.com>
 *
 *@description  Control Output flow
 */

package nightHawk

 import (
    "fmt"
    //"io/ioutil"
    "os"
    "path/filepath"
 )


func ProcessOutput(computername string, auditname string, data []byte) {
    if OPCONTROL == OP_CONSOLE_ONLY {

        fmt.Println(string(data))
    } else if OPCONTROL == OP_DATASTORE_ONLY {  

        ExportToElastic(computername, auditname, []byte(data))

    } else if OPCONTROL == OP_CONSOLE_DATASTORE {
        
        fmt.Println(string(data))
        ExportToElastic(computername, auditname, []byte(data))

    } else if OPCONTROL == OP_WRITE_FILE {
        ConsoleMessage("INFO", "Starting writing to file", VERBOSE)
        opfile := fmt.Sprintf("%s_%s.txt", computername, auditname)
        opfilename := filepath.Join(TMP, opfile)

        fd,err := os.OpenFile(opfilename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
        if err != nil {
            panic(err.Error())
        }
        fd.WriteString(string(data))
        ConsoleMessage("INFO", "Complete writing to file", VERBOSE)
        defer fd.Close()
    }
}