package cmd

import (
	"context"

	"github.com/guionardo/go-gstools/gist"
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name:    AppName,
		Version: Version,
		Usage:   "A simple todo app with a cli interface and github gist persistence",
		Authors: []*cli.Author{
			{
				Name:  "Guionardo Furlan",
				Email: "guionardo@gmail.com",
			},
		},
		ExtraInfo: func() map[string]string {
			return map[string]string{
				"Build Host": BuildHost,
				"Build Date": BuildDate,
			}
		},
		EnableBashCompletion: true,
		Before: func(c *cli.Context) error {
			logger.SetLogger(c.Bool("debug"))
			if c.Bool("debug") {
				logger.Debugf("Debug mode enabled")
				gist.SetDefaultLogger()
			}
			c.Context = context.WithValue(c.Context, ctx.Key, ctx.ContextFromCli(c))
			return nil
		},
		Suggest: true,
		Action:  ActionList,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "data-folder",
				Usage:   "Load configuration from `FOLDER`",
				Value:   ctx.GetDefaultDataFolder(),
				EnvVars: []string{"TODO_CONFIG"},
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},
		},
		Commands: []*cli.Command{
			AddCommand,
			UpdateCommand,
			ListCommand,
			DeleteCommand,
			CompleteCommand,
			ActCommand,
			SetupCommand,
			NotifyCommand,
			SyncCommand,
			BackupCommand,
			InitCommand,
		},
	}
}
