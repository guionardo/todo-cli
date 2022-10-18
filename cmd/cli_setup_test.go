package cmd

import (
	"testing"
)

func TestApp(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		app := App()
		if app.Name != "todo-cli" {
			t.Errorf("App.Name = %v, want %v", app.Name, "todo-cli")
		}

	})

}
