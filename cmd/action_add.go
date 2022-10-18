package cmd

import (
	"fmt"
	"time"

	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

var (
	AddCommand = &cli.Command{
		Name:     "add",
		Usage:    "Add a new todo item",
		Aliases:  []string{"a"},
		Before:   ctx.ChainedContext(ctx.AssertLocalConfig, ctx.AssertAutoSychronization),
		Action:   ActionAdd,
		Category: "Tasks",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "Title of the task",
				Required: true,
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

func ActionAdd(c *cli.Context) error {
	context := ctx.ContextFromCtx(c)
	title := c.String("title")

	dueDate := time.Time{}
	if c.IsSet("due-date") {
		dueDate = *c.Timestamp("due-date")
	}
	parentId := ""
	if c.IsSet("parent-id") {
		pId := c.Int("parent-id")
		if pId > 0 {
			parentItem := context.Collection.Get(pId)
			if parentItem == nil {
				return fmt.Errorf("Parent ID %d not found", pId)
			}
			parentId = parentItem.Id
		}
	}

	tags := c.StringSlice("tags")
	item := &todo.ToDoItem{
		Id:        todo.NewItemId(),
		Title:     title,
		DueTo:     dueDate,
		Tags:      tags,
		UpdatedAt: time.Now(),
		ParentId:  parentId,
	}
	context.Collection.Add(item)

	context.Collection.LastUpdate = time.Now()
	context.SetExit(nil, "Added to-do %s", item)
	return nil
}
