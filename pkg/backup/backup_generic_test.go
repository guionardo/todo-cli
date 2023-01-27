package backup

import (
	"testing"

	"github.com/guionardo/todo-cli/pkg/config"
	"github.com/guionardo/todo-cli/pkg/interfaces"
	"github.com/guionardo/todo-cli/pkg/todo"
)

func TestNewBackup(t *testing.T) {
	cfg := config.BackupConfig{
		BackupFolder:   t.TempDir(),
		BackupMaxCount: 2,
	}

	tests := []struct {
		name      string
		instances []interfaces.Backupable
	}{
		{"LocalConfig", []interfaces.Backupable{
			&config.LocalConfig{
				ToDoListName: "test0",
				Gist:         config.GistConfig{},
				Backup:       config.BackupConfig{},
			},
			&config.LocalConfig{
				ToDoListName: "test1",
				Gist:         config.GistConfig{},
				Backup:       config.BackupConfig{},
			},
		}},
		{"Collection", []interfaces.Backupable{&todo.Collection{}}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, inst := range tt.instances {
				got := NewBackup(inst, cfg)
				if got == nil {
					t.Errorf("NewBackup() = %v, want %v", got, tt.instances)
					return
				}
				if err := got.Run(); err != nil {
					t.Errorf("NewBackup() = %v, want %v", got, tt.instances)
					return
				}
				backupFiles, _ := got.GetBackupFiles()
				if len(backupFiles) != 1 {
					t.Errorf("Backup files() = %v, want %v", backupFiles, []string{})
					return
				}
			}
		})
	}
}
