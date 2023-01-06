package utils

import (
	"fmt"
	"os"
	"path"
)

func GetDefaultDataFolder() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("error getting user home dir: %s", err))
	}
	return path.Join(home, ".config", "todo-cli")
}