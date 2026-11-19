package password_test

import (
	"testing"

	"github.com/chaihanij/backend-challenge/app/util/password"
)

func TestCheckPasswordHash_CorrectPassword(t *testing.T) {
	const plain = "s3cret-passw0rd"

	hash, err := password.HashPassword(plain)
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	if !password.CheckPasswordHash(plain, hash) {
		t.Fatal("CheckPasswordHash() = false, want true for the password that produced the hash")
	}
}

func TestCheckPasswordHash_IncorrectPassword(t *testing.T) {
	hash, err := password.HashPassword("s3cret-passw0rd")
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	if password.CheckPasswordHash("wrong-password", hash) {
		t.Fatal("CheckPasswordHash() = true, want false for a mismatched password")
	}
}

func TestCheckPasswordHash_MalformedHash(t *testing.T) {
	if password.CheckPasswordHash("s3cret-passw0rd", "not-a-bcrypt-hash") {
		t.Fatal("CheckPasswordHash() = true, want false for a malformed hash")
	}
}

func TestCheckPasswordHash_EmptyHash(t *testing.T) {
	if password.CheckPasswordHash("s3cret-passw0rd", "") {
		t.Fatal("CheckPasswordHash() = true, want false for an empty hash")
	}
}

func TestCheckPasswordHash_EmptyPasswordAgainstRealHash(t *testing.T) {
	hash, err := password.HashPassword("s3cret-passw0rd")
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	if password.CheckPasswordHash("", hash) {
		t.Fatal("CheckPasswordHash() = true, want false when checking an empty password against a real hash")
	}
}
