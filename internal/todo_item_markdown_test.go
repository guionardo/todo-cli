package internal

import (
	"testing"
	"time"
)

var (
	singleItem         = ToDoItem{Index: 1, Title: "Pending"}
	completeSingleItem = ToDoItem{Index: 2, Title: "Complete", Completed: true}
	taggedItem         = ToDoItem{Index: 3, Title: "Tagged", Tags: []string{"tag1", "tag2"}}
	dueToItem          = ToDoItem{Index: 4, Title: "Due to", DueTo: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	taggedDueToItem    = ToDoItem{Index: 5, Title: "Tagged and due to", Tags: []string{"tag1", "tag2"}, DueTo: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	lastActionItem     = ToDoItem{Index: 6, Title: "Last action", LastAction: time.Date(2020, 1, 1, 23, 44, 0, 0, time.UTC)}
)

const (
	singleMD         = "- [ ] 001 Pending"
	completeSingleMD = "- [x] 002 Complete"
	taggedMD         = "- [ ] 003 #tag1 #tag2 Tagged"
	dueToMD          = "- [ ] 004 (2020-01-01) Due to"
	taggedDueToMD    = "- [ ] 005 #tag1 #tag2 (2020-01-01) Tagged and due to"
	lastActionMD     = "- [ ] 006 Last action {2020-01-01 23:44:00}"

	singleStr         = "001 ⌛ Pending"
	completeSingleStr = "002 ✔ Complete"
	taggedStr         = "003 ⌛ #tag1 #tag2 Tagged"
	dueToStr          = "004 ⌛ (2020-01-01) Due to"
	taggedDueToStr    = "005 ⌛ #tag1 #tag2 (2020-01-01) Tagged and due to"
	lastActionStr     = "006 ⌛ Last action {2020-01-01 23:44:00}"
)

func TestToDoItem_ToMarkDown(t *testing.T) {

	tests := []struct {
		name string
		item ToDoItem
		want string
	}{
		{name: "Single", item: singleItem, want: singleMD},
		{name: "Complete", item: completeSingleItem, want: completeSingleMD},
		{name: "Tagged", item: taggedItem, want: taggedMD},
		{name: "Due to", item: dueToItem, want: dueToMD},
		{name: "Tagged and due to", item: taggedDueToItem, want: taggedDueToMD},
		{name: "Last action", item: lastActionItem, want: lastActionMD},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.item.ToMarkDown(); got != tt.want {
				t.Errorf("ToDoItem.ToMarkDown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDoItem_FromMarkDown(t *testing.T) {
	tests := []struct {
		name    string
		item    ToDoItem
		line    string
		wantErr bool
	}{
		{name: "Single", item: singleItem, line: singleMD, wantErr: false},
		{name: "Completed", item: completeSingleItem, line: completeSingleMD, wantErr: false},
		{name: "Tagged", item: taggedItem, line: taggedMD, wantErr: false},
		{name: "Due to", item: dueToItem, line: dueToMD, wantErr: false},
		{name: "Tagged and due to", item: taggedDueToItem, line: taggedDueToMD, wantErr: false},
		{name: "Last action", item: lastActionItem, line: lastActionMD, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			parsedItem := &ToDoItem{}
			if err := parsedItem.FromMarkDown(tt.line); (err != nil) != tt.wantErr {
				t.Errorf("ToDoItem.FromMarkDown() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.item.ToMarkDown() != parsedItem.ToMarkDown() {
				t.Errorf("ToDoItem.FromMarkDown(\"%s\") got = %v, want %v", tt.line, parsedItem, tt.item)
			}
		})
	}
}

func compareArrays(a, b []string) bool {
	if a == nil {
		a = make([]string, 0)
	}
	if b == nil {
		b = make([]string, 0)
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
