package todo

import (
	"fmt"
	"time"
)

func (collection *ToDoCollection) Sorted() (items []*ToDoItem) {
	items = make([]*ToDoItem, len(collection.Items))
	i := 0
	for _, item := range collection.Items {
		items[i] = item
		i++
	}
	SortList(items)
	return
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
	collection.Items[item.Id] = item
	collection.LastUpdate = time.Now()
}

func (collection *ToDoCollection) Get(index int) *ToDoItem {
	for _, item := range collection.Items {
		if item.Index == index {
			return item
		}
	}
	return nil
}

func (collection *ToDoCollection) Remove(index int) error {
	item := collection.Get(index)
	if item == nil {
		return fmt.Errorf("Id not found %d", index)
	}
	delete(collection.Items, item.Id)
	collection.DeletedItems = append(collection.DeletedItems, item.Id)
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *ToDoCollection) DoAct(index int) error {
	item := collection.Get(index)
	if item == nil {
		return fmt.Errorf("Id not found %d", index)
	}
	item.LastAction = time.Now()
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *ToDoCollection) Complete(index int, undo bool) error {
	item := collection.Get(index)
	if item == nil {
		return fmt.Errorf("Id not found %d", index)
	}
	item.Completed = !undo
	collection.LastUpdate = time.Now()
	return nil
}
