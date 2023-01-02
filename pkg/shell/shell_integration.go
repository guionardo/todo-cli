package shell

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/guionardo/todo-cli/pkg/utils"
)

const (
	shellIntegrationCmd = `eval "$(todo-cli init bash)"`
	shellIntegrationPs1 = `eval "$(todo-cli init ps1)"`
)

var (
	//go:embed action_setup_ps1.sh
	action_setup_ps1 string
)

func setupShellCmd(remove bool, command string) error {
	shellProfile, _ := utils.GetShellData()
	content, err := os.ReadFile(shellProfile)
	if err != nil {
		return fmt.Errorf("error reading shell profile: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, command) {
			if strings.HasPrefix(line, "#") {
				if remove {
					logger.Infof("Shell integration already removed")
				} else {
					logger.Infof("Shell integration redefined")
					lines[i] = command
				}
			} else {
				if remove {
					logger.Infof("Shell integration removed")
					lines[i] = "# " + line
				} else {
					logger.Infof("Shell integration already set")
				}
			}
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, "", command)
	}
	content = []byte(strings.Join(lines, "\n"))
	err = os.WriteFile(shellProfile, content, 0644)
	if err == nil {
		logger.Infof("Shell integration saved to: %s", shellProfile)
	} else {
		logger.Warnf("Error setting shell integration: %v", err)
	}
	return err
}
func SetupShellInit(remove bool) error {
	return setupShellCmd(remove, shellIntegrationCmd)
}

func SetupShellPS1(remove bool) error {
	return setupShellCmd(remove, shellIntegrationPs1)
}
