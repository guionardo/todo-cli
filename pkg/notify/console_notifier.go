package notify

import (
	"fmt"

	"github.com/guionardo/todo-cli/internal"
)

type ConsoleNotifier struct {
}

func (n *ConsoleNotifier) BeforeNotify() {
	fmt.Println()
}
func (n *ConsoleNotifier) Notify(title string, message string) {
	fmt.Printf("%s\n", message)
}
func (n *ConsoleNotifier) NotifyItem(item internal.ToDoItem) {
	fmt.Println(item.String() + " " + item.NotifyText())
}
func (n *ConsoleNotifier) Info() NotifierInfo {
	return NotifierInfo{
		Name:        "console",
		Description: "Prints the notification to the stdout",
		Alias:       "c",
		Enabled:     true,
		Default:     true,
	}
}

func init() {
	Factory.Register("console", &ConsoleNotifier{})
}
