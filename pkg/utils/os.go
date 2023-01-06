package utils

import (
	"fmt"
	"os"
)

// FileExists returns true if the file exists and is not a directory
// TODO: Substituir todos os os.Stat por esta funçãos
func FileExists(fileName string) error {
	if stat, err := os.Stat(fileName); err == nil {
		if stat.IsDir() {
			return fmt.Errorf("File %s is a directory", fileName)
		}
		return nil
	} else {
		return err
	}
}

func DirectoryExists(dirName string) error {
	if stat, err := os.Stat(dirName); err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("Directory %s is a file", dirName)
		}
		return nil
	} else {
		return err
	}
}
