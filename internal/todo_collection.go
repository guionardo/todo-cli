package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ToDoCollection struct {
	Items      []*ToDoItem `yaml:"items"`
	LastUpdate time.Time   `yaml:"last_update"`
	LastSave   time.Time   `yaml:"last_save"`
	LastSync   time.Time   `yaml:"last_sync"`
	Config     Config      `yaml:"config"`
}

func (c *ToDoCollection) String() string {
	return fmt.Sprintf("%s\nItems: %d\nLastSave: %s\nLastSync: %s\n", c.Config, len(c.Items), c.LastSave, c.LastSync)
}

func ParseCollectionData(collectionData []byte) (*ToDoCollection, error) {
	collection := &ToDoCollection{}
	if err := yaml.Unmarshal(collectionData, collection); err != nil {
		return nil, err
	}
	for _, item := range collection.Items {
		if len(item.Id) == 0 {
			item.Id = NewItemId()
		}
	}
	return collection, nil
}

func ParseCollectionFile(collectionFile string) (*ToDoCollection, error) {
	file, err := os.Open(collectionFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	collection, err := ParseCollectionData(body)
	collection.Config.ConfigFileName = collectionFile
	return collection, err
}

func (collection *ToDoCollection) Save(collectionFile string) error {
	lastSave := collection.LastSave
	collection.LastSave = time.Now()
	body, err := yaml.Marshal(collection)
	if err != nil {
		collection.LastSave = lastSave
		return err
	}
	err = ioutil.WriteFile(collectionFile, body, 0644)
	return err
}

func (collection *ToDoCollection) Add(item *ToDoItem) {
	maxId := 0
	for _, item := range collection.Items {
		if item.Index > maxId {
			maxId = item.Index
		}
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = time.Now()
	}
	item.Index = maxId + 1
	collection.Items = append(collection.Items, item)
	collection.LastUpdate = time.Now()
}

func (collection *ToDoCollection) Complete(id int) error {
	for _, item := range collection.Items {
		if item.Index == id {
			if item.Completed {
				return fmt.Errorf("Item %d already completed", id)
			}
			item.Completed = true
			item.UpdatedAt = time.Now()
			collection.LastUpdate = time.Now()
			return nil
		}
	}
	return fmt.Errorf("Item %d not found", id)
}

func (collection *ToDoCollection) UndoComplete(id int) error {
	for _, item := range collection.Items {
		if item.Index == id {
			if !item.Completed {
				return fmt.Errorf("Item %d pending yet", id)
			}
			item.Completed = false
			collection.LastUpdate = time.Now()
			item.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("Item %d not found", id)
}

func (collection *ToDoCollection) DoAct(id int) error {
	item, err := collection.Get(id)
	if err != nil {
		return err
	}
	if item.Completed {
		return fmt.Errorf("Item %d already completed", id)
	}
	item.LastAction = time.Now()
	item.UpdatedAt = time.Now()
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
	for _, item := range collection.Items {
		if item.Index == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("Item %d not found", id)
}

func (collection *ToDoCollection) GetByFilter(tags []string, justCompleted bool, justPending bool) []*ToDoItem {
	items := make([]*ToDoItem, 0)
	for _, item := range collection.Items {
		if justCompleted && !item.Completed {
			continue
		}
		if justPending && item.Completed {
			continue
		}
		if len(tags) == 0 {
			items = append(items, item)
			continue
		}
		for _, tag := range tags {
			for _, itemTag := range item.Tags {
				if itemTag == tag {
					items = append(items, item)
					break
				}
			}
		}
	}
	return items
}
