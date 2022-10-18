package todo

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func LoadCollection(filename string) (collection ToDoCollection, err error) {
	var content []byte
	if content, err = os.ReadFile(filename); err == nil {
		return LoadCollectionFromData(content)
	}

	return
}

func LoadCollectionFromData(data []byte) (collection ToDoCollection, err error) {
	err = yaml.Unmarshal(data, &collection)
	if err == nil {
		collection.UpdateLevels()
	}
	return
}

func (c *ToDoCollection) Save(filename string) (err error) {
	c.UpdateLevels()
	lastSave := c.LastSave
	c.LastSave = time.Now()
	var content []byte
	if content, err = yaml.Marshal(c); err == nil {
		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			c.LastSave = lastSave
		}
	}
	if err != nil {
		c.LastSave = lastSave
	}
	return
}

func (c *ToDoCollection) getLevel(id string) int {
	item, ok := c.Items[id]
	if !ok {
		return -1
	}
	if item.ParentId == "" {
		item.Level = 0
	} else {
		item.Level = c.getLevel(item.ParentId) + 1
	}
	return item.Level
}
func (c *ToDoCollection) UpdateLevels() {
	for _, item := range c.Items {
		c.getLevel(item.Id)
	}
}
