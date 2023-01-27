package todo

import (
	"github.com/guionardo/todo-cli/pkg/config"
	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/guionardo/todo-cli/pkg/utils"
)

const BackupFilePrefix = "todo_backup"

func validCollectionBackupFile(backupFile string) error {
	_, err := LoadFromFile(backupFile)
	return err
}

// Backup creates a backup of the collection file
func (c *Collection) Backup(cfg *config.BackupConfig) (err error) {
	backupFiles, err := utils.NewBackupFiles(cfg.BackupFolder, BackupFilePrefix,
		cfg.BackupMaxCount, validCollectionBackupFile)
	if err != nil {
		return err
	}

	if lastBackupFilesIsSameCollection(backupFiles.GetLastBackupFile(), *c) {
		return exceptions.NoNeedBackupError
	}

	newBackupFile := backupFiles.GetNewBackupFileName()
	if err = c.Save(newBackupFile); err != nil {
		return err
	}
	logger.Infof("Backup file created: %s", newBackupFile)

	backupFiles.Purge()

	return
}

func (c *Collection) ListBackup(cfg *config.BackupConfig) ([]string, error) {
	backupFiles, err := utils.NewBackupFiles(cfg.BackupFolder, BackupFilePrefix,
		cfg.BackupMaxCount, validCollectionBackupFile)
	if err != nil {
		return nil, err
	}
	return backupFiles.GetBackupFiles()
}

func lastBackupFilesIsSameCollection(lastBackupFile string, c Collection) bool {
	if len(lastBackupFile) == 0 {
		return false
	}
	lastBackupCollection, err := LoadFromFile(lastBackupFile)
	if err != nil {
		return false
	}
	return c.Equal(&lastBackupCollection)
}
