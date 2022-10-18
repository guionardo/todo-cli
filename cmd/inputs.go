package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

func inputText(prompt string, defaultValue string) string {
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

func askYesNo(defaultValue bool, prompt string, args ...any) bool {
	var defaultText string
	if defaultValue {
		defaultText = "Y"
	} else {
		defaultText = "N"
	}
	text := inputText(fmt.Sprintf(prompt, args...), defaultText)
	return strings.ToUpper(text) == "Y"
}

func getToDoId(c *cli.Context) (int, error) {
	if c.NArg() == 0 {
		return 0, fmt.Errorf("Missing todo-id")
	}
	todoId := c.Args().Get(0)
	id, err := strconv.Atoi(todoId)
	if err != nil {
		return 0, fmt.Errorf("Invalid todo-id")
	}
	return id, nil
}
