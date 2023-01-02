package ctx

import (
	"context"
	"fmt"
	"time"

	"github.com/guionardo/go-gstools/gist"
)

type GistConfig struct {
	Authorization    string        `yaml:"authorization"`
	GistId           string        `yaml:"gist_id"`
	RawURL           string        `yaml:"raw_url"`
	AutoSync         bool          `yaml:"auto_sync"`
	GistDescription  string        `yaml:"gist_description"`
	LastSync         time.Time     `yaml:"last_sync"`
	AutoSyncInterval time.Duration `yaml:"auto_sync_interval"`
}

func (c *GistConfig) String() string {
	return fmt.Sprintf("GistId: %s\nAuthorization: %s\nURL: %s\nAutoSync: %t\n",
		c.GistId, c.Authorization, c.RawURL, c.AutoSync)
}

func GetDefaultGistConfig() GistConfig {
	return GistConfig{
		Authorization:    "",
		GistId:           "",
		RawURL:           "",
		AutoSync:         true,
		GistDescription:  "ToDo List",
		LastSync:         time.Time{},
		AutoSyncInterval: time.Hour * 6,
	}
}

func (c *GistConfig) SetToken(token string) error {
	if token == c.Authorization {
		return nil
	}
	gistCtx, err := gist.NewGitContext(token, context.Background())
	if err != nil {
		return err
	}

	c.Authorization = token

	remoteGist, err := gist.GetGistByDescription(gistCtx, c.GistDescription)
	if err == nil && remoteGist != nil {
		c.GistId = *remoteGist.ID
		c.GistDescription = *remoteGist.Description
		c.RawURL = *remoteGist.HTMLURL
	}

	return nil
}
