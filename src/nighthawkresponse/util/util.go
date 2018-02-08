/* nighthawk.util.util
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Contains common utility functions
 */

package util

import (
 	"fmt"
	"os"
	"time"
)

const Layout = "2006-01-02T15:04:05Z"

func FixEmptyTimestamp() string {
	return "1970-01-01T01:01:01Z"
}


func CurrentTimestamp() string {
	return fmt.Sprintf("%s", time.Now().UTC().Format(Layout))
}



func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	} 
	panic(err)
}

