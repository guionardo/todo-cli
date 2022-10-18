package github

import (
	"os"
	"testing"

	"github.com/joho/godotenv"

	"github.com/guionardo/todo-cli/pkg/logger"
)

func TestNewGistAPI(t *testing.T) {
	t.Run("Default", func(t *testing.T) {	
		t.Skip()	
		logger.SetLogger(true)
		godotenv.Load()
		auth:=os.Getenv("GITHUB_TOKEN")
		if auth == "" {
			t.Skip("Skipping test because GITHUB_TOKEN is not set")
		}
		config := &GistConfig{
			Authorization:   auth,
			GistDescription: "TEST GIST",
		}
		api, err := NewGistAPI(config)
		if err != nil {
			t.Errorf("NewGistAPI() error = %v", err)
			return
		}
		if len(api.GistContent) > 0 {
			api.Delete()
		}

		if err := api.Save([]byte("TEST GIST CONTENT")); err != nil {
			t.Errorf("Save() error = %v", err)
			return
		}

		newApi, err := NewGistAPI(api.Config)
		if err != nil {
			t.Errorf("NewGistAPI() error = %v", err)
			return
		}
		if len(newApi.GistContent) != len(api.GistContent) {
			t.Errorf("Expected same size content")
			return
		}

		if err := newApi.Delete(); err != nil {
			t.Errorf("Delete() error = %v", err)
			return
		}
	})

}
