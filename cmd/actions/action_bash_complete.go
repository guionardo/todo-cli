package actions

import (
	"fmt"

	_ "embed"

	"github.com/urfave/cli/v2"
)

//go:embed action_bash_complete.sh
var bashCompleteScript string

var (
	AutoCompleteCommand = &cli.Command{
		Name:     "autocomplete",
		Usage:    "autocomplete bash script",
		Action:   ActionAutoComplete,
		Category: "Tasks",
	}
)

func ActionAutoComplete(c *cli.Context) error {
	fmt.Print(bashCompleteScript)
	return nil
}
