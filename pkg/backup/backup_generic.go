package backup

import (
	"fmt"
	"os"
	"path"
	"sort"
	"time"

	"github.com/guionardo/todo-cli/pkg/config"
	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/guionardo/todo-cli/pkg/interfaces"
	"github.com/guionardo/todo-cli/pkg/logger"
)

type Backup struct {
	prefix   string
	folder   string
	maxFiles int
	instance interfaces.Backupable
}

func NewBackup[K interfaces.Backupable](instance K, cfg config.BackupConfig) *Backup {
	return &Backup{
		prefix:   instance.BackupPrefix(),
		folder:   cfg.BackupFolder,
		maxFiles: cfg.BackupMaxCount,
		instance: instance,
	}
}

func (b *Backup) assertBackupFolder() error {
	stat, err := os.Stat(b.folder)
	if err == nil {
		if stat.IsDir() {
			return nil
		}
		return exceptions.NewException(fmt.Errorf("Backup folder is not a directory: %s", b.folder))
	}
	return os.MkdirAll(b.folder, 0755)
}

func (b *Backup) GetNewBackupFileName() string {
	return path.Join(b.folder, b.prefix+"."+time.Now().Format("20060102_150405"))
}

func (b *Backup) ValidateFile(fileName string) error {
	if b.instance == nil {
		return nil
	}
	if content, err := os.ReadFile(fileName); err != nil {
		return err
	} else {
		return b.instance.Parse(content)
	}
}

func (b *Backup) GetBackupFiles() (backupFiles []string, err error) {
	var files []os.DirEntry
	if files, err = os.ReadDir(b.folder); err != nil || len(files) == 0 {
		return
	}

	backupFiles = make([]string, 0, len(files))
	match := false
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, err = path.Match(b.prefix+".*", file.Name()); err != nil {
			continue
		} else if match {
			backupFile := path.Join(b.folder, file.Name())

			if b.ValidateFile(backupFile) != nil {
				// Assert that the backup file is valid
				continue
			}

			backupFiles = append(backupFiles, backupFile)
		}
	}
	sort.Strings(backupFiles)
	return
}

func (b *Backup) Run() error {
	if err := b.assertBackupFolder(); err != nil {
		return err
	}
	if !b.NeedsBackup() {
		return nil
	}
	newFileName := b.GetNewBackupFileName()
	if err := b.instance.Save(newFileName); err == nil {
		return b.Purge()
	} else {
		return err
	}
}

func (b *Backup) Purge() (err error) {
	if err = b.assertBackupFolder(); err != nil {
		return err
	}
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

func (b *Backup) GetLastBackupFile() string {
	if backupFiles, err := b.GetBackupFiles(); err != nil || len(backupFiles) < 1 {
		return ""
	} else {
		return backupFiles[len(backupFiles)-1]
	}
}

func (b *Backup) NeedsBackup() bool {
	lastBackup := b.GetLastBackupFile()
	if lastBackup == "" {
		return true
	}
	if content, err := os.ReadFile(lastBackup); err != nil {
		return true
	} else {
		if other, err := b.instance.ParseNew(content); err != nil {
			return true
		} else {
			return b.instance.Equal(other)
		}
	}
}
