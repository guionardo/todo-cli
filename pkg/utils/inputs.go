package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/guionardo/todo-cli/pkg/logger"
)

func InputText(prompt string, defaultValue string) string {
	fmt.Printf("%s [%s]:", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		logger.Fatalf("Input error: %v", err)
	}
	text = strings.Replace(text, "\n", "", -1)

	if len(text) == 0 {
		text = defaultValue
	}
	return text
}

func AskYesNo(defaultValue bool, prompt string, args ...any) bool {
	var defaultText string
	if defaultValue {
		defaultText = "Y"
	} else {
		defaultText = "N"
	}
	text := InputText(fmt.Sprintf(prompt, args...), defaultText)
	return strings.ToUpper(text) == "Y"
}
