package actions

import (
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func GetCommandSetupShow() *cli.Command {
	return &cli.Command{
		Name:   "show",
		Usage:  "Show current setup",
		Action: ActionSetupShow,
		Before: ctx.ChainedContext(ctx.LocalConfigRequired),
	}
}

func ActionSetupShow(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	if setup, err := yaml.Marshal(c2.LocalConfig); err != nil {
		return err
	} else {
		logger.Infof(utils.AppDescription)
		logger.Infof("Build: %s @ %s\n", utils.BuildHost, utils.BuildDate)
		logger.Infof("Current config file: %s", c2.LocalConfigFile)
		logger.Infof("Current to-do collection file: %s", c2.LocalCollectionFile)
		logger.Infof("Current setup [%s]:\n%s", utils.AppDescription, string(setup))
	}

	return nil
}
