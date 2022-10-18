package todo

import (
	"time"

	"github.com/guionardo/todo-cli/pkg/github"
	"gopkg.in/yaml.v3"
)

func (c *ToDoCollection) GistSync(config *github.GistConfig) (diffCount int, log []string, err error) {
	var api *github.GistAPI
	if api, err = github.NewGistAPI(config); err != nil {
		return
	}
	// Parse Gist content
	var gistColl ToDoCollection

	if len(api.GistContent) == 0 {
		gistColl = *NewTodoCollection()
	} else if gistColl, err = LoadCollectionFromData(api.GistContent); err != nil {
		gistColl = *NewTodoCollection()
	}
	upload := false
	log, diffCount, upload = c.Merge(&gistColl)
	config.LastSync = time.Now()
	if !upload {
		return
	}
	var content []byte
	if content, err = yaml.Marshal(c); err == nil {
		if err = api.Save(content); err == nil {
			config.GistId = api.Config.GistId
			config.RawURL = api.Config.RawURL
			config.GistDescription = api.Config.GistDescription
		}

	}

	return
}
