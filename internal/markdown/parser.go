package markdown

import (
	"fmt"
	"time"

	"github.com/guionardo/todo-cli/internal"
)

func CreateMarkdown(config internal.Config, items []internal.ToDoItem, updated time.Time) (output []string, err error) {
	if len(config.ToDoListName) == 0 {
		err = fmt.Errorf("No list name defined")
		return
	}
	output = make([]string, 4+len(items))
	output[0] = fmt.Sprintf("# %s", config.ToDoListName)
	output[2] = fmt.Sprintf("## Updated: %s", updated.Format("2006-01-02 15:04:05"))
	for i, item := range items {
		output[i+3] = item.ToMarkDown()
	}

	return nil, nil
}
