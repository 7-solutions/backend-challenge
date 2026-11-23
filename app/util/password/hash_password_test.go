package password_test

import (
	"strings"
	"testing"

	"github.com/chaihanij/backend-challenge/app/util/password"
)

func TestHashPassword_ReturnsHashDifferentFromInput(t *testing.T) {
	hash, err := password.HashPassword("s3cret-passw0rd")
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned an empty hash")
	}

	if hash == "s3cret-passw0rd" {
		t.Fatal("HashPassword() returned the plaintext password instead of a hash")
	}
}

func TestHashPassword_SamePasswordProducesDifferentHashes(t *testing.T) {
	const plain = "s3cret-passw0rd"

	first, err := password.HashPassword(plain)
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	second, err := password.HashPassword(plain)
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	if first == second {
		t.Fatal("HashPassword() produced identical hashes for two calls with the same password; each call should use a fresh salt")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := password.HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error for empty password: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned an empty hash for an empty password")
	}
}

func TestHashPassword_PasswordExceedsBcryptLimit(t *testing.T) {
	// bcrypt rejects passwords longer than 72 bytes.
	tooLong := strings.Repeat("a", 73)

	_, err := password.HashPassword(tooLong)
	if err == nil {
		t.Fatal("HashPassword() expected an error for a password longer than 72 bytes, got nil")
	}
}
