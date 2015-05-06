// Package libenv provides environment related library functions.
package libenv

import (
	"os"
)

func EnvWithDefault(name string, defaultVal string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultVal
	}
	return value
}
