package ports

import "github.com/testGolang/backend-challenge/domain"

// UserRepository defines persistence operations.
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	Count() (int64, error)
	GetByEmail(email string) (*domain.User, error)
}
