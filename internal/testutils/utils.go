package testutils

import (
	"regexp"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
)

var rxCSRF = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+?)">`)

func ExtractCSRFToken(t *testing.T, body string) string {
	matches := rxCSRF.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatalf("no csrf token found in body")
	}

	return matches[1]
}

var hasher crypto.PasswordHasher

func MustHashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := hasher.GenerateFromPassword(password, crypto.DefaultPasswordHashParams)
	if err != nil {
		t.Fatalf("error generating hash for password %q: %v", password, err)
	}
	return hash
}
