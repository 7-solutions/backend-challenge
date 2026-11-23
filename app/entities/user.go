package entities

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UserFilter struct {
	ID    *string
	Name  *string
	Email *string
}
