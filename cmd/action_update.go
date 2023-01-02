package cmd

import (
	"fmt"
	"time"

	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/urfave/cli/v2"
)

var (
	UpdateCommand = &cli.Command{
		Name:     "update",
		Usage:    "Update a todo item",
		Aliases:  []string{"u"},
		Action:   ActionUpdate,
		Category: "Tasks",
		Before:   ctx.ChainedContext(ctx.AssertLocalConfig, ctx.AssertValidId),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the task",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "title",
				Usage:    "Title of the task",
				Required: false,
			},
			&cli.TimestampFlag{
				Name:     "due-date",
				Aliases:  []string{"d"},
				Usage:    "Due date for the todo item",
				Layout:   "2006-01-02",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "tags",
				Aliases:  []string{"t"},
				Usage:    "Tags for the todo item",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "parent-id",
				Usage:    "Parent ID for the todo item",
				Required: false,
				Aliases:  []string{"p"},
			},
		},
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
		context.CurrentToDo.Tags = c.StringSlice("tags")
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
		context.CancelSaving = true
		context.CancelSync = true
		return fmt.Errorf("nothing to change")
	}
	context.CurrentToDo.UpdatedAt = time.Now()

	return nil

}
