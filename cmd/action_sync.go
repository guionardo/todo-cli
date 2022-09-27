package cmd

import (
	"fmt"
	"log"
	"syscall"

	"github.com/guionardo/todo-cli/internal"
	"github.com/guionardo/todo-cli/pkg/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func ActionSync(c *cli.Context) error {
	fmt.Println("Starting todo setup")
	context := internal.GetRunningContext(c)

	configFile := c.String("config")
	if len(configFile) == 0 {
		configFile = internal.DefaultCollectionFilePath
		context.DebugLog("Using default config file: %s", configFile)
	} else {
		context.DebugLog("Using config file: %s", configFile)
	}
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
	api := github.NewGitHubGistAPI(string(auth), context.DebugMode)

	err = api.GetToDoConfigFileGist()
	gistItems := make([]*internal.ToDoItem, 0)
	if err != nil {
		if err.Error() == "Invalid token" {
			fmt.Println("Invalid token")
			return err
		}
	} else {
		if api.GetConfigFileContent() == nil {
			gistConfig, err := internal.ParseCollectionData(api.ConfigFileContent)
			if err == nil {
				gistItems = gistConfig.Items
			}
		}
	}
	config := internal.Config{
		Authorization: string(auth),
		ToDoListName:  collectionName,
	}
	collection := internal.ToDoCollection{
		Config: config,
	}

	existentConfig, err := internal.ParseCollectionFile(configFile)
	if err == nil {
		if len(existentConfig.Items) > 0 {
			collection.Items = existentConfig.Items
		}
	}
	collection.Items = internal.MergeToDoItems(collection.Items, gistItems)
	err = collection.Save(configFile)
	if err == nil {
		log.Printf("Collection saved to %s", configFile)
		err = api.SetConfigFileGist(configFile)
		if err == nil {
			if collection.Config.GistId != api.GistId {
				collection.Config.GistId = api.GistId
				collection.Save(configFile)
			}
			log.Printf("Collection saved to gist %s", api.GistId)
		} else {
			log.Printf("Error saving collection to gist %v", err)
		}
	} else {
		log.Printf("Error saving collection to %s: %v", configFile, err)
	}
	return err

}

