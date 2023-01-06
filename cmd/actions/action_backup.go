package actions

import (
	"path"
	"strings"

	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)



var (
	BackupCommand = &cli.Command{
		Name:     "backup",
		Usage:    "Run a backup of the todo list",
		Before:   ctx.ChainedContext(ctx.LocalConfigRequired, ctx.AssertAutoSychronization),
		Category: "Tasks",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Path to backup directory",
				Required: false,
				Value:    path.Join(utils.GetDefaultDataFolder(), "backup"),
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run a backup of the todo list",
				Action: ActionBackupRun,
			},
			{
				Name:   "list",
				Usage:  "List all backups",
				Action: ActionBackupList,
			},
		},
	}
)

func ActionBackupRun(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	bkpConfig := c2.LocalConfig.Backup
	bkpConfig.BackupFolder = c.String("path")
	if err := c2.Collection.Backup(&bkpConfig); err != nil {
		if err == exceptions.NoNeedBackupError {
			return c2.SetExitWarning(err.Error())
		}
		return c2.SetExitError(err, "Error running backup - %v", err)
	}

	return c2.SetExitSuccess("Backup completed")
}

func ActionBackupList(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	bkpConfig := c2.LocalConfig.Backup
	bkpConfig.BackupFolder = c.String("path")
	backupFiles, err := c2.Collection.ListBackup(&bkpConfig)

	if err != nil {
		return c2.SetExitError(err, "Error listing backups - %v", err)
	}
	
	return c2.SetExitSuccess("Backups:\n" + strings.Join(backupFiles, "\n"))
}
