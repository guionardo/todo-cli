package cmd

import (
	"context"

	"github.com/guionardo/todo-cli/internal"
	"github.com/urfave/cli/v2"
)

const Version = "0.0.1"

func App() *cli.App {
	extra_info := map[string]string{
		"github": "https://github.com/guionardo/todo-cli"}
	return &cli.App{
		Name:  "todo-cli",
		Usage: "A simple todo app with a cli interface and github gist persistence",
		Authors: []*cli.Author{
			{
				Name:  "Guionardo Furlan",
				Email: "guionardo@gmail.com",
			},
		},
		EnableBashCompletion: true,
		Version:              Version,
		ExtraInfo: func() map[string]string {
			return extra_info
		},
		Before: func(c *cli.Context) error {
			running_context := internal.NewRunningContext(c)
			c.Context = context.WithValue(c.Context, "running_context", running_context)
			return nil
		},
		Action: ActionList,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				Value:   internal.DefaultCollectionFilePath,
				EnvVars: []string{"TODO_CONFIG"},
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},
		},
		Commands: []*cli.Command{
			AddCommand,
			ListCommand,
			DeleteCommand,
			CompleteCommand,
			ActCommand,
			SetupCommand,
			NotifyCommand,
		},
	}
}

func GetCollection(c *cli.Context) *internal.ToDoCollection {
	return c.Context.Value("collection").(*internal.ToDoCollection)
}
