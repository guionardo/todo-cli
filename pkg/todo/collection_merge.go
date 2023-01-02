package todo

import (
	"fmt"
	"reflect"
)

type (
	mergeData struct {
		newCol        map[string]*Item
		deleted_items map[string]void
		log           map[string]string
		diffCount     int
		upload        bool
	}
	void struct{}
)

func findItem(col []*Item, id string) *Item {
	for _, item := range col {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func mergeList(col1 []*Item, col2 []*Item, data *mergeData) {
	for _, item := range col1 {
		diff := false
		col2Item := findItem(col2, item.Id)
		if col2Item == nil {
			diff = true
			col2Item = item
			data.log[item.Id] = ">"
		} else if reflect.DeepEqual(item, col2Item) {
			data.log[item.Id] = "="
		} else if col2Item.UpdatedAt.After(item.UpdatedAt) {
			item = col2Item
			data.log[item.Id] = "<"
			diff = true
		} else {
			data.log[item.Id] = ">"
			diff = true
		}
		if item.Deleted {
			data.deleted_items[item.Id] = void{}
			data.log[item.Id] = "X"
			diff = true
		} else {
			data.newCol[item.Id] = item
		}
		if diff {
			data.diffCount++
		}
		if data.log[item.Id] == ">" {
			data.upload = true
		}
	}
}

func MergeCollections(col1 []*Item, col2 []*Item) (newCol []*Item, deleted_items []string, log []string, diffCount int, upload bool) {
	data := &mergeData{
		newCol:        make(map[string]*Item),
		deleted_items: make(map[string]void),
		log:           make(map[string]string),
	}
	mergeList(col1, col2, data)
	mergeList(col2, col1, data)
	newCol = make([]*Item, len(data.newCol))
	i := 0
	for _, item := range data.newCol {
		newCol[i] = item
		i++
	}
	i = 0
	deleted_items = make([]string, len(data.deleted_items))
	for item := range data.deleted_items {
		deleted_items[i] = item
		i++
	}
	log = make([]string, len(data.log))
	i = 0
	for item, action := range data.log {
		log[i] = fmt.Sprintf("%s %s", action, item)
		i++
	}
	diffCount = data.diffCount
	upload = data.upload
	return
}

func (collection *Collection) Merge(other *Collection) (log []string, diffCount int, upload bool) {
	thisItems := collection.Sorted()
	otherItems := other.Sorted()

	newCol, deleted_items, log, diffCount, upload := MergeCollections(thisItems, otherItems)

	collection.Items = map[string]*Item{}
	for _, item := range newCol {
		collection.Items[item.Id] = item
	}
	collection.DeletedItems = deleted_items
	return log, diffCount, upload
}
