package config

import (
	"os"
	"strings"
)

type LocalConfig struct {
	ToDoListName string       `yaml:"todo_list_name"`
	Gist         GistConfig   `yaml:"gist"`
	Backup       BackupConfig `yaml:"backup"`
}

func (c *LocalConfig) String() string {
	return strings.Join([]string{
		"LocalConfig{",
		" ToDoListName: " + c.ToDoListName,
		" Gist: " + c.Gist.String(),
		" Backup: " + c.Backup.String(),
		"}",
	}, "\n")
}

func GetDefaultLocalConfig(dataFolder string) *LocalConfig {
	return &LocalConfig{
		ToDoListName: "todo",
		Gist:         GetDefaultGistConfig(),
		Backup:       GetDefaultBackupConfig(dataFolder),
	}
}

func LoadFromFile(filename string) (cfg LocalConfig, err error) {
	var content []byte
	if content, err = os.ReadFile(filename); err == nil {
		err = cfg.Parse(content)
	}

	return
}

func validConfigFile(backupFile string) error {
	_, err := LoadFromFile(backupFile)
	return err
}
func lastBackupConfigFileIsSameConfig(lastBackupFile string, c LocalConfig) bool {
	if len(lastBackupFile) == 0 {
		return false
	}
	cfg, err := LoadFromFile(lastBackupFile)

	if err != nil {
		return false
	}
	return c.Equal(&cfg)
}

// func (c *LocalConfig) DoBackup(fileName string, cfg *BackupConfig) error {
// 	prefix := strings.TrimSuffix(path.Base(fileName), path.Ext(fileName))
// 	backupFiles, err := utils.NewBackupFiles(cfg.BackupFolder, prefix, cfg.BackupMaxCount, validConfigFile)
// 	if err != nil {
// 		return err
// 	}

// }
