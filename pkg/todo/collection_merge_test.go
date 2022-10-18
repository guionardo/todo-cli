package todo

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestMergeCollections(t *testing.T) {
	t.Run("MergeToDoItems", func(t *testing.T) {
		items1 := []*ToDoItem{
			{Id: "I01", Title: "item1", UpdatedAt: time.Now()},
			{Id: "I02", Title: "item2", UpdatedAt: time.Now(), Deleted: true},
			{Id: "I03", Title: "item3", UpdatedAt: time.Now()},
		}
		items2 := []*ToDoItem{
			{Id: items1[0].Id, Title: "itemA", UpdatedAt: time.Now().Add(time.Hour)},
			{Id: items1[1].Id, Title: "itemB", UpdatedAt: time.Now().Add(time.Hour), Deleted: true},
			{Id: "I04", Title: "itemC", UpdatedAt: time.Now()},
		}
		new_items, deleted, logs, diffCount, upload := MergeCollections(items1, items2)
		t.Logf("changes: %v", logs)

		expected := []*ToDoItem{
			items2[0],
			items1[2],
			items2[2],
		}
		SortList(new_items)
		SortList(expected)
		if !reflect.DeepEqual(new_items, expected) {
			t.Errorf("MergeToDoItems() = %v, want %v", new_items, expected)
		}

		expected_deleted := []string{items1[1].Id}
		sort.Strings(deleted)
		sort.Strings(expected_deleted)
		if !reflect.DeepEqual(deleted, expected_deleted) {
			t.Errorf("MergeToDoItems() = %v, want %v", deleted, expected_deleted)
		}
		if diffCount != 6 {
			t.Errorf("MergeToDoItems() diffCount = %v, want %v", diffCount, 3)
		}
		if !upload {
			t.Errorf("MergeToDoItems() upload = %v, want %v", upload, true)
		}
	})
}

func TestMergeEqualCollections(t *testing.T) {
	t.Run("MergeEqual", func(t *testing.T) {
		i0 := &ToDoItem{Id: "I01", Title: "item1", UpdatedAt: time.Now()}
		i1 := &ToDoItem{Id: "I02", Title: "item2", UpdatedAt: time.Now(), Deleted: true}
		i2 := &ToDoItem{Id: "I03", Title: "item3", UpdatedAt: time.Now()}
		items1 := []*ToDoItem{i0, i1, i2}
		items2 := []*ToDoItem{i0, i1, i2}
		new_items, deleted, logs, diffCount, _ := MergeCollections(items1, items2)
		t.Logf("changes: %v", logs)
		if len(deleted) != 1 {
			t.Errorf("MergeEqual() deleted = %v, want %v", deleted, 1)
		}
		if len(new_items) != 2 {
			t.Errorf("MergeEqual() new_items = %v, want %v", new_items, 2)
		}
		if diffCount != 2 {
			t.Errorf("MergeEqual() diffCount = %v, want %v", diffCount, 2)
		}
	})
}
