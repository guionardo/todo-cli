package cmd

import (
	"fmt"
	"path"

	"github.com/guionardo/todo-cli/pkg/backup"
	"github.com/guionardo/todo-cli/pkg/consts"
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

var (
	BackupCommand = &cli.Command{
		Name:     "backup",
		Usage:    "Run a backup of the todo list",
		Aliases:  []string{"a"},
		Before:   ctx.ChainedContext(ctx.AssertLocalConfig, ctx.AssertAutoSychronization),
		Category: "Tasks",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Path to backup directory",
				Required: false,
				Value:    path.Join(ctx.GetDefaultDataFolder(), "backup"),
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
	bkpColl, err := backup.NewBackupCollection(c2.DataFolder, c2.LocalConfig.Backup, consts.DefaultLocalConfigFile, consts.DefaultLocalCollectionFile)
	if err != nil {
		return err
	}
	if err = bkpColl.DoBackup(); err == nil {
		logger.Infof("Backup created") // TODO: Implementar log dos arquivos que foram backupados
	}

	return err
}

func ActionBackupList(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	bkpConfig := c2.LocalConfig.Backup
	bkpConfig.BackupFolder = c.String("path")
	bkpColl, err := backup.NewBackupCollection(c2.DataFolder, c2.LocalConfig.Backup, consts.DefaultLocalConfigFile, consts.DefaultLocalCollectionFile)
	if err != nil {
		return err
	}
	for sourceFile, bkp := range bkpColl.Backups {
		fmt.Printf("%s: %s\n", sourceFile, bkp.LastBackups)
	}

	return nil
}
