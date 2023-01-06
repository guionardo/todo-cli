package ctx

import (
	"github.com/urfave/cli/v2"
)

func ChainedContext(functions ...func(c *cli.Context) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		for _, f := range functions {
			err := f(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
