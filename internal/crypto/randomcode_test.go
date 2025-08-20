package crypto

import (
	"strings"
	"testing"
	"unicode"
)

func TestGenerateEmailVerificationCode(t *testing.T) {
	code := GenerateEmailVerificationCode()
	if len(code) != 8 {
		t.Errorf("length of email verification code is invalid. got=%q, want=%q", len(code), 8)
	}

	for _, r := range code {
		if !unicode.IsDigit(r) {
			t.Errorf("code contains non-digit character: %q", r)
		}
	}
}

func TestGeneratePasswordResetCode(t *testing.T) {
	code := GeneratePasswordResetCode()
	if len(code) != 12 {
		t.Errorf("length of password reset code is invalid. got=%q, want=%q", len(code), 12)
	}

	for _, r := range code {
		if !unicode.IsDigit(r) {
			t.Errorf("code contains non-digit character: %q", r)
		}
	}
}

func TestGenerateRandomUint32(t *testing.T) {
	const max uint32 = 100

	for range 1000 {
		num := GenerateRandomUint32(max)
		if num >= max {
			t.Errorf("generated number is out of range. got=%d, want < %d", num, max)
		}
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length uint32
	}{
		{
			name:   "Length of 0",
			length: 0,
		},
		{
			name:   "Length of 5",
			length: 5,
		},
		{
			name:   "Length of 10",
			length: 10,
		},
		{
			name:   "Length of 20",
			length: 20,
		},
		{
			name:   "Length of 50",
			length: 50,
		},
		{
			name:   "Length of 100",
			length: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := GenerateRandomString(tt.length)
			if len(str) != int(tt.length) {
				t.Errorf("length of random string is incorrect. got=%q, want=%q", len(str), tt.length)
			}

			for _, r := range str {
				if !strings.ContainsRune(alphabet, r) {
					t.Errorf("code contains invalid character: %q", r)
				}
			}
		})
	}
}
