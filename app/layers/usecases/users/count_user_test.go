package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	_mockuserrepo "github.com/chaihanij/backend-challenge/app/mocks/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"go.uber.org/mock/gomock"
)

func TestUseCase_CountUser(t *testing.T) {
	testhelpers.Init()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {

		mockUserRepo.EXPECT().CountUser(gomock.Any()).Return(int64(10), nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)

		count, err := u.CountUser(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if count != 10 {
			t.Errorf("expected count 10, got %d", count)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		mockUserRepo.EXPECT().CountUser(gomock.Any()).Return(int64(0), errors.New("failed to count users")).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)

		count, err := u.CountUser(ctx)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}
	})
}
