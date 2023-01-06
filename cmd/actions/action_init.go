package actions

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

//go:embed action_init.sh
var ActionInitSh string

var (
	InitCommand = &cli.Command{
		Name:  "init",
		Usage: "Outputs the shell initialization script",

		Category: "Tasks",
		Subcommands: []*cli.Command{
			{
				Name:   "bash",
				Action: ActionInitBash,
			},
			{
				Name:   "ps1",
				Action: ActionInitPs1,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "period",
						Usage: "Period in seconds to show notifications",
						Value: 60,
					},
				},
			},
		},
	}
)

func ActionInitBash(c *cli.Context) error {
	_, thisPath := utils.GetShellData()
	fmt.Printf("%s list", thisPath)
	return nil
}

func ActionInitPs1(c *cli.Context) error {
	_, thisPath := utils.GetShellData()
	ps1 := strings.ReplaceAll(ActionInitSh, "#PERIOD#", fmt.Sprintf("%d", c.Int("period")))
	ps1 = strings.ReplaceAll(ps1, "#THISPATH#", thisPath)
	fmt.Print(ps1)
	return nil
}
