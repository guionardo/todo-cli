package todo

import (
	"github.com/guionardo/todo-cli/pkg/interfaces"
	"gopkg.in/yaml.v3"
)

func (c *Collection) Equal(other any) bool {
	if o, ok := other.(*Collection); ok {
		if c == nil || o == nil || len(c.Items) != len(o.Items) {
			return false
		}

		for _, item := range c.Items {
			if otherItem, ok := o.Items[item.Id]; !ok || !otherItem.Equal(item) {
				return false
			}
		}
		for _, otherItem := range o.Items {
			if item, ok := c.Items[otherItem.Id]; !ok || !item.Equal(otherItem) {
				return false
			}
		}
		return true

	}
	return false
}

func (c *Collection) Parse(data []byte) error {
	if parsed, err := c.ParseNew(data); err == nil {
		p := parsed.(*Collection)
		c.Items = p.Items
		c.DeletedItems = p.DeletedItems
		c.LastUpdate = p.LastUpdate
		c.LastSave = p.LastSave
		c.LastSync = p.LastSync
		return nil
	} else {
		return err
	}
}

func (c *Collection) ParseNew(data []byte) (interfaces.Backupable, error) {
	var parsed Collection
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (c *Collection) BackupPrefix() string {
	return "todo"
}
