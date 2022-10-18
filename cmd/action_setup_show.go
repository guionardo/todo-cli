package cmd

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

func GetCommandSetupShow() *cli.Command {
	return &cli.Command{
		Name:   "show",
		Usage:  "Show current setup",
		Action: ActionSetupShow,
		Before: ctx.ChainedContext(ctx.AssertLocalConfig),
	}
}

func ActionSetupShow(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	logger.Infof(AppDescription)	
	logger.Infof("Build: %s @ %s", BuildHost, BuildDate)
	logger.Infof("Current config file: %s", c2.LocalConfigFile)
	logger.Infof("Current to-do collection file: %s", c2.LocalCollectionFile)
	logger.Infof("%v", c2.LocalConfig)

	return nil
}
