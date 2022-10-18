package github

import (
	"fmt"
	"time"
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

func (c GistConfig) String() string {
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
	oldAuth := c.Authorization
	c.Authorization = token

	api, err := NewGistAPI(c)
	if err != nil {
		c.Authorization = oldAuth
		return err
	}

	c.GistDescription = api.Config.GistDescription
	c.GistId = api.Config.GistId
	c.RawURL = api.Config.RawURL
	return nil
}
