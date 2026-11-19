package users_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	_mockuserrepo "github.com/chaihanij/backend-challenge/app/mocks/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/mock/gomock"
)

func TestUseCase_FindOneUser(t *testing.T) {
	testhelpers.Init()
	ctx := context.Background()

	t.Run("Returns User When User Exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		id := "f20db9be8e6f361ec3a93d9135175c4a"
		expectedUser := &entities.User{ID: id, Name: "Tester", Email: "test@test.com"}
		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), id).Return(expectedUser, nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindOneUser(ctx, id)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected non-nil result, got nil")
		}
		if result != nil && result.ID != id {
			t.Errorf("expected id %s, got %s", id, result.ID)
		}
	})

	t.Run("Returns NotFoundError When User Is Missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		id := "missing-user"
		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), id).Return(nil, mongo.ErrNoDocuments).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindOneUser(ctx, id)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result, got %#v", result)
		}
		if _, ok := errors.AsType[*customerrors.NotFoundError](err); !ok {
			t.Errorf("expected not found error, got %v", err)
		}
	})

	t.Run("Returns NotFoundError For Wrapped NoDocuments Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		id := "missing-user"
		wrappedErr := fmt.Errorf("lookup failed: %w", mongo.ErrNoDocuments)
		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), id).Return(nil, wrappedErr).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindOneUser(ctx, id)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result, got %#v", result)
		}
		if _, ok := errors.AsType[*customerrors.NotFoundError](err); !ok {
			t.Errorf("expected not found error, got %v", err)
		}
	})

	t.Run("Returns Repository Error When Lookup Fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		id := "f20db9be8e6f361ec3a93d9135175c4a"
		expectedErr := errors.New("database unavailable")
		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), id).Return(nil, expectedErr).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindOneUser(ctx, id)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result, got %#v", result)
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}
