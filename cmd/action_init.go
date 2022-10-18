package cmd

import (
	"fmt"
	"strings"

	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

var (
	InitCommand = &cli.Command{
		Name:  "init",
		Usage: "Outputs the shell initialization script",

		Category: "Tasks",
		Subcommands: []*cli.Command{
			{
				Name:   "bash",
				Action: ActionInitBash},
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
	//TODO: Usar o padrão init bash: eval "$(starship init bash)"

	_, thisPath := utils.GetShellData()
	fmt.Printf("%s notify", thisPath)
	return nil
}

func ActionInitPs1(c *cli.Context) error {
	_, thisPath := utils.GetShellData()
	ps1 := `
if [[ "$PROMPT_COMMAND" == *__tdc_ps1* ]]; then
	exit 0
fi

TCC_TIME=0

__tdc_ps1() {	
    local tcc_elapsed=$((SECONDS - TCC_TIME))
    if [[ "$TCC_TIME" -eq "0" || "$tcc_elapsed" -gt "#PERIOD#" ]]; then
        TCC_TIME=$SECONDS		
        #THISPATH# notify
    fi
}

PROMPT_COMMAND="__tdc_ps1 $PROMPT_COMMAND"
echo "todo-cli ps1 initialized (every #PERIOD# seconds)"
`
	ps1 = strings.ReplaceAll(ps1, "#PERIOD#", fmt.Sprintf("%d", c.Int("period")))
	ps1 = strings.ReplaceAll(ps1, "#THISPATH#", thisPath)
	fmt.Print(ps1)
	return nil
}
