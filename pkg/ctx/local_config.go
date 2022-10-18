package ctx

import (
	"os"

	"github.com/guionardo/todo-cli/pkg/backup"
	"github.com/guionardo/todo-cli/pkg/github"
	"gopkg.in/yaml.v3"
)

type LocalConfig struct {
	ToDoListName string              `yaml:"todo_list_name"`
	Gist         github.GistConfig   `yaml:"gist"`
	Backup       backup.BackupConfig `yaml:"backup"`
}

func GetDefaultLocalConfig(dataFolder string) *LocalConfig {
	return &LocalConfig{
		ToDoListName: "todo",
		Gist:         github.GetDefaultGistConfig(),
		Backup:       backup.GetDefaultBackupConfig(dataFolder),
	}
}

func LoadLocalConfig(filename string) (config LocalConfig, err error) {
	var content []byte
	content, err = os.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, &config)
	return
}

func (c *LocalConfig) Save(filename string) error {
	if content, err := yaml.Marshal(c); err == nil {
		return os.WriteFile(filename, content, 0644)
	}

	return nil
}
