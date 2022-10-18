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
		Before:    ctx.ChainedContext(ctx.AssertLocalConfig, ctx.AssertAutoSychronization),
		Action:    ActionAct,
		ArgsUsage: "[todo-id]",
		Category:  "Tasks",
		After:     ctx.ChainedContext(ctx.AssertSave, ctx.AssertAutoSychronization),
	}
)

func ActionAct(c *cli.Context) error {
	context := ctx.ContextFromCtx(c)

	id, err := getToDoId(c)
	if err != nil {
		return err
	}

	return context.SetExit(context.Collection.DoAct(id), "Set action to #%d", id)
}
