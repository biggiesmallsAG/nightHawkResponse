/*
 * Base Package for common function
 *
 * SendMessageToConsole(sender, level, message)
 * ExitOnError(sender, level, message, errcode)
 */

package common

import (
	"fmt"
	"os"
	"time"
)




func ConsoleMessage(sender string, level string, message string, verbose bool) {
	if verbose {
		fmt.Printf("%s - %s - %s - %s\n", time.Now().UTC().Format(Layout), sender, level, message)
	}
}


func ExitOnError(sender string, level string, message string, errcode int) {
	ConsoleMessage(sender, "ERROR", message, true)
	os.Exit(errcode)
}


