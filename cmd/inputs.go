package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func inputText(prompt string, defaultValue string) string {
	fmt.Printf("%s [%s]:", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Input error: %v", err)
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
