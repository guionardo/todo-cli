package cmd

import (
	"context"
	"fmt"

	"github.com/guionardo/todo-cli/internal"
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
		Before: func(c *cli.Context) error {
			collectionFile, err := internal.CollectionFile()
			if err != nil {
				return err
			}
			collection, err := internal.ParseCollectionFile(collectionFile)
			c.Context = context.WithValue(c.Context, "collection", collection)
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Ok")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",				
				Destination: &internal.CollectionConfigFile,
				EnvVars:     []string{"TODO_CONFIG"},
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug mode",
				Destination: &internal.DebugMode,
			},
		},
		Commands: []*cli.Command{
			AddCommand,
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
			{
				Name:    "setup",
				Usage:   "Setup todo app",
				Aliases: []string{"s"},
				Action:  ActionSetup,
			},
		},
	}
}

func GetCollection(c *cli.Context) *internal.ToDoCollection {
	return c.Context.Value("collection").(*internal.ToDoCollection)
}
