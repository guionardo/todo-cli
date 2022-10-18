package main

import (
	"os"

	"github.com/guionardo/todo-cli/cmd"
	"github.com/guionardo/todo-cli/pkg/logger"
)

func main() {
	if err := cmd.App().Run(os.Args); err != nil {
		logger.Fatalf("%v", err)
	}
}
