package password

import "golang.org/x/crypto/bcrypt"

// CheckPasswordHash reports whether password matches hash, a bcrypt hash
// produced by HashPassword. It returns false for any comparison failure,
// including a malformed or empty hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
