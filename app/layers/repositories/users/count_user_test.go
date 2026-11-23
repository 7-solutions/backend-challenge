package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chaihanij/backend-challenge/app/layers/repositories/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	mockmongo "github.com/chaihanij/backend-challenge/app/test_helpers/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestRepo_CountUser(t *testing.T) {
	testhelpers.Init()

	t.Run("Success: Returns Total Users", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			EstimatedDocumentCountFn: func(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error) {
				return 12, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		count, err := repo.CountUser(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 12 {
			t.Fatalf("expected count 12, got %d", count)
		}
	})

	t.Run("Success: Returns Zero When No Users", func(t *testing.T) {
		fakeCol := &mockmongo.FakeCollection{
			EstimatedDocumentCountFn: func(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error) {
				return 0, nil
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		count, err := repo.CountUser(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Fatalf("expected count 0, got %d", count)
		}
	})

	t.Run("Failure: Returns Error When Count Fails", func(t *testing.T) {
		expectedErr := errors.New("count failed")
		fakeCol := &mockmongo.FakeCollection{
			EstimatedDocumentCountFn: func(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error) {
				return 0, expectedErr
			},
		}

		repo := users.Wrap(fakeCol, testhelpers.Config.MongoDBRequestTimeout)

		count, err := repo.CountUser(context.Background())
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
		if count != 0 {
			t.Fatalf("expected count 0 on error, got %d", count)
		}
	})
}
