// Package password provides bcrypt-based helpers for hashing and verifying
// user passwords.
//
// HashPassword produces a salted bcrypt hash suitable for storage:
//
//	hash, err := password.HashPassword("s3cret")
//
// CheckPasswordHash later verifies a plaintext password against that stored
// hash:
//
//	ok := password.CheckPasswordHash("s3cret", hash) // true
package password
