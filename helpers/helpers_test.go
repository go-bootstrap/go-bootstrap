package helpers

import (
	"os"
	"strings"
	"testing"
)

func TestGoPaths(t *testing.T) {
	os.Setenv("GOPATH", "/foo:/bar")
	gopaths := GoPaths()
	if len(gopaths) != 2 || gopaths[0] != "/foo" || gopaths[1] != "/bar" {
		t.Errorf("Incorrectly split GOPATH. Output: %v", gopaths)
	}
}

func TestIsValidGoPath(t *testing.T) {
	os.Setenv("GOPATH", "/foo:/bar")
	if !IsValidGoPath("/foo") {
		t.Errorf("Valid GOPATH not recognized: /foo")
	}
	if IsValidGoPath("/baz") {
		t.Errorf("Invalid GOPATH recognized as valid: /baz")
	}
}

func TestBashEscape(t *testing.T) {
	os.Setenv("PGPASSWORD", "boo")
	dbName := DefaultPGDSN("default")
	escapedDbName := BashEscape(dbName)
	if !strings.Contains(escapedDbName, `\&`) {
		t.Errorf("Incorrectly escape dbName. Output: %v", escapedDbName)
	}
}

func TestRandString(t *testing.T) {
	str := RandString(12)
	if len(str) != 12 {
		t.Errorf("Generated string with incorrect length. Output: %v", str)
	}
}

func TestDefaultPGDSN(t *testing.T) {
	os.Setenv("PGPASSWORD", "boo")
	dbName := DefaultPGDSN("default")
	if !strings.HasPrefix(dbName, `postgres://`) {
		t.Errorf("Incorrectly generate DSN. Output: %v", dbName)
	}
	if !strings.Contains(dbName, "password=boo") {
		t.Errorf("Incorrectly generate DSN. Output: %v", dbName)
	}
}

func TestTrim(t *testing.T) {
	in := "/github.com/bro/tato/"
	out := strings.Trim(in, "/")
	if out != "github.com/bro/tato" {
		t.Errorf("Incorrectly trim string. Output: %v", out)
	}
}
