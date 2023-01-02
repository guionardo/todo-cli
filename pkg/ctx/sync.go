package ctx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/guionardo/go-gstools/gist"
	"github.com/guionardo/todo-cli/pkg/backup"
	"github.com/guionardo/todo-cli/pkg/todo"
)

func (config *Context) GistSync() (log []string, err error) {
	log = make([]string, 0, 10)
	if config.LocalConfig.Gist.Authorization == "" {
		err = errors.New("gist authorization not set")
		return
	}
	var (
		ctxGist    context.Context
		localGist  *github.Gist
		remoteGist *github.Gist
	)
	if ctxGist, err = gist.NewGitContext(
		config.LocalConfig.Gist.Authorization,
		context.Background()); err != nil {
		return
	}

	localCollection, err := todo.LoadCollection(config.LocalCollectionFile)

	if err != nil {
		// error on reading or parsing collection file
		newFilename, err := backup.MoveFileToBackup(config.LocalCollectionFile)
		log = append(log, fmt.Sprintf("Error on reading or parsing collection file. Moving to %s - %v", newFilename, err))
		localCollection = *todo.NewTodoCollection()
	}

	remoteCollection := *todo.NewTodoCollection()
	if config.LocalConfig.Gist.GistId != "" {
		log = append(log, "Loading remote collection... Gist ID: "+config.LocalConfig.Gist.GistId)
		remoteGist, err = gist.GetGistById(ctxGist, config.LocalConfig.Gist.GistId)
		if err == nil {
			if collectionFile, ok := remoteGist.Files[github.GistFilename(path.Base(config.LocalCollectionFile))]; ok {
				var collectionData string = *collectionFile.Content
				remoteCollection, err = todo.LoadCollectionFromData([]byte(collectionData))
			}
		}
		if err != nil {
			log = append(log, fmt.Sprintf("Error on reading or parsing remote collection file. %v", err))
			remoteCollection = *todo.NewTodoCollection()
		}
	}

	mergeLog, diffCount, _ := localCollection.Merge(&remoteCollection)

	if diffCount == 0 {
		log = append(log, "No changes detected")
		return
	}
	log = append(log, mergeLog...)
	tmpDir, err := os.MkdirTemp(config.DataFolder, "todo-sync-*")
	if err != nil {
		log = append(log, fmt.Sprintf("Error on creating temporary folder. %v", err))
		return
	}
	defer os.RemoveAll(tmpDir)
	tmpCol := path.Join(tmpDir, path.Base(config.LocalCollectionFile))
	tmpCfg := path.Join(tmpDir, path.Base(config.LocalConfigFile))

	if err = localCollection.Save(tmpCol); err != nil {
		log = append(log, fmt.Sprintf("Error on saving temporary collection file. %v", err))
		return
	}

	newConfig := *config.LocalConfig
	newConfig.Gist.LastSync = time.Now()

	if err = newConfig.Save(tmpCfg); err != nil {
		log = append(log, fmt.Sprintf("Error on saving temporary config file. %v", err))
		return
	}

	localGist, err = gist.CreateGistFromLocalFiles(
		config.LocalConfig.ToDoListName,
		tmpCol)

	if err != nil {
		log = append(log, fmt.Sprintf("Error on creating gist from local files. %v", err))
		return
	}
	localGist.ID = github.String(config.LocalConfig.Gist.GistId)
	localGist.Public = github.Bool(false)
	remoteGist, action, err := gist.SyncGistFiles(ctxGist, localGist, config.DataFolder)
	if err != nil {
		log = append(log, fmt.Sprintf("Error on syncing gist files. %v", err))
		return
	}
	config.LocalConfig.Gist.LastSync = *remoteGist.UpdatedAt
	config.LocalConfig.Gist.GistId = *remoteGist.ID
	config.LocalConfig.Gist.GistDescription = *remoteGist.Description
	config.LocalConfig.Save(config.LocalConfigFile)
	switch action {
	case gist.NoAction:
		log = append(log, "No changes detected")
	case gist.Upload:
		log = append(log, fmt.Sprintf("Uploading files to GIST #%s - %s", *remoteGist.ID, *remoteGist.Description))
	case gist.Download:
		log = append(log, fmt.Sprintf("Downloading files from GIST #%s - %s", *remoteGist.ID, *remoteGist.Description))
		localCollection, err = todo.LoadCollection(tmpCol)
		if err != nil {
			log = append(log, fmt.Sprintf("Error on reading or parsing collection file. %v", err))
		} else {
			err = localCollection.Save(config.LocalCollectionFile)
			if err != nil {
				log = append(log, fmt.Sprintf("Error on saving collection file. %v", err))
			} else {
				log = append(log, "Collection file updated")
			}
		}
	}
	return
}
