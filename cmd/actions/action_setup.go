package actions

import (
	"github.com/urfave/cli/v2"
)

var (
	SetupCommand *cli.Command
)

func init() {

	SetupCommand = &cli.Command{
		Name:     "setup",
		Usage:    "Setup todo app",
		Category: "Setup",
		Subcommands: []*cli.Command{
			GetCommandSetupShow(),
			GetCommandSetupNew(),
			GetCommandSetupSync(),
			GetCommandSetupBackup(),
			GetCommandSetupShell(),
		},
	}
}
