package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	usersv1 "github.com/chaihanij/backend-challenge/app/gen/proto/users/v1"
	usersgrpc "github.com/chaihanij/backend-challenge/app/layers/deliveries/grpc/users"
	_mockusers "github.com/chaihanij/backend-challenge/app/mocks/users"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name string
		req  *usersv1.CreateUserRequest
	}{
		{"nil request", nil},
		{"short name", &usersv1.CreateUserRequest{Name: "ab", Email: "a@b.com", Password: "123456"}},
		{"long name", &usersv1.CreateUserRequest{Name: string(make([]byte, 256)), Email: "a@b.com", Password: "123456"}},
		{"long email", &usersv1.CreateUserRequest{Name: "Alice", Email: string(make([]byte, 256)), Password: "123456"}},
		{"short password", &usersv1.CreateUserRequest{Name: "Alice", Email: "a@b.com", Password: "12345"}},
		{"invalid email", &usersv1.CreateUserRequest{Name: "Alice", Email: "invalid", Password: "123456"}},
		{"email with surrounding whitespace", &usersv1.CreateUserRequest{Name: "Alice", Email: " a@b.com", Password: "123456"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			uc := _mockusers.NewMockUseCase(ctrl)
			resp, err := usersgrpc.NewService(uc).CreateUser(context.Background(), tt.req)
			if resp != nil {
				t.Fatalf("expected nil response, got %v", resp)
			}
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("expected InvalidArgument, got %v", err)
			}
		})
	}

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		uc := _mockusers.NewMockUseCase(ctrl)
		input := &entities.User{Name: "Alice", Email: "alice@example.com", Password: "secret1"}
		created := &entities.User{ID: "user-1", Name: input.Name, Email: input.Email, Password: "hashed"}
		uc.EXPECT().CreateUser(gomock.Any(), input).Return(created, nil)

		resp, err := usersgrpc.NewService(uc).CreateUser(context.Background(), &usersv1.CreateUserRequest{
			Name: input.Name, Email: input.Email, Password: input.Password,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.GetUser().GetId() != created.ID || resp.GetUser().GetName() != created.Name || resp.GetUser().GetEmail() != created.Email {
			t.Fatalf("unexpected user response: %v", resp.GetUser())
		}
	})

	t.Run("maps use case errors", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			err  error
			code codes.Code
		}{
			{"conflict", apperrors.NewConflictError("users", "duplicate"), codes.AlreadyExists},
			{"validation", apperrors.NewValidationError(map[string]string{"name": "required"}), codes.InvalidArgument},
			{"internal", errors.New("database details"), codes.Internal},
		} {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				uc := _mockusers.NewMockUseCase(ctrl)
				uc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, tc.err)
				resp, err := usersgrpc.NewService(uc).CreateUser(context.Background(), &usersv1.CreateUserRequest{Name: "Alice", Email: "a@b.com", Password: "123456"})
				if resp != nil || status.Code(err) != tc.code {
					t.Fatalf("expected nil response and %s, got response=%v error=%v", tc.code, resp, err)
				}
			})
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("rejects missing ID", func(t *testing.T) {
		for _, req := range []*usersv1.GetUserRequest{nil, {}, {Id: "  "}} {
			ctrl := gomock.NewController(t)
			uc := _mockusers.NewMockUseCase(ctrl)
			_, err := usersgrpc.NewService(uc).GetUser(context.Background(), req)
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("expected InvalidArgument, got %v", err)
			}
		}
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		uc := _mockusers.NewMockUseCase(ctrl)
		user := &entities.User{ID: "user-1", Name: "Alice", Email: "alice@example.com", Password: "secret"}
		uc.EXPECT().FindOneUser(gomock.Any(), "user-1").Return(user, nil)
		resp, err := usersgrpc.NewService(uc).GetUser(context.Background(), &usersv1.GetUserRequest{Id: "user-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.GetUser().GetId() != user.ID || resp.GetUser().GetName() != user.Name || resp.GetUser().GetEmail() != user.Email {
			t.Fatalf("unexpected user response: %v", resp.GetUser())
		}
	})

	t.Run("maps not found and internal errors", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			err  error
			code codes.Code
		}{
			{"not found", apperrors.NewNotFoundError("users", "missing"), codes.NotFound},
			{"internal", errors.New("database details"), codes.Internal},
		} {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				uc := _mockusers.NewMockUseCase(ctrl)
				uc.EXPECT().FindOneUser(gomock.Any(), "user-1").Return(nil, tc.err)
				resp, err := usersgrpc.NewService(uc).GetUser(context.Background(), &usersv1.GetUserRequest{Id: "user-1"})
				if resp != nil || status.Code(err) != tc.code {
					t.Fatalf("expected nil response and %s, got response=%v error=%v", tc.code, resp, err)
				}
			})
		}
	})
}
