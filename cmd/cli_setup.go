package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func App() *cli.App {

	return &cli.App{
		Name:  "todo",
		Usage: "A simple todo app with a cli interface and github gist persistence",
		Authors: []*cli.Author{
			{
				Name:  "Guionardo Furlan",
				Email: "guionardo@gmail.com",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Ok")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Usage:   "Add a new todo item",
				Aliases: []string{"a"},
				Action:  ActionAdd,
			},
			{
				Name:    "list",
				Usage:   "List all todo items",
				Aliases: []string{"l"},
				Action:  ActionList,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Aliases:     []string{"a"},
						Usage:       "List all todo items",
						Destination: &ListAll,
					},
					&cli.BoolFlag{
						Name:        "done",
						Aliases:     []string{"d"},
						Usage:       "List all done todo items",
						Destination: &ListDone,
					},
				},
			},
			{
				Name:    "delete",
				Usage:   "Delete a todo item",
				Aliases: []string{"d"},
				Action:  ActionDelete,
			},
			{
				Name:      "complete",
				Usage:     "Complete a todo item",
				Aliases:   []string{"c"},
				Action:    ActionComplete,
				ArgsUsage: "[todo-id]",
			},
		},
	}
}
