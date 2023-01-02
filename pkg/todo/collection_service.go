package todo

import (
	"fmt"
	"time"
)

func (collection *Collection) Sorted() (items []*Item) {
	items = make([]*Item, len(collection.Items))
	i := 0
	for _, item := range collection.Items {
		items[i] = item
		i++
	}
	SortList(items)
	return
}

func (collection *Collection) Add(item *Item) {
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

func (collection *Collection) Get(index int) *Item {
	for _, item := range collection.Items {
		if item.Index == index {
			return item
		}
	}
	return nil
}

func (collection *Collection) Remove(index int) error {
	item := collection.Get(index)
	if item == nil {
		return fmt.Errorf("id not found %d", index)
	}
	delete(collection.Items, item.Id)
	collection.DeletedItems = append(collection.DeletedItems, item.Id)
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *Collection) DoAct(index int) error {
	item := collection.Get(index)
	if item == nil {
		return fmt.Errorf("id not found %d", index)
	}
	item.LastAction = time.Now()
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *Collection) Complete(index int, undo bool) error {
	item := collection.Get(index)
	if item == nil {
		return fmt.Errorf("id not found %d", index)
	}
	item.Completed = !undo
	collection.LastUpdate = time.Now()
	return nil
}
