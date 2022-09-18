package internal

import (
	"path"
	"reflect"
	"testing"
)

func TestToDoCollection_SaveLoad(t *testing.T) {

	t.Run("Save and load", func(t *testing.T) {
		config := Config{
			GistId:        "gist_id",
			Authorization: "github_token",
			ToDoListName: "todo.md",
		}
		collection := &ToDoCollection{Config: config}
		collection.Add(&ToDoItem{Title: "Title 1"})
		collection.Add(&ToDoItem{Title: "Title 2", Completed: true})

		fileName := path.Join(t.TempDir(), "collection.yaml")

		if err := collection.Save(fileName); err != nil {
			t.Errorf("ToDoCollection.Save() error = %v", err)
		}

		col2, err := ParseCollectionFile(fileName)
		if err != nil {
			t.Errorf("ToDoCollection.Load() error = %v", err)
		}
		if !reflect.DeepEqual(collection, col2) {
			t.Errorf("ToDoCollection.Load() = %v, want %v", col2, collection)
		}
	})

}
