package actions

import (
	"github.com/guionardo/todo-cli/pkg/config"
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

var (
	currentConfig = &config.LocalConfig{
		Backup: config.BackupConfig{},
	}
)

func GetCommandSetupBackup() *cli.Command {
	return &cli.Command{
		Name:   "backup",
		Usage:  "Setup of backup",
		Action: ActionSetupBackup,
		Before: ctx.ChainedContext(ctx.LocalConfigRequired, getCurrentConfig),

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Destination: &currentConfig.Backup.BackupFolder,
				Usage:       "Path to backup directory",
			},
			&cli.IntFlag{
				Name:    "needed_warning_days",
				Usage:   "Number of days before the backup is considered old",
				Value:   7,
			},
			&cli.IntFlag{
				Name:    "max_count",
				Value:   10,
			},
			&cli.BoolFlag{
				Name:    "auto_backup",
				Value:   true,
				Usage:   "Set to false to disable automatic backup",
			},
			&cli.IntFlag{
				Name:    "auto_backup_interval_days",
				Value:   1,
				Usage:   "Number of days between automatic backups",
			},
		},
	}
}

func getCurrentConfig(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)
	currentConfig = c2.LocalConfig
	return nil
}

func ActionSetupBackup(c *cli.Context) error {
	c2 := ctx.ContextFromCtx(c)

	setupChanged := false
	if c.IsSet("needed_warning_days") {
		neededWarningDays := c.Int("needed_warning_days")
		if neededWarningDays != currentConfig.Backup.BackupNeededWarningDays && neededWarningDays > 0 {
			currentConfig.Backup.BackupNeededWarningDays = neededWarningDays
			setupChanged = true
		}
	}
	if c.IsSet("max_count") {
		maxCount := c.Int("max_count")
		if maxCount != currentConfig.Backup.BackupMaxCount && maxCount >= 0 {
			currentConfig.Backup.BackupMaxCount = maxCount
			setupChanged = true
		}
	}
	if c.IsSet("auto_backup") {
		autoBackup := c.Bool("auto_backup")
		if autoBackup != currentConfig.Backup.AutoBackup {
			currentConfig.Backup.AutoBackup = autoBackup
			setupChanged = true
		}
	}
	if c.IsSet("auto_backup_interval_days") {
		autoBackupIntervalDays := c.Int("auto_backup_interval_days")
		if autoBackupIntervalDays != currentConfig.Backup.AutoBackupIntervalDays && autoBackupIntervalDays > 0 {
			currentConfig.Backup.AutoBackupIntervalDays = autoBackupIntervalDays
			setupChanged = true
		}
	}
	var err error
	if setupChanged {
		err = currentConfig.Save(c2.LocalConfigFile)
		if err == nil {
			logger.Infof("Backup setup changed: %v", currentConfig.Backup)
		}
	}
	return err

}
