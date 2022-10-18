package backup

import (
	"fmt"
	"path"
)

type BackupConfig struct {
	BackupNeededWarningDays int    `yaml:"backup_needed_warning_days"`
	BackupMaxCount          int    `yaml:"backup_max_count"`
	AutoBackup              bool   `yaml:"auto_backup"`
	AutoBackupIntervalDays  int    `yaml:"auto_backup_interval_days"`
	BackupFolder            string `yaml:"backup_folder"`
}

func (bc BackupConfig) String() string {
	return fmt.Sprintf("BackupNeededWarningDays: %d\nBackupMaxCount: %d\nAutoBackup: %v\nAutoBackupIntervalDays: %d\nBackupFolder: %s\n",
		bc.BackupNeededWarningDays, bc.BackupMaxCount, bc.AutoBackup, bc.AutoBackupIntervalDays, bc.BackupFolder)
}

func GetDefaultBackupConfig(dataFolder string) BackupConfig {
	return BackupConfig{
		BackupNeededWarningDays: 7,
		BackupMaxCount:          10,
		AutoBackup:              true,
		AutoBackupIntervalDays:  7,
		BackupFolder:            path.Join(dataFolder, "backup"),
	}
}
