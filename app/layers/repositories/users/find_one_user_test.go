package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	mockmongo "github.com/chaihanij/backend-challenge/app/test_helpers/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestRepo_FindOneUse(t *testing.T) {
	testhelpers.Init()

	t.Run("Success", func(t *testing.T) {
		id := bson.NewObjectID()
		name := "Tester"
		email := "test@example.com"
		password := "hashed"
		now := time.Now().UTC().Truncate(time.Millisecond)
		fakeCol := &mockmongo.FakeCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				doc := bson.M{"_id": id, "name": name, "email": email, "password": password, "created_at": now, "updated_at": now}
				return mongo.NewSingleResultFromDocument(doc, nil, nil)
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		result, err := repo.FindOneUser(context.Background(), id.Hex())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertUser(t, result, id.Hex(), name, email, password, now)
	})

	t.Run("Success: No Match", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		result, err := repo.FindOneUser(context.Background(), bson.NewObjectID().Hex())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "document is nil" {
			t.Fatalf("expected document is nil error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("expected nil result, got: %#v", result)
		}
	})

	t.Run("Failure: Filter Error", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				t.Fatal("FindOne should NOT be called when filter build fails")
				return nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		_, err := repo.FindOneUser(context.Background(), "not-an-object-id")
		if err == nil {
			t.Fatal("expected filter error, got nil")
		}
	})

	t.Run("Failure: Find Error", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(nil, errors.New("find failed"), nil)
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		_, err := repo.FindOneUser(context.Background(), bson.NewObjectID().Hex())
		if err == nil {
			t.Fatal("expected find error, got nil")
		}
	})

	t.Run("Failure: Decode Error", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				doc := bson.M{"_id": "not-an-object-id", "name": "Tester"}
				return mongo.NewSingleResultFromDocument(doc, nil, nil)
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		_, err := repo.FindOneUser(context.Background(), bson.NewObjectID().Hex())
		if err == nil {
			t.Fatal("expected decode error, got nil")
		}
	})
}

func assertUser(t *testing.T, result *entities.User, id, name, email, password string, ts time.Time) {
	t.Helper()

	if result == nil {
		t.Fatal("expected non-nil user")
	}
	if result.ID != id {
		t.Fatalf("expected id %q, got %q", id, result.ID)
	}
	if result.Name != name {
		t.Fatalf("expected name %s, got %s", name, result.Name)
	}
	if result.Email != email {
		t.Fatalf("expected email %s, got %s", email, result.Email)
	}
	if result.Password != password {
		t.Fatalf("expected password %s, got %s", password, result.Password)
	}
	if result.CreatedAt == nil || !result.CreatedAt.Equal(ts) {
		t.Fatalf("expected created_at %s, got %s", ts, result.CreatedAt)
	}
	if result.UpdatedAt == nil || !result.UpdatedAt.Equal(ts) {
		t.Fatalf("expected updated_at %s, got %s", ts, result.UpdatedAt)
	}
}
