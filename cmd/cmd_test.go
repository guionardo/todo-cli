package cmd

import (
	"path"
	"testing"
	"time"

	"github.com/guionardo/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func TestActions(t *testing.T) {
	app := App()
	setup_file := path.Join(t.TempDir(), "todo-cli.yaml")

	tests := []struct {
		name string
		ctx  *cli.Context
	}{
		{name: "setup", ctx: cli.NewContext(app, nil, nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.Run([]string{"todo-cli", "--config", setup_file, "setup new --token 123"})
			if err != nil {
				t.Errorf("App.Run() error = %v", err)
			}

		})
	}
}

func CmdTest_(t *testing.T) {

	dayAgo := time.Now().Add(-24 * time.Hour)
	tests := []struct {
		name string
		item todo.ToDoItem
		want string
	}{
		{name: "completed", item: todo.ToDoItem{Title: "Title", Completed: true, UpdatedAt: dayAgo}, want: "Completed @ " + dayAgo.Format(todo.DateTimeFormat)},
		{name: "new item", item: todo.ToDoItem{Title: "Title", UpdatedAt: dayAgo}, want: "New (1 days)"},
		{name: "due to", item: todo.ToDoItem{Title: "Title", DueTo: dayAgo, UpdatedAt: dayAgo}, want: "Overdue -1 days"},
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
