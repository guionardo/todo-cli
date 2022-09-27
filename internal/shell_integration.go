package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func SetupShellIntegration(c *RunningContext) error {
	return setProfileCommand()
}

func RemoveShellIntegration(c *RunningContext) error {
	return unsetProfileCommand()
}

var (
	shellProfile = ""
	thisPath     = ""
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home dir: %v", err)
	}
	shellProfile = path.Join(home, ".bashrc")
	if _, err := os.Stat(shellProfile); os.IsNotExist(err) {
		shellProfile = ""
	}
	thisPath, err = os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}
}

func setProfileCommand() error {
	basename := path.Base(thisPath)

	content, err := os.ReadFile(shellProfile)
	if err != nil {
		return fmt.Errorf("Error reading shell profile: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	notifyCmdOk := -1
	for i, line := range lines {
		if notifyCmdOk < 0 && strings.Contains(line, basename) {
			if strings.Contains(line, basename+" notify") {
				log.Print("Notify command already set")
				notifyCmdOk = i
				break
			}
		}
	}

	command := fmt.Sprintf("[ -f %s ] && %s notify", thisPath, thisPath)

	if notifyCmdOk < 0 {
		lines = append(lines, command)
	} else {
		lines[notifyCmdOk] = command
	}
	content = []byte(strings.Join(lines, "\n"))
	err = os.WriteFile(shellProfile, content, 0644)
	if err == nil {
		log.Printf("Shell integration set: %s", thisPath)
	} else {
		log.Printf("Error setting shell integration: %v", err)
	}
	return err

}

func unsetProfileCommand() error {
	basename := path.Base(thisPath)

	content, err := os.ReadFile(shellProfile)
	if err != nil {
		return fmt.Errorf("Error reading shell profile: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	notifyCmdOk := -1
	for i, line := range lines {
		if notifyCmdOk < 0 && strings.Contains(line, basename) {
			if strings.Contains(line, basename+" notify") {
				log.Print("Notify command already set")
				notifyCmdOk = i
				break
			}
		}
	}

	if notifyCmdOk < 0 {
		log.Printf("Shell integration not set")
		return nil
	}
	lines = append(lines[:notifyCmdOk], lines[notifyCmdOk+1:]...)

	content = []byte(strings.Join(lines, "\n"))
	err = os.WriteFile(shellProfile, content, 0644)
	if err == nil {
		log.Printf("Shell integration set: %s", thisPath)
	} else {
		log.Printf("Error setting shell integration: %v", err)
	}
	return err
}
