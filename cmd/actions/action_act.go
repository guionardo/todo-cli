package actions

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/urfave/cli/v2"
)

var (
	ActCommand = &cli.Command{
		Name:    "act",
		Usage:   "Set current timestamp as action for item",
		Aliases: []string{"a"},
		Before: ctx.ChainedContext(ctx.LocalConfigRequired,
			ctx.AssertAutoSychronization,
			ctx.RequiredTodoId),
		Action:    ActionAct,
		ArgsUsage: "[todo-id]",
		Category:  "Tasks",
		After: ctx.ChainedContext(
			ctx.AssertSave,
			ctx.AssertAutoSychronization),
	}
)

func ActionAct(c *cli.Context) error {
	context := ctx.ContextFromCtx(c)

	if err := context.Collection.DoAct(context.CurrentToDo.Index); err != nil {
		err = exceptions.NewException(err)
		return context.SetExitError(err, "Error registering act for item %d: %v", context.CurrentToDo.Index, err)
	}
	return context.SetExitSuccess("Registered act for item %d", context.CurrentToDo.Index)
}
