package utils

import (
	"os"
	"path"
	"sort"
	"time"

	"github.com/guionardo/todo-cli/pkg/logger"
)

type BackupFiles struct {
	backupFolder   string
	fileNamePrefix string
	maxFiles       int
	validateFile   func(string) error
}

func NewBackupFiles(backupFolder string, fileNamePrefix string, maxFiles int, validateFile func(string) error) (*BackupFiles, error) {
	if stat, err := os.Stat(backupFolder); os.IsNotExist(err) || !stat.IsDir() {
		if err = os.MkdirAll(backupFolder, 0755); err != nil {
			return nil, err
		}
	}
	return &BackupFiles{
		backupFolder:   backupFolder,
		fileNamePrefix: fileNamePrefix,
		maxFiles:       maxFiles,
		validateFile:   validateFile,
	}, nil
}

func (b *BackupFiles) GetBackupFiles() (backupFiles []string, err error) {
	var files []os.DirEntry
	if files, err = os.ReadDir(b.backupFolder); err != nil || len(files) == 0 {
		return
	}

	backupFiles = make([]string, 0, len(files))
	match := false
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, err = path.Match(b.fileNamePrefix+".*", file.Name()); err != nil {
			continue
		} else if match {
			backupFile := path.Join(b.backupFolder, file.Name())
			if b.validateFile != nil && b.validateFile(backupFile) != nil {
				// Assert that the backup file is valid
				continue
			}

			backupFiles = append(backupFiles, backupFile)
		}
	}
	sort.Strings(backupFiles)
	return
}

func (b *BackupFiles) Purge() (err error) {
	var backupFiles []string
	if backupFiles, err = b.GetBackupFiles(); err != nil || b.maxFiles < 1 || len(backupFiles) < b.maxFiles {
		return nil
	}
	removedBackupFiles := backupFiles[:len(backupFiles)-b.maxFiles]
	for _, backupFile := range removedBackupFiles {
		if err = os.Remove(backupFile); err != nil {
			logger.Warnf("Error removing backup file %s: %v", backupFile, err)
			return nil
		}
	}
	return nil
}

func (b *BackupFiles) GetNewBackupFileName() string {
	return path.Join(b.backupFolder, b.fileNamePrefix+"."+time.Now().Format("2006-01-02_15:04:05"))
}

func (b *BackupFiles) GetLastBackupFile() string {
	if backupFiles, err := b.GetBackupFiles(); err != nil || len(backupFiles) < 1 {
		return ""
	} else {
		return backupFiles[len(backupFiles)-1]
	}
}
