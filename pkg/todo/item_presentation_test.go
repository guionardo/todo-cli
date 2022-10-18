package todo

import (
	"testing"
	"time"
)

func TestToDoCollection_GetTreeList(t *testing.T) {
	t.Run("TreeList", func(t *testing.T) {
		collection := &ToDoCollection{
			Items:        make(map[string]*ToDoItem),
			DeletedItems: make([]string, 0),
		}
		collection.Add(&ToDoItem{
			Id:    "1",
			Index: 1,
			Title: "Item 1",
			DueTo: time.Now().Add(time.Hour * 24),
		})
		collection.Add(&ToDoItem{
			Id:    "2",
			Index: 2,
			Title: "Item 2",
		})
		collection.Add(&ToDoItem{
			Id:       "3",
			Index:    3,
			Title:    "Item 2.1",
			ParentId: "2",
			DueTo:    time.Now().Add(time.Hour * 24 * 2),
		})
		collection.Add(&ToDoItem{
			Id:    "4",
			Index: 4,
			Title: "Item 3",
		})
		collection.Add(&ToDoItem{
			Id:       "5",
			Index:    5,
			Title:    "Item 2.2",
			ParentId: "2",
		})
		items := collection.GetByFilter([]string{}, false, false)
		got := collection.GetTreeList(items)
		if len(got) != 5 {
			t.Errorf("GetTreeList() = %v, want %v", got, 6)
		}
		for _, item := range got {
			t.Logf("%v", item)
		}

	})

}

func TestToDoItem_StringNoColor(t *testing.T) {
	tests := []struct {
		name string
		item *ToDoItem
		want string
	}{
		{name: "Just title", item: &ToDoItem{Title: "Test"}, want: "#000 Test"},
		{
			name: "Title with tags",
			item: &ToDoItem{Index: 1, Title: "Test", Tags: []string{"tag1", "tag2"}},
			want: "#001 (tag1 tag2) Test",
		},
		{
			name: "Title with due date delayed",
			item: &ToDoItem{Index: 1, Title: "Test", DueTo: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
			want: "#001 ⌛ (2021-01-01) Test",
		},
		{
			name: "Title with due date future",
			item: &ToDoItem{Index: 1, Title: "Test", DueTo: time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)},
			want: "#001 ⏱ (2099-12-31) Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.StringNoColor()
			if got != tt.want {
				t.Errorf("ToDoItem.StringNoColor() = %v, want %v", got, tt.want)
			} else {
				t.Logf("%s", got)
			}
		})
	}
}
