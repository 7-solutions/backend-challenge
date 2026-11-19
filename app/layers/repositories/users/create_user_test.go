package users_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
	mockmongo "github.com/chaihanij/backend-challenge/app/test_helpers/mongo"
)

func TestRepo_CreateUser(t *testing.T) {

	testhelpers.Init()

	t.Run("Success", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			InsertOneFn: func(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
				user := doc.(*models.User)
				return &mongo.InsertOneResult{InsertedID: user.ID}, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			Email:    "test@example.com",
			Name:     "Tester",
			Password: "P@ssword",
		}

		result, err := repo.CreateUser(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID == "" {
			t.Fatalf("expected id to be set, got empty")
		}
		if result.Email != input.Email {
			t.Fatalf("expected email %s, got %s", input.Email, result.Email)
		}
		if result.Name != input.Name {
			t.Fatalf("expected name %s, got %s", input.Name, result.Name)
		}
		if result.Password == "" {
			t.Fatalf("expected password to be hashed, got empty")
		}
		if result.CreatedAt == nil {
			t.Fatalf("expected created_at to be set")
		}
		if result.UpdatedAt == nil {
			t.Fatalf("expected updated_at to be set")
		}
	})

	t.Run("Failure: Model Error", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			InsertOneFn: func(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
				t.Fatalf("InsertOne should NOT be called when model error occurs")
				return nil, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			Email:    "test@example.com",
			Name:     "Tester",
			Password: strings.Repeat("A", 100), // bcrypt error
		}

		_, err := repo.CreateUser(context.Background(), input)
		if err == nil {
			t.Fatalf("expected model error, got nil")
		}
	})

	t.Run("Failure: InsertOne Error", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			InsertOneFn: func(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
				return nil, errors.New("insert failed")
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			Email:    "test@example.com",
			Name:     "Tester",
			Password: "P@ssword",
		}

		_, err := repo.CreateUser(context.Background(), input)
		if err == nil {
			t.Fatalf("expected insert error, got nil")
		}
	})

	t.Run("Failure: InsertedID is nil", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			InsertOneFn: func(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
				return &mongo.InsertOneResult{InsertedID: nil}, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			Email:    "test@example.com",
			Name:     "Tester",
			Password: "P@ssword",
		}

		_, err := repo.CreateUser(context.Background(), input)
		if err == nil {
			t.Fatalf("expected nil InsertedID error, got nil")
		}
	})
}
