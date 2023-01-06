package backup

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/guionardo/todo-cli/pkg/config"
	"github.com/guionardo/todo-cli/pkg/utils"
)

type File struct {
	Source         string
	BackupPath     string
	LastBackups    []string
	LastBackup     string
	LastBackupTime time.Time
}

func CreateBackup(source string, backupPath string, cfg config.BackupConfig) (bkp *File, err error) {
	// Check if source path exists
	if err = utils.DirectoryExists(path.Dir(source)); err != nil {
		return nil, err
	}

	if err = utils.DirectoryExists(backupPath); err != nil {
		err = os.Mkdir(backupPath, 0755)
	}

	if err != nil {
		return nil, err
	}

	bkp = &File{
		Source:     source,
		BackupPath: backupPath,
	}
	bkp.PurgeOldBackups(cfg)

	return
}

func (b *File) DoBackup() (backupFile string, err error) {
	if err = utils.FileExists(b.Source); err != nil {
		return "", fmt.Errorf("source file %s does not exists", b.Source)
	}

	if !b.NeedsBackup() {
		return b.LastBackup, nil
	}
	backupTime := time.Now()
	backupFile = fmt.Sprintf("%s.%s", path.Join(b.BackupPath, path.Base(b.Source)), backupTime.Format("20060102150405"))
	for _, err := os.Stat(backupFile); !os.IsNotExist(err); {
		backupTime = backupTime.Add(time.Second)
		backupFile = fmt.Sprintf("%s.%s", path.Join(b.BackupPath, path.Base(b.Source)), backupTime.Format("20060102150405"))
	}

	content, err := os.ReadFile(b.Source)
	if err != nil {
		return "", fmt.Errorf("error reading config file %s: %v", b.Source, err)
	}
	err = os.WriteFile(backupFile, content, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing backup file %s: %v", backupFile, err)
	}
	return backupFile, nil
}

func (b *File) NeedsBackup() bool {
	if b.ReadBackups() != nil || len(b.LastBackups) == 0 {
		return true
	}
	sourceContent, _ := os.ReadFile(b.Source)
	backupContent, _ := os.ReadFile(b.LastBackup)
	return string(sourceContent) != string(backupContent)
}

func (b *File) ReadBackups() error {
	files, err := filepath.Glob(path.Join(b.BackupPath, path.Base(b.Source)) + ".*")
	if err != nil {
		return fmt.Errorf("error listing backup files: %v", err)
	}
	if len(files) > 0 {
		sort.Strings(files)
		b.LastBackups = files
		b.LastBackup = files[len(files)-1]

		if stat, err := os.Stat(b.LastBackup); err == nil {
			b.LastBackupTime = stat.ModTime()
		}
	} else {
		b.LastBackups = make([]string, 0)
		b.LastBackup = ""
	}
	return nil
}

func (b *File) PurgeOldBackups(cfg config.BackupConfig) error {

	if err := b.ReadBackups(); err != nil {
		return err
	}
	if len(b.LastBackups) > cfg.BackupMaxCount {
		for _, backup := range b.LastBackups[:len(b.LastBackups)-cfg.BackupMaxCount] {
			os.Remove(backup)
		}
	}
	return nil
}

func (b *File) AutoBackup(cfg config.BackupConfig) (done bool, err error) {
	if !cfg.AutoBackup {
		return
	}

	if err = b.ReadBackups(); err != nil {
		return
	}
	if len(b.LastBackups) > 0 && time.Since(b.LastBackupTime) < time.Duration(cfg.AutoBackupIntervalDays)*24*time.Hour {
		return
	}
	_, err = b.DoBackup()
	return err == nil, err
}
