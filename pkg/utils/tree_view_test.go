package utils

import (
	"bytes"
	"testing"
)

func TestTreeView_Tests(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		tv := NewTreeNode("root")
		c1 := tv.AddChild("c1", "child1")
		c2 := tv.AddChild("c2", "child2")
		c3 := tv.AddChild("c3", "child3")
		c1_1 := c1.AddChild("c1_1", "child1_1")
		c1_2 := c1.AddChild("c1_2", "child1_2")
		c2_1 := c2.AddChild("c2_1", "child2_1")
		if c3 == nil {
			t.Error("c3 is nil")
		}
		if c1_1 == nil || c1_2 == nil || c2_1 == nil {
			t.Error("c1_1, c1_2 or c2_1 is nil")
		}
		out := bytes.NewBufferString("")
		expected := `├── child1
│   ├── child1_1
│   └── child1_2
│   child2
│   └── child2_1
└── child3
`
		tv.Print(out)
		if out.String() != expected {
			t.Errorf("Expected: %s, got: %s", expected, out.String())
		}

	})

}
