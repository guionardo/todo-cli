package config

import (
	"os"

	"github.com/guionardo/todo-cli/pkg/interfaces"
	"gopkg.in/yaml.v3"
)

func (c *LocalConfig) Equal(other any) bool {
	if o, ok := other.(*LocalConfig); ok {
		return c.ToDoListName == o.ToDoListName &&
			c.Gist.Equal(&o.Gist) &&
			c.Backup.Equal(&o.Backup)
	}
	return false
}

func (c *LocalConfig) Parse(data []byte) error {
	if parsed, err := c.ParseNew(data); err != nil {
		return nil
	} else {
		p := parsed.(*LocalConfig)
		c.ToDoListName = p.ToDoListName
		c.Gist = p.Gist
		c.Backup = p.Backup
	}
	return nil
}

func (c *LocalConfig) ParseNew(data []byte) (interfaces.Backupable, error) {
	var parsed LocalConfig
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}
func (c *LocalConfig) Save(fileName string) error {
	if content, err := yaml.Marshal(c); err != nil {
		return err
	} else {
		return os.WriteFile(fileName, content, 0644)
	}
}

func (c *LocalConfig) BackupPrefix() string {
	return "config"
}
