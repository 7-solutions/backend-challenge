package models_test

import (
	"testing"
	"time"

	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUser_ToEntity(t *testing.T) {
	now := time.Now().UTC().Round(time.Second)
	id := bson.NewObjectID()

	input := models.User{
		ID:        id,
		Name:      "Tester",
		Email:     "tester@example.com",
		Password:  "hashed-password",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	entity, err := input.ToEntity()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entity == nil {
		t.Fatal("expected non-nil entity")
	}

	if entity.ID != id.Hex() {
		t.Fatalf("expected ID %q, got %q", id.Hex(), entity.ID)
	}
	if entity.Name != input.Name {
		t.Fatalf("expected Name %q, got %q", input.Name, entity.Name)
	}
	if entity.Email != input.Email {
		t.Fatalf("expected Email %q, got %q", input.Email, entity.Email)
	}
	if entity.Password != input.Password {
		t.Fatalf("expected Password %q, got %q", input.Password, entity.Password)
	}
	if entity.CreatedAt == nil || !entity.CreatedAt.Equal(now) {
		t.Fatalf("expected CreatedAt %v, got %#v", now, entity.CreatedAt)
	}
	if entity.UpdatedAt == nil || !entity.UpdatedAt.Equal(now) {
		t.Fatalf("expected UpdatedAt %v, got %#v", now, entity.UpdatedAt)
	}
}

func TestUsers_ToEntity(t *testing.T) {
	t.Run("maps all users", func(t *testing.T) {
		now := time.Now().UTC().Round(time.Second)
		id1 := bson.NewObjectID()
		id2 := bson.NewObjectID()

		input := models.Users{
			{
				ID:        id1,
				Name:      "A",
				Email:     "a@example.com",
				Password:  "hash-a",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			{
				ID:        id2,
				Name:      "B",
				Email:     "b@example.com",
				Password:  "hash-b",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		}

		entities, err := input.ToEntity()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(entities) != 2 {
			t.Fatalf("expected 2 entities, got %d", len(entities))
		}

		if entities[0].ID != id1.Hex() {
			t.Fatalf("expected first ID %q, got %q", id1.Hex(), entities[0].ID)
		}
		if entities[1].ID != id2.Hex() {
			t.Fatalf("expected second ID %q, got %q", id2.Hex(), entities[1].ID)
		}
	})

	t.Run("empty collection returns empty slice", func(t *testing.T) {
		input := models.Users{}

		entities, err := input.ToEntity()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(entities) != 0 {
			t.Fatalf("expected empty result, got %d items", len(entities))
		}
	})
}
