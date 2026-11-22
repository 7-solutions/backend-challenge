package models_test

import (
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users/models"
)

func TestUserResponseFromEntity(t *testing.T) {
	user := &entities.User{
		ID:       "507f1f77bcf86cd799439011",
		Email:    "user@example.com",
		Name:     "Example User",
		Password: "must-not-be-exposed",
	}

	got := models.UserResponseFromEntity(user)
	want := &models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
	if *got != *want {
		t.Errorf("UserResponseFromEntity() = %+v, want %+v", got, want)
	}
}

func TestUsersResponseFromEntity(t *testing.T) {
	users := []entities.User{
		{ID: "id-1", Email: "first@example.com", Name: "First User", Password: "secret-1"},
		{ID: "id-2", Email: "second@example.com", Name: "Second User", Password: "secret-2"},
	}

	got := models.UsersResponseFromEntity(users)
	want := models.UsersResponse{
		{ID: "id-1", Email: "first@example.com", Name: "First User"},
		{ID: "id-2", Email: "second@example.com", Name: "Second User"},
	}
	if len(got) != len(want) {
		t.Fatalf("UsersResponseFromEntity() returned %d users, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("UsersResponseFromEntity()[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestUsersResponseFromEmptyEntitySlice(t *testing.T) {
	got := models.UsersResponseFromEntity([]entities.User{})
	if len(got) != 0 {
		t.Errorf("UsersResponseFromEntity(empty) returned %d users, want 0", len(got))
	}
}
