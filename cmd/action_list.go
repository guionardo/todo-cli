package cmd

import (
	"fmt"

	"github.com/guionardo/todo-cli/internal"
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
		Flags: []cli.Flag{
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
		},
		Action: ActionList,
	}
)

func ActionList(c *cli.Context) error {
	context := internal.GetRunningContext(c).AssertExist()
	justDone := c.Bool("done")
	justPending := c.Bool("pending")

	tags := c.StringSlice("tags")
	items := context.Collection.GetByFilter(tags, justDone, justPending)
	fmt.Printf("%s\n", context.Collection.Config.ToDoListName)
	for _, item := range items {
		fmt.Printf("%s\n", item)
	}
	return nil
}
