package actions

import (
	"time"

	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/guionardo/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

var (
	AddCommand = &cli.Command{
		Name:    "add",
		Usage:   "Add a new todo item",
		Aliases: []string{"a"},
		Before: ctx.ChainedContext(
			ctx.LocalConfigRequired,
			ctx.AssertAutoSychronization),
		Action:   ActionAdd,
		Category: "Tasks",
		Flags:    ItemFlags,
		After:    ctx.ChainedContext(ctx.AssertSave, ctx.AssertSychronization),
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
				err := exceptions.ParentTodoNotFoundError(pId)
				return context.SetExitError(err, "%v", err)
			}
			parentId = parentItem.Id
		}
	}

	tags := c.StringSlice("tags")
	item := &todo.Item{
		Id:        todo.NewItemId(),
		Title:     title,
		DueTo:     dueDate,
		UpdatedAt: time.Now(),
		ParentId:  parentId,
	}
	item.SetTags(tags)
	context.Collection.Add(item)

	context.Collection.LastUpdate = time.Now()
	return context.SetExitSuccess("Added to-do %s", item)
}
