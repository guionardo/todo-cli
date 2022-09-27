package internal

import (
	"reflect"
	"testing"
	"time"
)

func Test_extractCompleted(t *testing.T) {

	tests := []struct {
		name  string
		line  string
		want  bool
		want1 string
	}{
		{name: "completed", line: completeSingleMD, want: true, want1: "Complete"},
		{name: "not completed", line: singleMD, want: false, want1: "Pending"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, noIndex := extractIndex(tt.line)
			got, got1 := extractCompleted(noIndex)
			if got != tt.want {
				t.Errorf("extractCompleted() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractCompleted() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_extractTags(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		want  []string
		want1 string
	}{
		{name: "one tag", line: "#tag Title", want: []string{"tag"}, want1: "Title"},
		{name: "no tag", line: "Title", want: []string{}, want1: "Title"},
		{name: "two tags", line: "#tag1 #tag2 Title", want: []string{"tag1", "tag2"}, want1: "Title"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractTags(tt.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractTags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_extractDueTo(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		want  time.Time
		want1 string
	}{
		{name: "No due to", line: "Title", want: time.Time{}, want1: "Title"},
		{name: "Due to", line: "(2020-01-01) Title", want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), want1: "Title"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractDueTo(tt.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractDueTo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractDueTo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestToDoItem_NotifyText(t *testing.T) {
	dayAgo := time.Now().Add(-24 * time.Hour)
	tests := []struct {
		name string
		item ToDoItem
		want string
	}{
		{name: "completed", item: ToDoItem{Title: "Title", Completed: true, UpdatedAt: dayAgo}, want: "Completed @ " + dayAgo.Format(DateTimeFormat)},
		{name: "new item", item: ToDoItem{Title: "Title", UpdatedAt: dayAgo}, want: "New (1 days)"},
		{name: "due to", item: ToDoItem{Title: "Title", DueTo: dayAgo, UpdatedAt: dayAgo}, want: "Overdue -1 days"},
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := tt.item
			if got := item.NotifyText(); got != tt.want {
				t.Errorf("ToDoItem.NotifyText() = %v, want %v", got, tt.want)
			}
		})
	}
}
