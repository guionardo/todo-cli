package actions

import (
	"strings"

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
		Before:   ctx.ChainedContext(ctx.LocalConfigRequired, ctx.AssertAutoSychronization, ctx.OptionalId),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Usage:    "List a specific item",
				Required: false,
			},
			&cli.BoolFlag{
				Name:    "all",
				Usage:   "List all tasks",
			},
			&cli.BoolFlag{
				Name:    "done",
				Usage:   "List all done tasks",
			},
			&cli.BoolFlag{
				Name:    "pending",
				Usage:   "List all pending tasks",
			},
			&cli.StringSliceFlag{
				Name:    "tags",
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

	var items []*todo.Item

	if c2.CurrentToDo != nil {
		items = []*todo.Item{c2.CurrentToDo}
	} else {
		items = c2.Collection.GetByFilter(c.StringSlice("tags"), justDone, justPending)
	}
	if !c.Bool("no-header") {
		logger.Logf("%s\n", c2.LocalConfig.ToDoListName)
	}
	for _, row := range c2.Collection.GetTreeView(items) {
		logger.Logf("%s", strings.ReplaceAll(row, "\n", ""))
	}
	return nil
}
