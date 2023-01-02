package cmd

import (
	"fmt"

	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

func GetCommandSetupSync() *cli.Command {
	return &cli.Command{
		Name:   "sync",
		Usage:  "Sync with gist",
		Before: ctx.ChainedContext(ctx.AssertLocalConfig),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "auto-sync",
				Aliases: []string{"a"},
				Usage:   "Enable auto sync on every change",
				Value:   true,
			},
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "Set Github token. Create a new github token at https://github.com/settings/tokens/new with gist permission",
			},
		},
		Action: ActionSetupSync,
	}
}
func ActionSetupSync(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	changed := false
	if c.IsSet("token") {
		token := c.String("token")
		if len(token) > 0 && token != c2.LocalConfig.Gist.Authorization {
			err := c2.LocalConfig.Gist.SetToken(token)
			if err != nil {
				return err
			}
			changed = true
		}
	}
	if c.IsSet("auto-sync") {
		autoSync := c.Bool("auto-sync")
		if c2.LocalConfig.Gist.AutoSync != autoSync {
			if autoSync && len(c2.LocalConfig.Gist.Authorization) == 0 {
				return fmt.Errorf("cannot enable auto-sync without a valid token")
			}
			c2.LocalConfig.Gist.AutoSync = autoSync
			changed = true
		}
	}
	if !changed {
		logger.Infof("Nothing to do")
		return nil
	}
	err := c2.LocalConfig.Save(c2.LocalConfigFile)
	if err == nil {
		logger.Infof("Config file %s saved - %v", c2.LocalConfigFile, c2.LocalConfig.Gist)
	} else {
		logger.Warnf("Error saving config file %s: %v", c2.LocalConfigFile, err)
		return err
	}
	if c2.LocalConfig.Gist.AutoSync {
		log, err := c2.GistSync()

		if err != nil {
			logger.Warnf("Error syncing with gist: %v", err)
			return err
		}

		for _, l := range log {
			logger.Infof("  %s", l)
		}

	}

	return nil
}
