package cmd

import (
	"github.com/guionardo/todo-cli/pkg/shell"
	"github.com/urfave/cli/v2"
)

func GetCommandSetupShell() *cli.Command {
	return &cli.Command{
		Name:   "shell",
		Usage:  "Setup shell integration",
		Action: ActionSetupShell,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "init",
				Usage: "Initialize shell integration",
			},
			&cli.BoolFlag{
				Name:  "ps1",
				Usage: "Setup PS1 prompt_command integration",
			},
			&cli.BoolFlag{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "Remove shell integration",
			},
		},
	}
}

func ActionSetupShell(c *cli.Context) error {
	if c.Bool("init") {
		return shell.SetupShellInit(c.Bool("remove"))
	}
	if c.Bool("ps1") {
		return shell.SetupShellPS1(c.Bool("remove"))
	}
	return cli.Exit("--init or --ps1 is required", 1)
}
