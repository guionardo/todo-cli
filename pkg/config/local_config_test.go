package config

import (
	"os"
	"path"
	"testing"
)

func TestLocalConfig_Save(t *testing.T) {

	t.Run("Save and load data", func(t *testing.T) {
		dataFolder := path.Join(t.TempDir(), "data")
		config := &LocalConfig{
			ToDoListName: "todo",
			Gist:         GetDefaultGistConfig(),
			Backup:       GetDefaultBackupConfig(dataFolder),
		}
		os.MkdirAll(config.Backup.BackupFolder, 0755)
		filename := path.Join(dataFolder, "todo-cli.yaml")
		err := config.Save(filename)
		if err != nil {
			t.Errorf("Save() error = %v", err)
			return
		}

		config2, err := LoadFromFile(filename)
		if err != nil {
			t.Errorf("LoadLocalConfig() error = %v", err)
			return
		}
		if config.ToDoListName != config2.ToDoListName {
			t.Errorf("ToDoListName = %v, want %v", config2.ToDoListName, config.ToDoListName)
		}

	})

}
