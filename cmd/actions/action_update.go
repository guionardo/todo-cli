package actions

import (
	"fmt"
	"time"

	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/urfave/cli/v2"
)

var (
	UpdateCommand = &cli.Command{
		Name:     "update",
		Usage:    "Update a todo item",
		Aliases:  []string{"u"},
		Action:   ActionUpdate,
		Category: "Tasks",
		Before:   ctx.ChainedContext(ctx.LocalConfigRequired, ctx.RequiredTodoId),
		Flags: append([]cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the task",
				Required: true,
			}},
			ItemFlags...,
		),
		After: ctx.ChainedContext(ctx.AssertSave, ctx.AssertSychronization),
	}
)

func ActionUpdate(c *cli.Context) error {
	context := ctx.ContextFromCtx(c)
	changed := false
	if c.IsSet("title") {
		title := c.String("title")
		if context.CurrentToDo.Title != title {
			context.CurrentToDo.Title = title
			changed = true
		}
	}

	if c.IsSet("due-date") {
		dueDate := c.Timestamp("due-date")
		if context.CurrentToDo.DueTo != *dueDate {
			context.CurrentToDo.DueTo = *dueDate
			changed = true
		}
	}

	if c.IsSet("tags") {
		context.CurrentToDo.SetTags(c.StringSlice("tags"))
		changed = true
	}

	if c.IsSet("parent-id") {
		pId := c.Int("parent-id")
		if pId > 0 {
			parentItem := context.Collection.Get(pId)
			if parentItem == nil {
				return fmt.Errorf("parent ID %d not found", pId)
			}
			if parentItem.Id != context.CurrentToDo.ParentId {
				context.CurrentToDo.ParentId = parentItem.Id
				changed = true
			}
		}
	}
	if !changed {
		return exceptions.NoChangedTodoError(context.CurrentToDo.Index)
	}
	context.CurrentToDo.UpdatedAt = time.Now()

	return nil

}
