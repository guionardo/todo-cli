package notify

import "github.com/guionardo/todo-cli/pkg/todo"

type (
	Notifier interface {
		Notify(title string, message string)
		NotifyItem(item todo.ToDoItem)
		Info() NotifierInfo
		BeforeNotify()
	}
	NotifierInfo struct {
		Name        string
		Description string
		Alias       string
		Enabled     bool
		Default     bool
	}
	NotifierFactory struct {
		notifiers map[string]Notifier
	}
)

var Factory = &NotifierFactory{
	notifiers: make(map[string]Notifier),
}

func (n *NotifierFactory) Register(name string, notifier Notifier) {
	n.notifiers[name] = notifier
}
func (n *NotifierFactory) Notifiers() []Notifier {
	notifiers := make([]Notifier, len(n.notifiers))
	i := 0
	for _, notifier := range n.notifiers {
		notifiers[i] = notifier
		i++
	}
	return notifiers
}
