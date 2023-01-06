package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
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

func LoadLocalConfig(filename string) (cfg LocalConfig, err error) {
	var content []byte
	if content, err = os.ReadFile(filename); err == nil {
		err = yaml.Unmarshal(content, &cfg)
	}

	return
}

func (c *LocalConfig) Save(filename string) error {
	if content, err := yaml.Marshal(c); err == nil {
		return os.WriteFile(filename, content, 0644)
	}

	return nil
}
