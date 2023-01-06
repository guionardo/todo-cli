package todo

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/guionardo/todo-cli/pkg/utils"
)

const (
	TimeFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	CompletedChar  = "✅ "
	OpenedChar     = ""
	PendingChar    = "⌛"
	ClockChar      = "⏱"
)

func (item *Item) String() string {
	var clr = item.ItemColor()

	return clr(item.StringNoColor())
}

func (item *Item) LevelString(level int) string {
	return strings.Repeat("  ", level) + "└ " + item.String()
}

func (item *Item) ItemColor() func(any) string {
	if item.Completed {
		return color.InGreen
	}
	if item.DueTo.IsZero() || item.DueTo.After(time.Now()) {
		return func(s any) string {
			return color.InBold(color.InYellow(s))
		}
	}
	if !item.DueTo.IsZero() && item.DueTo.Before(time.Now()) {
		return func(s any) string {
			return color.InBold(color.InRed(s))
		}
	}
	return color.InWhite
}

func (item *Item) StringNoColor() string {
	completed := utils.Tern(item.Completed, CompletedChar, OpenedChar)

	tags := utils.Tern(len(ParseTags(item.Tags)) > 0, fmt.Sprintf("(%s) ", strings.Join(item.Tags, " ")), "")
	lastAction := utils.Tern(item.LastAction.IsZero(), "", fmt.Sprintf("Last action: %s ", item.LastAction.Format(DateTimeFormat)))
	dueTo := utils.Tern(item.DueTo.IsZero() || item.Completed, "",
		utils.Tern(item.DueTo.Before(time.Now()), PendingChar, ClockChar)+fmt.Sprintf(" (%s) ", item.DueTo.Format(TimeFormat)))

	return fmt.Sprintf("#%03d %s%s%s%s%s", item.Index, completed, tags, dueTo, lastAction, item.Title)
}

func (item *Item) NotifyText() string {
	if item.Completed {
		return fmt.Sprintf("Completed @ %s", item.UpdatedAt.Format(DateTimeFormat))
	}
	if !item.DueTo.IsZero() {
		if item.DueTo.Before(time.Now()) {
			return fmt.Sprintf("Overdue %d days", int(item.DueTo.Sub(time.Now()).Hours()/24))
		}
		daysToComplete := time.Now().Sub(item.DueTo).Hours() / 24
		return fmt.Sprintf("Due in %d days @ %s", int(daysToComplete), item.DueTo.Format(DateTimeFormat))
	}
	lastAction := item.LastAction
	if lastAction.IsZero() {
		lastAction = item.UpdatedAt
	}
	if lastAction.IsZero() {
		lastAction = time.Now()
	}

	return fmt.Sprintf("New (%s)", utils.DurationString(time.Now().Sub(lastAction)))
}

func getSubList(allItems map[string]*Item, item *Item, level int) []string {
	delete(allItems, item.Id)
	// Header of item
	prefix := ""
	if level > 0 {
		prefix = strings.Repeat("  ", level) + "+ "
	}
	lines := []string{fmt.Sprintf("%s%s", prefix, item.String())}

	// Get list of children
	for _, child := range allItems {
		if child.ParentId == item.Id {
			childLines := getSubList(allItems, child, level+1)
			lines = append(lines, childLines...)
		}
	}
	return lines
}

func (collection *Collection) GetTreeList(items []*Item) []string {
	// Get list of items with no parents
	root := make([]*Item, 0)
	for _, item := range items {
		if item.ParentId == "" {
			root = append(root, item)
		}
	}
	// Sort by due date
	SortList(root)

	// Get map of items to avoid duplicates
	itemMap := make(map[string]*Item)
	for _, item := range items {
		itemMap[item.Id] = item
	}

	lines := make([]string, 0)
	for _, item := range root {
		sublist := getSubList(itemMap, item, 0)
		lines = append(lines, sublist...)
	}
	return lines

}

func parseChildren(allItems map[string]*Item, item *Item, tvRoot *utils.TreeNode) {
	delete(allItems, item.Id)
	node := tvRoot.FindById(item.Id)
	parent := tvRoot.FindById(item.ParentId)
	if parent == nil {
		parent = tvRoot
	}
	if node == nil {
		node = parent.AddChild(item.Id, item.String())
	} else {
		node.Text = item.String()
	}

	for id, child := range allItems {
		if child.ParentId == item.Id {
			node.AddChild(id, child.String())
			parseChildren(allItems, child, tvRoot)
		}
	}

}
func (collection *Collection) GetTreeView(items []*Item) []string {
	tv := utils.NewTreeNode("root")
	// Get list of items with no parents
	roots := make([]*Item, 0)
	for _, item := range items {
		if item.ParentId == "" {
			roots = append(roots, item)
			tv.AddChild(item.Id, item.String())
		}
	}

	// Sort by due date
	SortList(roots)

	// Get map of items to avoid duplicates
	itemMap := make(map[string]*Item)
	for _, item := range items {
		itemMap[item.Id] = item
	}

	for _, item := range roots {
		parseChildren(itemMap, item, tv)
	}
	out := bytes.NewBufferString("")
	tv.Print(out)
	return strings.Split(out.String(), "\n")

}

func getChildren(allItems map[string]*Item, item *Item) []*Item {
	delete(allItems, item.Id)
	children := make([]*Item, 0, 10)
	for _, child := range allItems {
		if child.ParentId == item.Id {
			children = append(children, child)
		}
	}
	return children
}

func getSubTree(allItems map[string]*Item, item *Item, level int) []string {
	delete(allItems, item.Id)
	// Header of item
	prefix := ""
	if level > 0 {
		prefix = strings.Repeat("  ", level) + "+ "
	}
	lines := []string{fmt.Sprintf("%s%s", prefix, item.String())}

	// Get list of children
	for _, child := range allItems {
		if child.ParentId == item.Id {
			childLines := getSubTree(allItems, child, level+1)
			lines = append(lines, childLines...)
		}
	}
	return lines
}
