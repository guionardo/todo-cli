package utils

import (
	"os"
	"path"

	"github.com/guionardo/go-gstools/git"
	"github.com/guionardo/todo-cli/pkg/logger"
)

// GetUser get username from git config, from OS or UNDEFINED
func GetUser() string {
	gitUser, err := git.GetCurrentGitUser()
	if err == nil && len(gitUser.Name) > 0 {
		return gitUser.Name
	}
	user := os.Getenv("USER")
	if len(user) == 0 {
		user = "UNDEFINED"
	}
	return user
}

func Tern(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func GetShellData() (shellProfile string, thisPath string) {
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Fatalf("Error getting user home dir: %v", err)
	}
	shellProfile = path.Join(home, ".bashrc")
	if _, err := os.Stat(shellProfile); os.IsNotExist(err) {
		shellProfile = ""
	}
	thisPath, err = os.Executable()
	if err != nil {
		logger.Fatalf("Error getting executable path: %v", err)
	}
	return
}
