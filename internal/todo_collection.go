package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ToDoCollection struct {
	Items      []*ToDoItem
	LastUpdate time.Time
	LastSave   time.Time
	LastSync   time.Time
	Config     Config
}

func CollectionFile() (string, error) {
	if todo_collection := os.Getenv("TODO_COLLECTION"); todo_collection != "" {
		if _, err := os.Stat(todo_collection); err == nil {
			_, err := ParseCollectionFile(todo_collection)
			if err == nil {
				return todo_collection, nil
			}
			log.Printf("Error parsing collection file %s: %s", todo_collection, err)
		}
	}
	return "", nil
}

func ParseCollectionFile(collectionFile string) (*ToDoCollection, error) {
	file, err := os.Open(collectionFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	collection := &ToDoCollection{}
	return collection, yaml.Unmarshal(body, collection)
}

func (collection *ToDoCollection) Save(collectionFile string) error {
	body, err := yaml.Marshal(collection)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(collectionFile, body, 0644)
}

func (collection *ToDoCollection) Add(item *ToDoItem) {
	collection.Items = append(collection.Items, item)
}

func (collection *ToDoCollection) Complete(id int) error {
	if id < 0 || id >= len(collection.Items) {
		return fmt.Errorf("Invalid id %d", id)
	}
	if collection.Items[id].Completed {
		return fmt.Errorf("Item %d already completed", id)
	}
	collection.Items[id].Completed = true
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *ToDoCollection) Remove(id int) error {
	if id < 0 || id >= len(collection.Items) {
		return fmt.Errorf("Invalid id %d", id)
	}
	collection.Items = append(collection.Items[:id], collection.Items[id+1:]...)
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *ToDoCollection) Get(id int) (*ToDoItem, error) {
	if id < 0 || id >= len(collection.Items) {
		return nil, fmt.Errorf("Invalid id %d", id)
	}
	return collection.Items[id], nil
}

func (collection *ToDoCollection) GetByTag(tag string) []*ToDoItem {
	items := make([]*ToDoItem, 0)

	for _, item := range collection.Items {
		for _, itemTag := range item.Tags {
			if itemTag == tag {
				items = append(items, item)
			}
		}
	}
	return items
}
