package cmd

import (
	"context"
	"os"
	"path"

	"time"

	"github.com/guionardo/go-gstools/gist"
	"github.com/guionardo/todo-cli/cmd/actions"
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name:    utils.AppName,
		Version: utils.Version,
		Usage:   "A simple todo app with a cli interface and github gist persistence",
		Authors: []*cli.Author{
			{
				Name:  "Guionardo Furlan",
				Email: "guionardo@gmail.com",
			},
		},
		ExtraInfo: func() map[string]string {
			return map[string]string{
				"Build Host": utils.BuildHost,
				"Build Date": utils.BuildDate,
			}
		},
		EnableBashCompletion: true,
		Before:               BeforeApp,
		After:                AfterApp,
		Suggest:              true,
		Action:               actions.ActionList,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "data-folder",
				Usage:   "Load configuration from `FOLDER`",
				Value:   utils.GetDefaultDataFolder(),
				EnvVars: []string{"TODO_CONFIG"},
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable color output",
			},
			&cli.StringFlag{
				Name:  "log-file",
				Value: path.Join(utils.GetDefaultDataFolder(), "todo_"+time.Now().Format("2006-01-02")+".log"),
				Usage: "Log file path (set to '_' to disable)",
			},
		},
		Commands: []*cli.Command{
			actions.AddCommand,
			actions.UpdateCommand,
			actions.ListCommand,
			actions.DeleteCommand,
			actions.CompleteCommand,
			actions.ActCommand,
			actions.SetupCommand,
			actions.NotifyCommand,
			actions.SyncCommand,
			actions.BackupCommand,
			actions.InitCommand,
			actions.AutoCompleteCommand,
		},
	}
}

func BeforeApp(c *cli.Context) error {
	if c.Bool("no-color") {
		logger.DisableColor()
	}
	logFile := c.String("log-file")
	if logFile != "_" {
		logger.SetFileOutput(logFile)
		logger.FileLog.Printf("LINE %s", os.Args)
	}
	logger.SetLogger(c.Bool("debug"))
	if c.Bool("debug") {
		logger.Debugf("Debug mode enabled")
		gist.SetDefaultLogger()
	}
	c.Context = context.WithValue(c.Context, ctx.Key, ctx.ContextFromCli(c))
	return nil
}

func AfterApp(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	if c2.Error != nil {
		logger.Warnf(utils.Tern(len(c2.ErrorPrefix) > 0, c2.ErrorPrefix, "Error: ")+"%v", c2.Error)
	} else {
		if c2.ExitWarning {
			logger.Warnf(c2.ExitMessage)
		} else {
			logger.Infof(c2.ExitMessage)
		}
	}
	logger.CloseLogging()
	return nil
}
