package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/guionardo/todo-cli/pkg/github"
)

func (collection *ToDoCollection) GISTSync(debugMode bool) error {
	if len(collection.Config.Authorization) == 0 {
		return errors.New("Authorization is empty")
	}
	api := github.NewGitHubGistAPI(collection.Config.Authorization, debugMode)
	err := api.GetToDoConfigFileGist()
	gistItems := make([]*ToDoItem, 0)
	if err != nil {
		if err.Error() == "Invalid token" {
			return err
		}
	} else {
		if api.GetConfigFileContent() == nil {
			gistConfig, err := ParseCollectionData(api.ConfigFileContent)
			if err == nil {
				gistItems = gistConfig.Items
			}
		}
	}
	collection.Items = MergeToDoItems(collection.Items, gistItems)
	old_last_sync := collection.LastSync
	collection.LastSync = time.Now()
	err = collection.Save(collection.Config.ConfigFileName)
	if err != nil {
		return errors.New(fmt.Sprintf("Error saving collection to %s: %v", collection.Config.ConfigFileName, err))
	}

	err = api.SetConfigFileGist(collection.Config.ConfigFileName)
	if err != nil {
		collection.LastSync = old_last_sync
		collection.Save(collection.Config.ConfigFileName)
		return errors.New(fmt.Sprintf("Error saving collection to gist: %v", err))
	}
	collection.LastSync = time.Now()
	collection.Config.GistId = api.GistId
	err = collection.Save(collection.Config.ConfigFileName)
	if err != nil {
		return errors.New(fmt.Sprintf("Error saving collection to %s: %v", collection.Config.ConfigFileName, err))
	}

	return nil
}
