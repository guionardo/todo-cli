package todo

import (
	"testing"
)

func TestCollection_Equal(t *testing.T) {
	items1A := make(map[string]*Item, 0)
	items1A["A"] = &Item{
		Id: "A",
	}
	items1B := make(map[string]*Item, 0)
	items1B["B"] = &Item{
		Id: "B",
	}
	tests := []struct {
		name       string
		collection *Collection
		other      any
		want       bool
	}{
		{"Equal", &Collection{
			Items:        make(map[string]*Item, 0),
			DeletedItems: []string{},
		}, &Collection{
			Items:        make(map[string]*Item, 0),
			DeletedItems: []string{},
		}, true},
		{"Different", &Collection{
			Items:        items1A,
			DeletedItems: []string{"A"},
		}, &Collection{
			Items:        items1B,
			DeletedItems: []string{"B"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collection.Equal(tt.other); got != tt.want {
				t.Errorf("Collection.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
