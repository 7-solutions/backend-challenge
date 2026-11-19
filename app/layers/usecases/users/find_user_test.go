package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	_mockuserrepo "github.com/chaihanij/backend-challenge/app/mocks/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"go.uber.org/mock/gomock"
)

func TestUseCase_FindUser(t *testing.T) {
	testhelpers.Init()
	ctx := context.Background()

	t.Run("Returns Matching Users When Repository Returns Results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		name := "Tester"
		input := &entities.UserFilter{Name: pointer.ToPointer(name)}
		expected := []entities.User{
			{ID: "1", Name: "Tester", Email: "tester1@test.com"},
			{ID: "2", Name: "Tester", Email: "tester2@test.com"},
		}

		mockUserRepo.EXPECT().FindUser(gomock.Any(), input).Return(expected, nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindUser(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(result) != len(expected) {
			t.Errorf("expected %d users, got %d", len(expected), len(result))
		}
	})

	t.Run("Returns Empty Slice When No Users Match Filter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		email := "missing@test.com"
		input := &entities.UserFilter{Email: pointer.ToPointer(email)}
		expected := []entities.User{}

		mockUserRepo.EXPECT().FindUser(gomock.Any(), input).Return(expected, nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindUser(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected non-nil result slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("expected empty result, got %d items", len(result))
		}
	})

	t.Run("Returns Users When Filter Is Nil And Repository Succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		expected := []entities.User{{ID: "1", Name: "Tester", Email: "tester@test.com"}}

		mockUserRepo.EXPECT().FindUser(gomock.Any(), nil).Return(expected, nil).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindUser(ctx, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 user, got %d", len(result))
		}
	})

	t.Run("Returns Error When Repository Find Fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		name := "Tester"
		input := &entities.UserFilter{Name: pointer.ToPointer(name)}
		expectedErr := errors.New("database unavailable")

		mockUserRepo.EXPECT().FindUser(gomock.Any(), input).Return(nil, expectedErr).Times(1)

		u := users.New(mockUserRepo, testhelpers.Logger)
		result, err := u.FindUser(ctx, input)
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
