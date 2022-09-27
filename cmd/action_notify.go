package cmd

import (
	"github.com/guionardo/todo-cli/internal"
	"github.com/guionardo/todo-cli/pkg/notify"
	"github.com/urfave/cli/v2"
)

var (
	NotifyCommand *cli.Command
)

func init() {
	notifiers := notify.Factory.Notifiers()
	flags := make([]cli.Flag, len(notifiers))
	for i, notifier := range notifiers {
		flags[i] = &cli.BoolFlag{
			Name:    notifier.Info().Name,
			Aliases: []string{notifier.Info().Alias},
			Usage:   notifier.Info().Description,
			Hidden:  !notifier.Info().Enabled,
			Value:   notifier.Info().Default,
		}
	}
	NotifyCommand = &cli.Command{
		Name:    "notify",
		Aliases: []string{"n"},
		Usage:   "Notify about pending tasks",
		Action:  ActionNotify,
		Flags:   flags,
	}
}

func ActionNotify(c *cli.Context) error {
	context := internal.GetRunningContext(c).AssertExist()
	items := context.Collection.GetByFilter(nil, false, true)
	if len(items) == 0 {
		return nil
	}
	// fmt.Printf("%s\n", context.Collection.Config.ToDoListName)
	for _, notifier := range notify.Factory.Notifiers() {
		if c.Bool(notifier.Info().Name) || c.Bool(notifier.Info().Alias) {
			notifier.BeforeNotify()
			for _, item := range items {
				notifier.NotifyItem(*item)
			}
		}
	}
	return nil
}
