package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	_mockuserrepo "github.com/chaihanij/backend-challenge/app/mocks/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"go.uber.org/mock/gomock"
)

func TestUseCase_CreateUser(t *testing.T) {
	testhelpers.Init()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		input := &entities.User{Name: "Tester", Email: "test@test.com", Password: "P@ssw0rd"}
		now := time.Now().UTC().Truncate(time.Millisecond)
		expected := &entities.User{ID: "f20db9be8e6f361ec3a93d9135175c4a", Name: "Tester", Email: "test@test.com", Password: "hashed", CreatedAt: pointer.ToPointer(now), UpdatedAt: pointer.ToPointer(now)}

		mockUserRepo.EXPECT().FindUser(gomock.Any(), &entities.UserFilter{Email: pointer.ToPointer(input.Email)}).Return([]entities.User{}, nil).Times(1)
		mockUserRepo.EXPECT().CreateUser(gomock.Any(), input).Return(expected, nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.CreateUser(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result == nil {
			t.Errorf("expected non-nil result, got nil")
		}
	})

	t.Run("Failure: Find User Exiting", func(t *testing.T) {
		input := &entities.User{Name: "Tester", Email: "test@test.com", Password: "P@ssw0rd"}
		mockUserRepo.EXPECT().FindUser(gomock.Any(), &entities.UserFilter{Email: pointer.ToPointer(input.Email)}).Return(nil, errors.New("find error")).Times(1)
		u := users.New(mockUserRepo, testhelpers.Logger)
		_, err := u.CreateUser(ctx, input)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Failure: User Already Exists", func(t *testing.T) {
		input := &entities.User{Name: "Tester", Email: "test@test.com", Password: "P@ssw0rd"}
		now := time.Now().UTC().Truncate(time.Millisecond)
		mockuser := entities.User{ID: "f20db9be8e6f361ec3a93d9135175c4a", Name: "Tester", Email: "test@test.com", Password: "hashed", CreatedAt: pointer.ToPointer(now), UpdatedAt: pointer.ToPointer(now)}

		mockUserRepo.EXPECT().FindUser(gomock.Any(), &entities.UserFilter{Email: pointer.ToPointer(input.Email)}).Return([]entities.User{mockuser}, nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		_, err := u.CreateUser(ctx, input)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if _, ok := errors.AsType[*customerrors.ConflictError](err); !ok {

			t.Errorf("expected conflict error, got %v", err)
		}
	})

	t.Run("Failure: User Creation Error", func(t *testing.T) {
		input := &entities.User{Name: "Tester", Email: "test@test.com", Password: "P@ssw0rd"}
		mockUserRepo.EXPECT().FindUser(gomock.Any(), &entities.UserFilter{Email: pointer.ToPointer(input.Email)}).Return([]entities.User{}, nil).Times(1)
		mockUserRepo.EXPECT().CreateUser(gomock.Any(), input).Return(nil, errors.New("creation error")).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.CreateUser(ctx, input)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result, got: %#v", result)
		}
	})
}
