package backup

import (
	"os"
	"path"
	"testing"
)

func TestBackup(t *testing.T) {
	t.Run("Backup", func(t *testing.T) {
		t.Skip()
		//  Create data file
		dataFile := path.Join(t.TempDir(), "data.txt")
		if err := os.WriteFile(dataFile, []byte("Hello World"), 0644); err != nil {
			t.Errorf("Error writing data file: %v", err)
			return
		}
		t.Logf("Data file: %s", dataFile)
		// Create backup dir
		backupDir := path.Join(t.TempDir(), "backup")

		backup, err := CreateBackup(dataFile, backupDir, BackupConfig{})
		if err != nil {
			t.Errorf("Error creating backup: %v", err)
			return
		}

		//  Create backup file
		backupFile, err := backup.DoBackup()

		if err != nil {
			t.Errorf("Error creating backup file: %v", err)
			return
		}

		// Tries to do another backup
		newBackupFile, err := backup.DoBackup()
		if err != nil {
			t.Errorf("Error creating backup file: %v", err)
			return
		}
		if newBackupFile != backupFile {
			t.Errorf("Backup file changed after second backup")
		}

		// change data file
		if err := os.WriteFile(dataFile, []byte("Hello World 2"), 0644); err != nil {
			t.Errorf("Error writing data file: %v", err)
			return
		}

		newBackupFile, err = backup.DoBackup()
		if err != nil {
			t.Errorf("Error creating backup file: %v", err)
			return
		}
		if newBackupFile == backupFile {
			t.Errorf("Backup file did not change after data file change")
		}

	})
}
