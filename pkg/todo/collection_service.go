package todo

import (
	"time"

	"github.com/guionardo/todo-cli/pkg/exceptions"
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
		return exceptions.TodoNotFoundError(index)
	}
	children := collection.GetChildren(item.Id)
	if len(children) > 0 {
		return exceptions.TodoItemHasChildrenError(index, len(children))
	}
	delete(collection.Items, item.Id)
	collection.DeletedItems = append(collection.DeletedItems, item.Id)
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *Collection) DoAct(index int) error {
	item := collection.Get(index)
	if item == nil {
		return exceptions.TodoNotFoundError(index)
	}
	item.LastAction = time.Now()
	collection.LastUpdate = time.Now()
	return nil
}

func (collection *Collection) Complete(index int, undo bool, recursive bool) error {
	item := collection.Get(index)
	if item == nil {
		return exceptions.TodoNotFoundError(index)
	}
	if item.Completed != undo {
		return exceptions.NoChangedTodoError(index)
	}
	if undo {
		if item.Completed {
			item.Completed = false
			item.UpdatedAt = time.Now()
			collection.LastUpdate = time.Now()
			return nil
		}
		return exceptions.NoChangedTodoError(index)
	}
	children := collection.GetChildren(item.Id)
	openChildren := make([]*Item, 0, len(children))
	if children != nil {
		for _, child := range children {
			if !child.Completed {
				openChildren = append(openChildren, child)
			}
		}
	}

	if len(openChildren) > 0 {
		if recursive {
			for _, child := range openChildren {
				if err := collection.Complete(child.Index, false, true); err != nil {
					return err
				}
			}
		} else {
			return exceptions.TodoItemHasChildrenError(index, len(children))
		}
	}

	item.Completed = true
	item.UpdatedAt = time.Now()
	collection.LastUpdate = time.Now()
	return nil
}
