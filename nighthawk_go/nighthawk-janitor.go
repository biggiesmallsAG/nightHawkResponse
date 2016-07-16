

package main


import (
	"path/filepath"
	"os"
	"fmt"
	"time"
)

const MAX_DIRSIZE = 2097152
const TIMEDIFF_URG = 3600
const TIMEDIFF_STD = 86400


func main() {
	var media_dir string = "/opt/nighthawk/var/media"
	var workspace_dir string = "/opt/nighthawk/var/workspace"

	mediaSize := SizeOfDir(media_dir)
	if mediaSize >= MAX_DIRSIZE {
		err := CleanUpDir(media_dir, TIMEDIFF_URG)
		if err != nil {
			fmt.Println("Error cleaning ", media_dir)
		}
	} else {
		err := CleanUpDir(media_dir, TIMEDIFF_STD)
		if err != nil {
			fmt.Println("Error cleaning ", media_dir)
		}
	}

	workspaceSize := SizeOfDir(workspace_dir)
	if workspaceSize >= MAX_DIRSIZE {
		err := CleanUpDir(workspace_dir, TIMEDIFF_URG)
		if err != nil {
			fmt.Println("Error cleaning ", workspace_dir)
		}
	} else {
		err := CleanUpDir(workspace_dir, TIMEDIFF_STD)
		if err != nil {
			fmt.Println("Error cleaning ", workspace_dir)
		}
	}
	
}


func SizeOfDir(dirPath string) int64 {
	var dirSize int64 = 0
	
	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}

	filepath.Walk(dirPath, readSize)
	return dirSize	
}

func CleanUpDir(dirPath string, TimeDiff int64) error {
	curUnixTime := time.Now().Unix()

	fileList,err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return err
	}

	for _,file := range fileList {
		fileinfo, err := os.Stat(file)
		if err != nil {
			return nil
		}
		fileUnixTime := fileinfo.ModTime().Unix()

		if curUnixTime > fileUnixTime && (curUnixTime - fileUnixTime) >= TimeDiff {
			fmt.Println("Deleting file ", file)
			os.RemoveAll(file)
		}
	}
	return nil
}
