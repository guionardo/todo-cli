package actions

import (
	"fmt"

	"github.com/guionardo/todo-cli/pkg/backup"
	"github.com/guionardo/todo-cli/pkg/config"
	"github.com/guionardo/todo-cli/pkg/ctx"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/guionardo/todo-cli/pkg/todo"
	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

func GetCommandSetupNew() *cli.Command {
	return &cli.Command{
		Name:   "new",
		Usage:  "Create a new setup",
		Action: ActionSetupNew,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Name of the todo list",
				Value:   fmt.Sprintf("%s's TODO", utils.GetUser()),
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force creation of new config file",
			},
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "Set Github token. Create a new github token at https://github.com/settings/tokens/new with gist permission",
				Required: false,
			},
		},
	}
}

func checkFileExists(localFile string, backupFolder string, force bool, fileTitle string, backupCfg config.BackupConfig) error {
	if err := utils.FileExists(localFile); err != nil {
		return nil
	}

	if !force {
		return fmt.Errorf("%s %s already exists. Use --force to overwrite", fileTitle, localFile)
	}
	logger.Warnf("%s %s already exists. Overwriting", fileTitle, localFile)

	// cfg,err:=ctx.LoadLocalConfig(localFile)
	// if err != nil {
	// 	logger.Warnf("Error loading %s %s: %v", fileTitle, localFile, err)
	// 	return err
	// }

	bkp, err := backup.CreateBackup(localFile, backupFolder, backupCfg)
	if err != nil {
		return err
	}
	backupFile, err := bkp.DoBackup()
	if err != nil && len(backupFile) == 0 {
		return fmt.Errorf("error backing up %s %s: %v", fileTitle, localFile, err)
	}
	logger.Infof("Backup file %s -> %s", localFile, backupFile)
	return nil
}
func checkConfigExists(context *ctx.Context, force bool) error {
	if err := checkFileExists(context.LocalConfigFile, context.LocalConfig.Backup.BackupFolder, force, "Config", context.LocalConfig.Backup); err != nil {
		return err
	}
	if err := checkFileExists(context.LocalCollectionFile, context.LocalConfig.Backup.BackupFolder, force, "Todo", context.LocalConfig.Backup); err != nil {
		return err
	}

	return nil
}

func ActionSetupNew(c *cli.Context) (err error) {
	context := ctx.ContextFromCli(c)

	if err = checkConfigExists(context, c.Bool("force")); err != nil {
		return
	}

	name := c.String("name")
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}

	config := config.LocalConfig{
		ToDoListName: name,
		Gist:         config.GetDefaultGistConfig(),
		Backup:       config.GetDefaultBackupConfig(context.DataFolder),
	}
	if err = config.Save(context.LocalConfigFile); err != nil {
		return
	}
	logger.Infof("Config file %s created", context.LocalConfigFile)

	collection := todo.NewTodoCollection()
	if err = collection.Save(context.LocalCollectionFile); err != nil {
		return
	}
	logger.Infof("Todo file %s created", context.LocalCollectionFile)

	token := c.String("token")
	if len(token) == 0 {
		return
	}

	if err = config.Gist.SetToken(token); err != nil {
		logger.Warnf("error setting token: %v", err)
	}
	err = config.Save(context.LocalConfigFile)

	return
}
