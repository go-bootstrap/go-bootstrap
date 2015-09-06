package libunix

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

func CurrentUser() (string, error) {
	var stdout bytes.Buffer
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("sh", "-c", "eval echo $USERNAME")
	} else {
		cmd = exec.Command("sh", "-c", "eval echo $USER")
	}
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading current user.")
	}

	return result, nil
}
