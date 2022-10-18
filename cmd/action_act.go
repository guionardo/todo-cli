package cmd

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/urfave/cli/v2"
)

var (
	ActCommand = &cli.Command{
		Name:      "act",
		Usage:     "Set current timestamp as action for item",
		Aliases:   []string{"a"},
		Before:    ctx.ChainedContext(ctx.AssertLocalConfig, ctx.AssertAutoSychronization, ctx.AssertValidId),
		Action:    ActionAct,
		ArgsUsage: "[todo-id]",
		Category:  "Tasks",
		After:     ctx.ChainedContext(ctx.AssertSave, ctx.AssertAutoSychronization),
	}
)

func ActionAct(c *cli.Context) error {
	context := ctx.ContextFromCtx(c)

	return context.SetExit(context.Collection.DoAct(context.CurrentToDo.Index), "Set action to #%d", context.CurrentToDo.Index)
}
