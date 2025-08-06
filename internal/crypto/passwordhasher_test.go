package crypto_test

import (
	"strings"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
)

func TestPasswordHashing(t *testing.T) {
	tests := []struct {
		name            string
		password        string
		comparePassword string
		shouldMatch     bool
	}{
		{
			name:            "Empty password",
			password:        "",
			comparePassword: "",
			shouldMatch:     true,
		},
		{
			name:            "Password with only alphabetic characters",
			password:        "Password",
			comparePassword: "Password",
			shouldMatch:     true,
		},
		{
			name:            "Password with numbers",
			password:        "Password123",
			comparePassword: "Password123",
			shouldMatch:     true,
		},
		{
			name:            "Password with special characters",
			password:        "%P&a$$w0rd123*",
			comparePassword: "%P&a$$w0rd123*",
			shouldMatch:     true,
		},
		{
			name:            "Password with space",
			password:        "Pass word",
			comparePassword: "Pass word",
			shouldMatch:     true,
		},
		{
			name:            "Wrong password",
			password:        "Pa$$w0rd123",
			comparePassword: "Pa$$w0rd1234",
			shouldMatch:     false,
		},
		{
			name:            "Long password",
			password:        strings.Repeat("x", 1000),
			comparePassword: strings.Repeat("x", 1000),
			shouldMatch:     true,
		},
		{
			name:            "Unicode password",
			password:        "pāsswørd密码🔒",
			comparePassword: "pāsswørd密码🔒",
			shouldMatch:     true,
		},
	}

	hasher := crypto.PasswordHasher{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hash := testutils.MustHashPassword(t, hasher, tt.password)
			ok, err := hasher.ComparePasswordAndHash(tt.comparePassword, hash)
			if err != nil {
				t.Fatalf(
					"error while comparing password with hash. password=%s, comparePassword=%s, hash=%s, err=%v",
					tt.password, tt.comparePassword, hash, err,
				)
			}

			if tt.shouldMatch && !ok {
				t.Errorf(
					"comparing the generated hash with the compare password results in false. password=%s, comparePassword=%s",
					tt.password, tt.comparePassword,
				)
			}

			if !tt.shouldMatch && ok {
				t.Errorf(
					"comparing the generated hash with the compare password results in true. password=%s, comparePassword=%s",
					tt.password, tt.comparePassword,
				)
			}
		})
	}

	t.Run("Same password different hash", func(t *testing.T) {
		password := "Password123"
		hash1 := testutils.MustHashPassword(t, hasher, password)
		hash2 := testutils.MustHashPassword(t, hasher, password)
		if hash1 == hash2 {
			t.Logf("warning: hashes for same password are identical.")
		}
	})
}
