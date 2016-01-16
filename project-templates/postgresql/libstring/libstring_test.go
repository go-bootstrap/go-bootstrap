package libstring

import (
	"testing"
)

func TestRandString(t *testing.T) {
	generated := RandString(32)
	if len(generated) != 32 {
		t.Fatalf("Generated string has unexpected length. String: %v(%v)", generated, len(generated))
	}
}
