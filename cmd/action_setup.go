package cmd

import (
	"fmt"
	"os"

	"github.com/guionardo/todo-cli/internal"
	"github.com/urfave/cli/v2"
)

var ()

func getUser() string {
	gitUser, err := internal.GetCurrentGitUser()
	if err == nil && len(gitUser.Name) > 0 {
		return gitUser.Name
	}
	user := os.Getenv("USER")
	if len(user) == 0 {
		user = "UNDEFINED"
	}
	return user
}
func ActionSetup(c *cli.Context) error {
	fmt.Println("Starting todo setup")
	default
	return nil
}
