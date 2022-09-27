package cmd

import (
	"fmt"
	"strconv"

	"github.com/guionardo/todo-cli/internal"
	"github.com/urfave/cli/v2"
)

var (
	ActCommand = &cli.Command{
		Name:      "act",
		Usage:     "Set current timestamp as action for item",
		Aliases:   []string{"a"},
		Action:    ActionAct,
		ArgsUsage: "[todo-id]",
		Category:  "Tasks",
	}
)

func ActionAct(c *cli.Context) error {
	context := internal.GetRunningContext(c).AssertExist()
	if c.NArg() == 0 {
		return fmt.Errorf("Missing todo-id")
	}
	todoId := c.Args().Get(0)
	id, err := strconv.Atoi(todoId)
	if err != nil {
		return fmt.Errorf("Invalid todo-id")
	}

	err = context.Collection.DoAct(id)
	if err == nil {
		if err = context.Collection.Save(context.CollectionFileName); err == nil {
			context.Collection.GISTSync(context.DebugMode)
		}
	}
	return err
}
