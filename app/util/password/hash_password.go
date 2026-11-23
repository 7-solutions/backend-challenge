package password

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes password using bcrypt with the library's default cost
// factor. The returned hash embeds both the salt and the cost, so it can be
// verified later with CheckPasswordHash without any extra state.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
