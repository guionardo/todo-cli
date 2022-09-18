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
		{name: "completed", line: "- [x] Title", want: true, want1: "Title"},
		{name: "not completed", line: "- [ ] Title 2", want: false, want1: "Title 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractCompleted(tt.line)
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

func TestToDoItem_FromMarkDown(t *testing.T) {
	type fields struct {
		Title     string
		Completed bool
		DueTo     time.Time
		Tag       []string
	}

	tests := []struct {
		name    string
		fields  fields
		line    string
		wantErr bool
	}{
		{name: "Single", fields: fields{Title: "Title", Completed: false, DueTo: time.Time{}, Tag: []string{}}, line: "- [ ] Title", wantErr: false},
		{name: "Completed", fields: fields{Title: "Title", Completed: true, DueTo: time.Time{}, Tag: []string{}}, line: "- [x] Title", wantErr: false},
		{name: "Tagged", fields: fields{Title: "Title", Completed: false, DueTo: time.Time{}, Tag: []string{"tag"}}, line: "- [ ] #tag Title", wantErr: false},
		{name: "Due to", fields: fields{Title: "Title", Completed: false, DueTo: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Tag: []string{}}, line: "- [ ] (2020-01-01) Title", wantErr: false},
		{name: "Tagged and due to", fields: fields{Title: "Title", Completed: false, DueTo: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Tag: []string{"tag"}}, line: "- [ ] #tag (2020-01-01) Title", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &ToDoItem{
				Title:     tt.fields.Title,
				Completed: tt.fields.Completed,
				DueTo:     tt.fields.DueTo,
				Tags:      tt.fields.Tag,
			}
			parsedItem := &ToDoItem{}
			if err := parsedItem.FromMarkDown(tt.line); (err != nil) != tt.wantErr {
				t.Errorf("ToDoItem.FromMarkDown() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(item, parsedItem) {
				t.Errorf("ToDoItem.FromMarkDown() got = %v, want %v", parsedItem, item)
			}
		})
	}
}

func TestToDoItem_ToMarkDown(t *testing.T) {

	tests := []struct {
		name string
		item ToDoItem
		want string
	}{
		{name: "Single", item: ToDoItem{Title: "Title", Completed: false, DueTo: time.Time{}, Tags: []string{}}, want: "- [ ] Title"},
		{name: "Completed", item: ToDoItem{Title: "Title", Completed: true, DueTo: time.Time{}, Tags: []string{}}, want: "- [x] Title"},
		{name: "Tagged", item: ToDoItem{Title: "Title", Completed: true, DueTo: time.Time{}, Tags: []string{"tag"}}, want: "- [x] #tag Title"},
		{name: "Tagged 2", item: ToDoItem{Title: "Title", Completed: true, DueTo: time.Time{}, Tags: []string{"tag", "tag2"}}, want: "- [x] #tag #tag2 Title"},
		{name: "Due to", item: ToDoItem{Title: "Title", Completed: true, DueTo: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Tags: []string{}}, want: "- [x] (2020-01-01) Title"},
		{name: "Tagged and due to", item: ToDoItem{Title: "Title", Completed: true, DueTo: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Tags: []string{"tag", "tag2"}}, want: "- [x] #tag #tag2 (2020-01-01) Title"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.item.ToMarkDown(); got != tt.want {
				t.Errorf("ToDoItem.ToMarkDown() = %v, want %v", got, tt.want)
			}
		})
	}
}
