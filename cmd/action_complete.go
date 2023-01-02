package cmd

import (
	"fmt"
	"time"

	"github.com/guionardo/todo-cli/pkg/ctx"
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
		},
		Before: ctx.ChainedContext(ctx.AssertLocalConfig, ctx.AssertValidId),
		After:  ctx.ChainedContext(ctx.AssertSave, ctx.AssertAutoSychronization),
	}
)

func ActionComplete(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	undo := c.Bool("undo")
	if c2.CurrentToDo.Completed == undo {
		c2.CurrentToDo.Completed = !undo
		c2.CurrentToDo.UpdatedAt = time.Now()
		return nil
	}
	c2.CancelSaving = true
	c2.CancelSync = true
	return fmt.Errorf("to-do item #%d no change", c2.CurrentToDo.Index)
}
