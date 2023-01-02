package todo

import "time"

type Collection struct {
	Items        map[string]*Item `yaml:"items"`
	DeletedItems []string         `yaml:"deleted_items"`
	LastUpdate   time.Time        `yaml:"last_update"`
	LastSave     time.Time        `yaml:"last_save"`
	LastSync     time.Time        `yaml:"last_sync"`
}

func NewTodoCollection() *Collection {
	return &Collection{
		Items:        make(map[string]*Item),
		DeletedItems: make([]string, 0),
		LastUpdate:   time.Now(),
		LastSave:     time.Time{},
		LastSync:     time.Time{},
	}
}
