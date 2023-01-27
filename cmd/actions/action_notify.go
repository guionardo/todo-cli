package actions

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
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
		Usage:   "Notify about pending tasks",
		Action:  ActionNotify,
		Flags:   flags,
		Before:  ctx.ChainedContext(ctx.LocalConfigRequired),
	}
}

func ActionNotify(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	items := c2.Collection.GetByFilter(nil, false, true)

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
