/* nighthawk.util.srcfile
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Determine sourcefile provided by user: xml file or zip file
 */

package util

import (
	//"fmt"
	//"nighthawk/config"
	"os"

)

const (
	MOD_UNC = 0
	MOD_ZIP = 1
	MOD_XML = 2
)

var BUF_ZIP = []byte{80, 75, 3, 4}
var BUF_XML = []byte{60, 63, 120, 109}

func SourceDataFileType(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 4)
	file.Read(buf)

	if CompareSlice(buf, BUF_ZIP) {
		return MOD_ZIP
	}

	if CompareSlice(buf, BUF_XML) {
		return MOD_XML
	}

	return MOD_UNC

}

func CompareSlice(X, Y []byte) bool {
	for i := range X {
		if X[i] != Y[i] {
			return false
		}
	}
	return true
}
