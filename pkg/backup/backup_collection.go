package backup

import (
	"path"
	"time"
)

type BackupCollection struct {
	SourcePath string                 // Path to the source folder
	Backups    map[string]*BackupFile // Map of files and their backups
}

func NewBackupCollection(sourcePath string, config BackupConfig, files ...string) (bc *BackupCollection, err error) {
	bc = &BackupCollection{
		SourcePath: sourcePath,
		Backups:    make(map[string]*BackupFile),
	}
	for _, file := range files {
		bkp, err := CreateBackup(path.Join(sourcePath, file), config.BackupFolder, config)
		if err != nil {
			return nil, err
		}
		bc.Backups[file] = bkp
	}
	return
}

func (bc *BackupCollection) DoBackup() (err error) {
	for _, bkp := range bc.Backups {
		if _, err = bkp.DoBackup(); err != nil {
			return
		}
	}
	return
}

func (bc *BackupCollection) AutoBackup(config BackupConfig) (done bool, err error) {
	if !config.AutoBackup {
		return
	}
	for _, b := range bc.Backups {
		if err = b.ReadBackups(); err != nil {
			return false, err
		}
		if len(b.LastBackups) > 0 && time.Since(b.LastBackupTime) < time.Duration(config.AutoBackupIntervalDays)*24*time.Hour {
			continue
		}
		if _, err = b.DoBackup(); err != nil {
			return false, err
		}
	}
	return true, nil
}
