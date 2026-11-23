package models

import (
	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
)

type UserResponse struct {
	ID    string `json:"id" example:"1bf2e66a05118e911a36933125bbab4c"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

func UserResponseFromEntity(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}

type UsersResponse []UserResponse

func UsersResponseFromEntity(users []entities.User) UsersResponse {
	var usersResponse []UserResponse
	for _, user := range users {
		userResponse := UserResponseFromEntity(&user)
		usersResponse = append(usersResponse, pointer.ToValue(userResponse))
	}
	return usersResponse
}

type UserResponseJsonSwagger struct {
	Code      int            `json:"code" example:"200"`
	Message   string         `json:"message" example:"Success"`
	Data      UserResponse   `json:"data"`
	Meta      *response.Meta `json:"meta,omitempty"`
	Timestamp string         `json:"timestamp,omitempty"`
}

type UsersResponseJsonSwagger struct {
	Code      int            `json:"code" example:"200"`
	Message   string         `json:"message" example:"Success"`
	Data      UsersResponse  `json:"data"`
	Meta      *response.Meta `json:"meta"`
	Timestamp string         `json:"timestamp,omitempty"`
}
