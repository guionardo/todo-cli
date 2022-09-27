package main

import (
	"log"
	"os"

	"github.com/guionardo/todo-cli/cmd"
)

func main() {
	if err := cmd.App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
