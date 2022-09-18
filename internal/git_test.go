package internal

import (
	"testing"
)

func TestGetCurrentGitUser(t *testing.T) {

	t.Run("Current git", func(t *testing.T) {
		gotGitUser, err := GetCurrentGitUser()
		if err != nil {
			t.Errorf("GetCurrentGitUser() error = %v", err)
			return
		}
		if len(gotGitUser.Name) == 0 {
			t.Errorf("GetCurrentGitUser() = %v", gotGitUser)
		}

	})

}
