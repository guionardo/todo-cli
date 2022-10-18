package cmd

import (
	"fmt"

	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/github"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

var (
	SyncCommand = &cli.Command{
		Name:     "sync",
		Usage:    "Synchronize local collection with GIST",
		Aliases:  []string{"s"},
		Action:   ActionSync,
		Category: "Setup",
		Before:   ctx.ChainedContext(ctx.AssertLocalConfig),
		After:    ctx.AssertSave,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "simulate",
				Usage:    "Just for testing",
				Required: false,
			},
		},
	}
)

func ActionSync(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)

	_, err := github.NewGistAPI(&c2.LocalConfig.Gist)
	if err != nil {
		return err
	}
	simulate := c.Bool("simulate")
	diffCount, log, err := c2.Collection.GistSync(&c2.LocalConfig.Gist)
	if err != nil {
		return err
	}
	if diffCount == 0 {
		return fmt.Errorf("No changes detected")
	}
	for _, line := range log {
		logger.Infof("  %s", line)
	}
	if simulate {
		c2.CancelSaving = true
		return fmt.Errorf("Simulating sync: no changes will be made")
	}
	return nil

}
