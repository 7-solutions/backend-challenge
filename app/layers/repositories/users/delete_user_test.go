package users_test

import (
	"context"
	"errors"
	"testing"

	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	mockmongo "github.com/chaihanij/backend-challenge/app/test_helpers/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestRepo_DeleteUser(t *testing.T) {
	testhelpers.Init()

	t.Run("Success", func(t *testing.T) {
		id := bson.NewObjectID().Hex()
		called := false
		fakeCol := &mockmongo.FakeCollection{
			DeleteOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
				called = true
				return &mongo.DeleteResult{Acknowledged: true, DeletedCount: 1}, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		if err := repo.DeleteUser(context.Background(), id); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Fatal("expected DeleteOne to be called")
		}
	})

	t.Run("Failure: No Match", func(t *testing.T) {
		id := bson.NewObjectID().Hex()
		fakeCol := &mockmongo.FakeCollection{
			DeleteOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
				return &mongo.DeleteResult{Acknowledged: true, DeletedCount: 0}, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		err := repo.DeleteUser(context.Background(), id)
		if err == nil {
			t.Fatal("expected not found error, got nil")
		}

		var notFound *customerrors.NotFoundError
		if !errors.As(err, &notFound) {
			t.Fatalf("expected NotFoundError, got %T: %v", err, err)
		}
		if notFound.Resource != "user" {
			t.Fatalf("expected resource %q, got %q", "user", notFound.Resource)
		}
		if notFound.ID != id {
			t.Fatalf("expected id %q, got %q", id, notFound.ID)
		}
	})

	t.Run("Failure: Filter Error", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			DeleteOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
				t.Fatal("DeleteOne should NOT be called when filter build fails")
				return nil, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		if err := repo.DeleteUser(context.Background(), "not-an-object-id"); err == nil {
			t.Fatal("expected filter error, got nil")
		}
	})

	t.Run("Failure: DeleteOne Error", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			DeleteOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
				return nil, errors.New("delete failed")
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		if err := repo.DeleteUser(context.Background(), bson.NewObjectID().Hex()); err == nil {
			t.Fatal("expected delete error, got nil")
		}
	})
}
