package domain

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}

// UpdateUserRequest defines the payload for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" example:"Alice Updated"`
	Email    string `json:"email,omitempty" example:"alice.new@example.com"`
	Password string `json:"password,omitempty" example:"123456"`
}

// LoginRequest defines the payload for user login
type LoginRequest struct {
	Email    string `json:"email" example:"test2@gmail.com"`
	Password string `json:"password" example:"123456"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !re.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
