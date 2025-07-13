package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testGolang/backend-challenge/application/usecases"
	"github.com/testGolang/backend-challenge/domain"
	"golang.org/x/crypto/bcrypt"
)

// mockRepo implements UserRepository for testing.
type mockRepo struct {
	users map[string]*domain.User
}

func (m *mockRepo) Create(u *domain.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockRepo) GetByID(id string) (*domain.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockRepo) GetByEmail(email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockRepo) GetAll() ([]*domain.User, error) {
	list := make([]*domain.User, 0, len(m.users))
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, nil
}

func (m *mockRepo) Update(u *domain.User) error {
	if _, ok := m.users[u.ID]; !ok {
		return errors.New("user not found")
	}
	m.users[u.ID] = u
	return nil
}

func (m *mockRepo) Delete(id string) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *mockRepo) Count() (int64, error) {
	return int64(len(m.users)), nil
}

func TestRegister_Success(t *testing.T) {
	repo := &mockRepo{users: make(map[string]*domain.User)}
	uc := usecases.NewUserUseCase(repo)

	user, err := uc.Register("Alice", "alice@example.com", "secret")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.NotEmpty(t, user.Password)

	// Ensure the password was hashed correctly
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret"))
	assert.NoError(t, err)

	t.Logf("Registered user: %+v", user)
}

func TestAuthenticate_Success(t *testing.T) {
	repo := &mockRepo{users: make(map[string]*domain.User)}
	uc := usecases.NewUserUseCase(repo)

	// Register first then authenticate
	_, _ = uc.Register("Bob", "bob@example.com", "password123")
	token, err := uc.Authenticate("bob@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	t.Logf("Generated JWT token: %s", token)
}

func TestCRUD_ListGetUpdateDelete(t *testing.T) {
	repo := &mockRepo{users: make(map[string]*domain.User)}
	uc := usecases.NewUserUseCase(repo)

	// Create two users
	u1, _ := uc.Register("User1", "u1@example.com", "pwd1")
	u2, _ := uc.Register("User2", "u2@example.com", "pwd2")

	// List
	all, err := uc.List()
	assert.NoError(t, err)
	assert.Len(t, all, 2)
	for i, u := range all {
		t.Logf("List user %d: %+v", i, *u)
	}

	// GetByID
	got, err := uc.Get(u1.ID)
	assert.NoError(t, err)
	assert.Equal(t, u1.ID, got.ID)
	t.Logf("Fetched user: %+v", got)

	// Update
	updated, err := uc.Update(u1.ID, "User1New", "u1new@example.com", "test")
	assert.NoError(t, err)
	assert.Equal(t, "User1New", updated.Name)
	assert.Equal(t, "u1new@example.com", updated.Email)
	t.Logf("Updated user: %+v", updated)

	// Delete
	err = uc.Delete(u2.ID)
	assert.NoError(t, err)

	// Verify deletion
	all, _ = uc.List()
	assert.Len(t, all, 1)
	for i, u := range all {
		t.Logf("Remaining user %d: %+v", i, *u)
	}
}
