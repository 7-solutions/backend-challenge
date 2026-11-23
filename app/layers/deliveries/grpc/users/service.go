package users

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/chaihanij/backend-challenge/app/entities"
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	usersv1 "github.com/chaihanij/backend-challenge/app/gen/proto/users/v1"
	usersusecase "github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements the users.v1 gRPC API.
type Service struct {
	usersv1.UnimplementedUserServiceServer
	usersUseCase usersusecase.UseCase
}

func NewService(usersUseCase usersusecase.UseCase) *Service {
	return &Service{usersUseCase: usersUseCase}
}

func (s *Service) CreateUser(ctx context.Context, req *usersv1.CreateUserRequest) (*usersv1.CreateUserResponse, error) {
	if req == nil || len(strings.TrimSpace(req.GetName())) < 3 || len(req.GetName()) > 255 || len(req.GetEmail()) > 255 || len(req.GetPassword()) < 6 {
		return nil, status.Error(codes.InvalidArgument, "name must be at least 3 characters, email is required, and password must be at least 6 characters")
	}
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil || strings.TrimSpace(req.GetEmail()) != req.GetEmail() {
		return nil, status.Error(codes.InvalidArgument, "email must be valid")
	}

	user, err := s.usersUseCase.CreateUser(ctx, &entities.User{
		Name: req.GetName(), Email: req.GetEmail(), Password: req.GetPassword(),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &usersv1.CreateUserResponse{User: userToProto(user)}, nil
}

func (s *Service) GetUser(ctx context.Context, req *usersv1.GetUserRequest) (*usersv1.GetUserResponse, error) {
	if req == nil || strings.TrimSpace(req.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.usersUseCase.FindOneUser(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &usersv1.GetUserResponse{User: userToProto(user)}, nil
}

func userToProto(user *entities.User) *usersv1.User {
	return &usersv1.User{Id: user.ID, Name: user.Name, Email: user.Email}
}

func toStatusError(err error) error {
	var conflict *apperrors.ConflictError
	var notFound *apperrors.NotFoundError
	var validation *apperrors.ValidationError
	switch {
	case errors.As(err, &conflict):
		return status.Error(codes.AlreadyExists, conflict.Error())
	case errors.As(err, &notFound):
		return status.Error(codes.NotFound, notFound.Error())
	case errors.As(err, &validation):
		return status.Error(codes.InvalidArgument, validation.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
