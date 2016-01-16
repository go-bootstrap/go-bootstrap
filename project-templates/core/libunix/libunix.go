package libunix

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func CurrentUser() (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo $USER")
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
