package actions

import (
	"github.com/guionardo/todo-cli/pkg/ctx"

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
		Before:   ctx.ChainedContext(ctx.LocalConfigRequired),
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
	logs, err := c2.GistSync()
	for _, log := range logs {
		logger.Infof(log)
	}
	return err

}
