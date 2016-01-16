package libunix

import (
	"os"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	userEnv := os.Getenv("USER")
	username, err := CurrentUser()
	if userEnv != "" && err != nil {
		t.Fatalf("If $USER is not blank, error should not happen. Error: %v", err)
	}
	if userEnv != username {
		t.Errorf("Fetched the wrong username. $USER: %v, username: %v", userEnv, username)
	}
}
