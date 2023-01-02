package todo

import "sort"

func (collection *Collection) GetByFilter(tags []string, justCompleted bool, justPending bool) []*Item {
	items := make([]*Item, 0)
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

func (collection *Collection) GetChildren(id string) []*Item {
	items := make([]*Item, 0)
	for _, item := range collection.Items {
		if item.ParentId == id {
			items = append(items, item)
		}
	}
	return items
}

func SortList(list []*Item) {
	sort.SliceStable(list, func(i int, j int) bool {
		return list[i].DueTo.Before(list[j].DueTo)
	})
}
