package libenv

import (
	"testing"
)

func TestEnvWithDefault(t *testing.T) {
	if EnvWithDefault("BROTATO", "I expect this value") != "I expect this value" {
		t.Errorf("Default environment variable was not set correctly")
	}
	if EnvWithDefault("PATH", "I do not expect this value") == "I do not expect this value" {
		t.Errorf("Existing environment variable should not be overriden")
	}
}
