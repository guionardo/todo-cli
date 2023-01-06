package backup

import (
	"fmt"
	"os"
)

// MoveFileToBackup Move file to back up on same folder, with file modification time as suffix
func MoveFileToBackup(filename string) (newFileName string, err error) {	
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return
	} else if stat.IsDir() {
		err = fmt.Errorf("file %s is a directory", filename)
		return
	}

	newFileName = fmt.Sprintf("%s.%s", filename, stat.ModTime().Format("20060102150405"))
	err = os.Rename(filename, newFileName)
	return
}
