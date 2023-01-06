package actions

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

var (
	CompleteCommand = &cli.Command{
		Name:      "complete",
		Usage:     "Complete a todo item",
		Aliases:   []string{"c"},
		Action:    ActionComplete,
		Category:  "Tasks",
		ArgsUsage: "[todo-id]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "undo",
				Aliases: []string{"u"},
				Usage:   "Undo a completed todo item",
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "Complete all children of the todo item",
			},
		},
		Before: ctx.ChainedContext(ctx.LocalConfigRequired, ctx.RequiredTodoId),
		After:  ctx.ChainedContext(ctx.AssertSave, ctx.AssertAutoSychronization),
	}
)

func ActionComplete(c *cli.Context) (err error) {
	c2 := ctx.ContextFromCtx(c)

	if err = c2.Collection.Complete(c2.CurrentToDo.Index,
		c.Bool("undo"),
		c.Bool("recursive")); err != nil {
		err = exceptions.NewException(err)
	} else {
		logger.Infof("Completed #%d", c2.CurrentToDo.Index)
	}
	return

}
