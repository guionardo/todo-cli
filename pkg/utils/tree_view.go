package utils

import (
	"io"
	"strings"
)

const (
	TreeBranchChar = "├── "
	TreeNodeChar   = "│   "
	TreeLastChar   = "└── "
)

type TreeNode struct {
	id       string
	Text     string
	children []*TreeNode
}

func NewTreeNode(id string) *TreeNode {
	return &TreeNode{id: id, children: make([]*TreeNode, 0, 5)}
}

func (t *TreeNode) FindById(id string) *TreeNode {
	if t.id == id {
		return t
	}
	for _, child := range t.children {
		if tv := child.FindById(id); tv != nil {
			return tv
		}
	}
	return nil
}
func (t *TreeNode) AddChild(id string, childText string) (child *TreeNode) {
	if t.children == nil {
		t.children = make([]*TreeNode, 0)
	}
	child = &TreeNode{id: id, Text: childText, children: make([]*TreeNode, 0, 5)}
	t.children = append(t.children, child)
	return
}

func (t *TreeNode) Print(w io.Writer) {
	t.printLevel(0, w, "")
}

func (t *TreeNode) printLevel(level int, w io.Writer, prefix string) {
	if level > 0 {
		prefix = strings.Repeat(TreeNodeChar, level-1) + prefix
		io.WriteString(w, prefix)
		io.WriteString(w, t.Text+"\n")
	}

	for i, child := range t.children {
		switch i {
		case len(t.children) - 1:
			child.printLevel(level+1, w, TreeLastChar)
		case 0:
			child.printLevel(level+1, w, TreeBranchChar)
		default:
			child.printLevel(level+1, w, TreeNodeChar)
		}
	}
}
