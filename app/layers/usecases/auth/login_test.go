package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/usecases/auth"
	_mockuserrepo "github.com/chaihanij/backend-challenge/app/mocks/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"github.com/chaihanij/backend-challenge/app/util/password"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Login(t *testing.T) {
	testhelpers.Init()

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)

		ctx := context.Background()

		email := "tester@test.com"
		plain := "s3cret-passw0rd"

		hash, _ := password.HashPassword(plain)
		user := entities.User{
			Email:    email,
			Password: hash,
		}

		mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return([]entities.User{user}, nil).Times(1)

		u := auth.New(mockUserRepo, testhelpers.Logger, "secret", "iss", "aud", time.Hour)

		token, err := u.Login(ctx, email, plain)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == nil || len(pointer.ToValue(token)) == 0 {
			t.Errorf("expected token, got %v", token)
		}
	})

	t.Run("Failure: db error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)

		ctx := context.Background()

		email := "tester@test.com"
		plain := "s3cret-passw0rd"

		mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)

		u := auth.New(mockUserRepo, testhelpers.Logger, "secret", "iss", "aud", time.Hour)
		_, err := u.Login(ctx, email, plain)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Failure: user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)

		ctx := context.Background()

		email := "tester@test.com"
		plain := "s3cret-passw0rd"

		mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

		u := auth.New(mockUserRepo, testhelpers.Logger, "secret", "iss", "aud", time.Hour)
		_, err := u.Login(ctx, email, plain)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

	})

	t.Run("Failure: wrong password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)

		ctx := context.Background()

		email := "tester@test.com"
		plain := "s3cret-passw0rd"

		hash, _ := password.HashPassword("wrong-password")
		user := entities.User{
			Email:    email,
			Password: hash,
		}

		mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return([]entities.User{user}, nil).Times(1)

		u := auth.New(mockUserRepo, testhelpers.Logger, "secret", "iss", "aud", time.Hour)

		token, err := u.Login(ctx, email, plain)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if token != nil {
			t.Errorf("expected no token, got %v", token)
		}
	})
}
