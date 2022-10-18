package cmd

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/guionardo/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

var (
	ListAll     = false
	ListDone    = false
	ListCommand = &cli.Command{
		Name:     "list",
		Aliases:  []string{"l"},
		Usage:    "List all todo items",
		Category: "Tasks",
		Before:   ctx.ChainedContext(ctx.AssertLocalConfig, ctx.AssertAutoSychronization, ctx.OptionalId),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "List a specific item",
				Required: false,
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "List all tasks",
			},
			&cli.BoolFlag{
				Name:    "done",
				Aliases: []string{"d"},
				Usage:   "List all done tasks",
			},
			&cli.BoolFlag{
				Name:    "pending",
				Aliases: []string{"p"},
				Usage:   "List all pending tasks",
			},
			&cli.StringSliceFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "List all tasks with the specified tags",
			},
			&cli.BoolFlag{
				Name:  "no-header",
				Usage: "Do not print the header",
			},
			//TODO: Implementar comando para listar campos específicos. Ler campos por reflection: https://stackoverflow.com/questions/40864840/how-to-get-the-json-field-names-of-a-struct-in-golang
		},
		Action: ActionList,
	}
)

func ActionList(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	justDone := c.Bool("done")
	justPending := c.Bool("pending")

	var items []*todo.ToDoItem

	if c2.CurrentToDo != nil {
		items = []*todo.ToDoItem{c2.CurrentToDo}
	} else {
		items = c2.Collection.GetByFilter(c.StringSlice("tags"), justDone, justPending)
	}
	if !c.Bool("no-header") {
		logger.Logf("%s\n", c2.LocalConfig.ToDoListName)
	}
	for _, row := range c2.Collection.GetTreeList(items) {
		logger.Logf("%s", row)
	}
	return nil
}
