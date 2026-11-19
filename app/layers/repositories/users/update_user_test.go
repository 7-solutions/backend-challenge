package users_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	mockmongo "github.com/chaihanij/backend-challenge/app/test_helpers/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestRepo_UpdateUser(t *testing.T) {
	testhelpers.Init()

	t.Run("Success", func(t *testing.T) {
		id := bson.NewObjectID().Hex()
		called := false
		fakeCol := &mockmongo.FakeCollection{
			UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
				called = true
				return &mongo.UpdateResult{Acknowledged: true, ModifiedCount: 1}, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			ID:    id,
			Email: "updated@example.com",
			Name:  "Updated Name",
		}

		if err := repo.UpdateUser(context.Background(), input); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Fatal("expected UpdateOne to be called")
		}
	})

	t.Run("Failure: No Match", func(t *testing.T) {
		id := bson.NewObjectID().Hex()
		fakeCol := &mockmongo.FakeCollection{
			UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
				return &mongo.UpdateResult{Acknowledged: true, ModifiedCount: 0}, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			ID:   id,
			Name: "Updated Name",
		}

		err := repo.UpdateUser(context.Background(), input)
		if err == nil {
			t.Fatal("expected not found error, got nil")
		}

		var notFound *customerrors.NotFoundError
		if !errors.As(err, &notFound) {
			t.Fatalf("expected NotFoundError, got %T: %v", err, err)
		}
		if notFound.Resource != "User" {
			t.Fatalf("expected resource %q, got %q", "User", notFound.Resource)
		}
		if notFound.ID != id {
			t.Fatalf("expected id %q, got %q", id, notFound.ID)
		}
	})

	t.Run("Failure: Nil Input", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
				t.Fatal("UpdateOne should NOT be called when input is nil")
				return nil, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		err := repo.UpdateUser(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error for nil input, got nil")
		}
	})

	t.Run("Failure: Model Error (Password Hash)", func(t *testing.T) {
		id := bson.NewObjectID().Hex()
		fakeCol := &mockmongo.FakeCollection{
			UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
				t.Fatal("UpdateOne should NOT be called when model build fails")
				return nil, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			ID:       id,
			Password: strings.Repeat("A", 100), // bcrypt error
		}

		err := repo.UpdateUser(context.Background(), input)
		if err == nil {
			t.Fatal("expected model error, got nil")
		}
	})

	t.Run("Failure: UpdateOne Error", func(t *testing.T) {
		id := bson.NewObjectID().Hex()
		fakeCol := &mockmongo.FakeCollection{
			UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
				return nil, errors.New("update failed")
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			ID:   id,
			Name: "Updated Name",
		}

		err := repo.UpdateUser(context.Background(), input)
		if err == nil {
			t.Fatal("expected update error, got nil")
		}
	})

	t.Run("Failure: Filter Error", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
				t.Fatal("UpdateOne should NOT be called when filter build fails")
				return nil, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		input := &entities.User{
			ID:   "not-an-object-id",
			Name: "Updated Name",
		}

		if err := repo.UpdateUser(context.Background(), input); err == nil {
			t.Fatal("expected filter error, got nil")
		}
	})
}
