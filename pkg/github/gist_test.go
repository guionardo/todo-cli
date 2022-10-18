package github

import (
	"os"
	"path"
	"testing"
)

func TestNewGitHubGistAPI(t *testing.T) {
	t.Run("Create, publish and consume gist", func(t *testing.T) {
		api := NewGitHubGistAPI("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")
		api.GistDescription = "ToDoCli Test"
		err := api.GetToDoConfigFileGist()
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		configContent := "TEST CONTENT"
		configFile := path.Join(t.TempDir(), "todo-cli.yaml")
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Errorf("Error writing config file: %s", err)
		}

		err = api.SetConfigFileGist(configFile)
		if err != nil {
			t.Errorf("Error setting config file gist: %s", err)
		}
		api.ConfigFileContent = nil
		if err = api.GetConfigFileContent(); err != nil {
			t.Errorf("Error getting config file content: %s", err)
		}
		if string(api.ConfigFileContent) != configContent {
			t.Errorf("Expected config content %s, got %s", configContent, string(api.ConfigFileContent))
		}
		api.DeleteGist()
	})
}
