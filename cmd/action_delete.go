package cmd

import "github.com/urfave/cli/v2"

var (
	DeleteCommand = &cli.Command{
		Name:     "delete",
		Usage:    "Delete a todo item",
		Aliases:  []string{"d"},
		Action:   ActionDelete,
		Category: "Tasks",
	}
)

func ActionDelete(c *cli.Context) error {
	return nil
}
