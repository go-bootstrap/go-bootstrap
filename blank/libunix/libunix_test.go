package libunix

import (
	"os"
	"runtime"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	var userEnv string
	if runtime.GOOS == "windows" {
		userEnv = os.Getenv("USERNAME")
	} else {
		userEnv = os.Getenv("USER")
	}
	username, err := CurrentUser()
	if userEnv != "" && err != nil {
		t.Fatalf("If $USER is not blank, error should not happen. Error: %v", err)
	}
	if userEnv != username {
		t.Errorf("Fetched the wrong username. $USER: %v, username: %v", userEnv, username)
	}
}
