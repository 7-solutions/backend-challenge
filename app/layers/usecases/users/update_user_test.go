package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	_mockuserrepo "github.com/chaihanij/backend-challenge/app/mocks/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"go.uber.org/mock/gomock"
)

func TestUseCase_UpdateUser(t *testing.T) {
	testhelpers.Init()
	ctx := context.Background()

	t.Run("Returns Updated User When Update Succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		input := &entities.User{ID: "f20db9be8e6f361ec3a93d9135175c4a", Name: "Updated Tester", Email: "updated@test.com"}
		updated := &entities.User{ID: input.ID, Name: input.Name, Email: input.Email}

		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), input.ID).Return(&entities.User{ID: input.ID}, nil).Times(1)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any(), input).Return(nil).Times(1)
		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), input.ID).Return(updated, nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.UpdateUser(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected non-nil result, got nil")
		}
		if result != nil && result.Name != input.Name {
			t.Errorf("expected name %s, got %s", input.Name, result.Name)
		}
	})

	t.Run("Returns Error When Lookup Before Update Fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		input := &entities.User{ID: "missing-user", Name: "Tester", Email: "test@test.com"}
		expectedErr := errors.New("find before update error")

		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), input.ID).Return(nil, expectedErr).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.UpdateUser(ctx, input)
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

	t.Run("Returns Error When Update Operation Fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		input := &entities.User{ID: "f20db9be8e6f361ec3a93d9135175c4a", Name: "Tester", Email: "test@test.com"}
		expectedErr := errors.New("update error")

		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), input.ID).Return(&entities.User{ID: input.ID}, nil).Times(1)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any(), input).Return(expectedErr).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.UpdateUser(ctx, input)
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

	t.Run("Returns Error When Lookup After Update Fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		input := &entities.User{ID: "f20db9be8e6f361ec3a93d9135175c4a", Name: "Tester", Email: "test@test.com"}
		expectedErr := errors.New("find after update error")

		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), input.ID).Return(&entities.User{ID: input.ID}, nil).Times(1)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any(), input).Return(nil).Times(1)
		mockUserRepo.EXPECT().FindOneUser(gomock.Any(), input.ID).Return(nil, expectedErr).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.UpdateUser(ctx, input)
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
