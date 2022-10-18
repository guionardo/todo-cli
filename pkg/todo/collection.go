package todo

import "time"

type ToDoCollection struct {
	Items        map[string]*ToDoItem `yaml:"items"`
	DeletedItems []string             `yaml:"deleted_items"`
	LastUpdate   time.Time            `yaml:"last_update"`
	LastSave     time.Time            `yaml:"last_save"`
	LastSync     time.Time            `yaml:"last_sync"`
}

func NewTodoCollection() *ToDoCollection {
	return &ToDoCollection{
		Items:        make(map[string]*ToDoItem),
		DeletedItems: make([]string, 0),
		LastUpdate:   time.Now(),
		LastSave:     time.Time{},
		LastSync:     time.Time{},
	}
}
