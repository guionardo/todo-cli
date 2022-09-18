package internal

import (
	"path"
	"reflect"
	"testing"
)

func TestConfig_SaveAndLoad(t *testing.T) {

	t.Run("Save and load", func(t *testing.T) {
		c := Config{
			GistId:        "gist_id",
			Authorization: "github_token",
		}
		configFile := path.Join(t.TempDir(), "config.json")
		if err := c.Save(configFile); err != nil {
			t.Errorf("Config.Save() error = %v", err)
		}
		
		var c2 Config
		if err := c2.Load(configFile); err != nil {
			t.Errorf("Config.Load() error = %v", err)
		}
		if !reflect.DeepEqual(c, c2) {
			t.Errorf("Config.Load() = %v, want %v", c2, c)
		}
	})

}
