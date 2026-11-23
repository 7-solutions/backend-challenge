package auth

import (
	"github.com/chaihanij/backend-challenge/app/entities"
	usersmodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users/models"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) Register(c *gin.Context) {
	req, err := usersmodels.NewCreateUserRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := h.usersUseCase.CreateUser(c, &entities.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, usersmodels.UserResponseFromEntity(user))
}
