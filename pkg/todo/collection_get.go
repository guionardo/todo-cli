package todo

import "sort"

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

func (collection *ToDoCollection) GetChildren(id string) []*ToDoItem {
	items := make([]*ToDoItem, 0)
	for _, item := range collection.Items {
		if item.ParentId == id {
			items = append(items, item)
		}
	}
	return items
}

func SortList(list []*ToDoItem) {
	sort.SliceStable(list, func(i int, j int) bool {
		return list[i].DueTo.Before(list[j].DueTo)
	})
}