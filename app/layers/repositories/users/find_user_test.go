package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	mockmongo "github.com/chaihanij/backend-challenge/app/test_helpers/mongo"
)

func TestRepo_FindUser(t *testing.T) {
	testhelpers.Init()

	t.Run("Success", func(t *testing.T) {

		id := bson.NewObjectID()
		now := time.Now()

		fakeCol := &mockmongo.FakeCollection{
			FindFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
				docs := []any{
					bson.M{
						"_id":        id,
						"name":       "Tester",
						"email":      "test@example.com",
						"password":   "hashed",
						"created_at": now,
						"updated_at": now,
					},
				}
				return mongo.NewCursorFromDocuments(docs, nil, nil)
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		name := "Tester"
		result, err := repo.FindUser(context.Background(), &entities.UserFilter{Name: &name})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 1 {
			t.Fatalf("expected 1 user, got %d", len(result))
		}
		if result[0].ID != id.Hex() {
			t.Fatalf("expected id %s, got %s", id.Hex(), result[0].ID)
		}
		if result[0].Email != "test@example.com" {
			t.Fatalf("expected email test@example.com, got %s", result[0].Email)
		}
	})

	t.Run("Success: No Match", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			FindFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
				return mongo.NewCursorFromDocuments([]any{}, nil, nil)
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		name := "Nobody"
		result, err := repo.FindUser(context.Background(), &entities.UserFilter{Name: &name})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected 0 users, got %d", len(result))
		}
	})

	t.Run("Failure: Filter Error", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			FindFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
				t.Fatalf("Find should NOT be called when filter build fails")
				return nil, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		// an all-nil filter produces an empty mongo filter, which NewUserFilter rejects.
		_, err := repo.FindUser(context.Background(), &entities.UserFilter{})
		if err == nil {
			t.Fatalf("expected filter error, got nil")
		}
	})

	t.Run("Failure: Find Error", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			FindFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
				return nil, errors.New("find failed")
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		name := "Tester"
		_, err := repo.FindUser(context.Background(), &entities.UserFilter{Name: &name})
		if err == nil {
			t.Fatalf("expected find error, got nil")
		}
	})

	t.Run("Failure: Decode Error", func(t *testing.T) {

		fakeCol := &mockmongo.FakeCollection{
			FindFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
				docs := []any{
					bson.M{
						"_id":  "not-an-object-id", // wrong type, should fail to decode into primitive.ObjectID
						"name": "Tester",
					},
				}
				return mongo.NewCursorFromDocuments(docs, nil, nil)
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		name := "Tester"
		_, err := repo.FindUser(context.Background(), &entities.UserFilter{Name: &name})
		if err == nil {
			t.Fatalf("expected decode error, got nil")
		}
	})
}
