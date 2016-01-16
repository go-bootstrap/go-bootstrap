package libhttp

import (
	"testing"
)

func TestParseBasicAuth(t *testing.T) {
	auth := "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ=="
	username, password, ok := ParseBasicAuth(auth)

	if username != "Aladdin" {
		t.Errorf("Username is not as expected. Received: %v", username)
	}
	if password != "open sesame" {
		t.Errorf("Password is not as expected. Received: %v", password)
	}
	if !ok {
		t.Error("Parsing basic auth should work.")
	}
}
