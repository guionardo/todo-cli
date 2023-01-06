package todo

import (
	"testing"
)

func TestParseTags(t *testing.T) {

	t.Run("Default", func(t *testing.T) {
		got := ParseTags([]string{"tag1", "tag2", "tag2", "tag3", ""})
		wanted := []string{"tag1", "tag2", "tag3"}
		if !EqualTags(got, wanted) {
			t.Errorf("ParseTags() = %v, want %v", got, wanted)
		}
	})

}
