package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/guionardo/todo-cli/internal"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
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
func ActionSetup(c *cli.Context) error {
	fmt.Println("Starting todo setup")
	defaultCollectionName := fmt.Sprintf("%s's TODO", getUser())
	collectionName := inputText("Collection name", defaultCollectionName)

	fmt.Println("Collection name:", collectionName)

	fmt.Println("Create a new github token at https://github.com/settings/tokens/new with gist permission")
	fmt.Print("Github authentication token: ")
	auth, err := term.ReadPassword(int(syscall.Stdin))
	if len(auth) == 0 {
		err = fmt.Errorf("Authentication token is required")
	}
	if err != nil {
		return err
	}
	config := internal.Config{
		Authorization: string(auth),
		ToDoListName:  collectionName,
	}
	collection := internal.ToDoCollection{
		Config: config,
	}
	configFile, err := internal.CollectionFile()
	err = collection.Save(configFile)
	if err == nil {
		log.Printf("Collection saved to %s", configFile)
	} else {
		log.Printf("Error saving collection to %s: %v", configFile, err)
	}
	return err

}
