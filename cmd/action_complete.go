package cmd

import (
	"fmt"
	"strconv"

	"github.com/guionardo/todo-cli/internal"
	"github.com/urfave/cli/v2"
)

var (
	CompleteCommand = &cli.Command{
		Name:      "complete",
		Usage:     "Complete a todo item",
		Aliases:   []string{"c"},
		Action:    ActionComplete,
		Category:  "Tasks",
		ArgsUsage: "[todo-id]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "undo",
				Aliases: []string{"u"},
				Usage:   "Undo a completed todo item",
			},
		},
	}
)

func ActionComplete(c *cli.Context) error {
	context := internal.GetRunningContext(c).AssertExist()
	if c.NArg() == 0 {
		return fmt.Errorf("Missing todo-id")
	}
	todoId := c.Args().Get(0)
	id, err := strconv.Atoi(todoId)
	if err != nil {
		return fmt.Errorf("Invalid todo-id")
	}
	undo := c.Bool("undo")
	msg := "Completed"
	if undo {
		err = context.Collection.UndoComplete(id)
	} else {
		err = context.Collection.Complete(id)
		msg = "Undo completed"
	}
	if err == nil {
		err = context.Collection.Save(context.CollectionFileName)
		if err == nil {
			context.Collection.GISTSync(context.DebugMode)
			fmt.Printf("%s todo: %d", msg, id)
		}
	}
	return err
}
