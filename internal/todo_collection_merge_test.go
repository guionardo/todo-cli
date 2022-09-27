package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestMergeToDoItems(t *testing.T) {

	t.Run("MergeToDoItems", func(t *testing.T) {
		items1 := []*ToDoItem{
			{Id: NewItemId(), Title: "item1", UpdatedAt: time.Now()},
			{Id: NewItemId(), Title: "item2", UpdatedAt: time.Now()},
			{Id: NewItemId(), Title: "item3", UpdatedAt: time.Now()},
		}
		items2 := []*ToDoItem{
			{Id: items1[0].Id, Title: "itemA", UpdatedAt: time.Now().Add(time.Hour)},
			{Id: items1[1].Id, Title: "itemB", UpdatedAt: time.Now().Add(time.Hour)},
			{Id: NewItemId(), Title: "itemC", UpdatedAt: time.Now()},
		}
		new_items := MergeToDoItems(items1, items2)

		expected := []*ToDoItem{
			items2[0],
			items2[1],
			items1[2],
			items2[2],
		}
		if !reflect.DeepEqual(new_items, expected) {
			t.Errorf("MergeToDoItems() = %v, want %v", new_items, expected)
		}
	})

}
