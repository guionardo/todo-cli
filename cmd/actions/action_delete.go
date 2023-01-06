package actions

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

var (
	DeleteCommand = &cli.Command{
		Name:     "delete",
		Usage:    "Delete a todo item",
		Aliases:  []string{"d"},
		Action:   ActionDelete,
		Category: "Tasks",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force delete without confirmation",
			},
		},
		Before: ctx.ChainedContext(ctx.LocalConfigRequired, ctx.RequiredTodoId),
		After:  ctx.ChainedContext(ctx.AssertSave, ctx.AssertAutoSychronization),
	}
)

func ActionDelete(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)

	force := c.Bool("force")
	if !force && !utils.AskYesNo(false, "Are you sure you want to delete this item: %s ?", c2.CurrentToDo.String()) {
		return nil
	}
	if err := c2.Collection.Remove(c2.Id); err != nil {
		logger.Warnf("Error deleting item %v: %v", c2.CurrentToDo, err)
		return err
	}
	logger.Infof("Item %v deleted", c2.CurrentToDo)

	return nil

}
