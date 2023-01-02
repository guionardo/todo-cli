package notify

import (
	"github.com/0xAX/notificator"
	"github.com/guionardo/todo-cli/pkg/todo"
)

type DesktopNotifier struct {
}

func (n *DesktopNotifier) BeforeNotify() {
	// Do nothing
}

func (n *DesktopNotifier) Notify(title string, message string) {
	notify := notificator.New(notificator.Options{
		AppName:     "Todo CLI",
		DefaultIcon: notificator.UR_NORMAL,
	})
	notify.Push("title", message, "", notificator.UR_NORMAL)
}

func (n *DesktopNotifier) NotifyItem(item todo.Item) {
	n.Notify("ToDo Cli", item.NotifyText()+"\n"+item.StringNoColor())
}

func (n *DesktopNotifier) Info() NotifierInfo {
	return NotifierInfo{
		Name:        "desktop",
		Description: "Desktop notifications",
		Alias:       "d",
		Enabled:     true,
	}
}
func init() {
	Factory.Register("desktop", &DesktopNotifier{})
}
