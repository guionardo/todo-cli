package internal

import "sort"

func MergeToDoItems(items1 []*ToDoItem, items2 []*ToDoItem) []*ToDoItem {
	items_map := make(map[string]*ToDoItem)

	for _, item1 := range items1 {
		items_map[item1.Id] = item1
	}
	for _, item2 := range items2 {
		if item, ok := items_map[item2.Id]; ok && item.UpdatedAt.After(item2.UpdatedAt) {
			continue
		}
		items_map[item2.Id] = item2
	}
	new_items := make([]*ToDoItem, len(items_map))
	index := 0
	for _, v := range items_map {
		new_items[index] = v
		index++
	}
	sort.SliceStable(new_items, func(i int, j int) bool {
		return new_items[i].DueTo.Before(new_items[j].DueTo)
	})
	for i, item := range new_items {
		item.Index = i+1
	}
	return new_items
}
