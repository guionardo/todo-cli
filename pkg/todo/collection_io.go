package todo

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func LoadCollection(filename string) (collection Collection, err error) {
	var content []byte
	if content, err = os.ReadFile(filename); err == nil {
		return LoadCollectionFromData(content)
	}

	return
}

func LoadCollectionFromData(data []byte) (collection Collection, err error) {
	err = yaml.Unmarshal(data, &collection)
	if err == nil {
		collection.UpdateLevels()
	}
	return
}

func (collection *Collection) Save(filename string) (err error) {
	collection.UpdateLevels()
	lastSave := collection.LastSave
	collection.LastSave = time.Now()
	var content []byte
	if content, err = yaml.Marshal(collection); err == nil {
		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			collection.LastSave = lastSave
		}
	}
	if err != nil {
		collection.LastSave = lastSave
	}
	return
}

func (collection *Collection) getLevel(id string) int {
	item, ok := collection.Items[id]
	if !ok {
		return -1
	}
	if item.ParentId == "" {
		item.Level = 0
	} else {
		item.Level = collection.getLevel(item.ParentId) + 1
	}
	return item.Level
}

func (collection *Collection) UpdateLevels() {
	for _, item := range collection.Items {
		collection.getLevel(item.Id)
	}
}
