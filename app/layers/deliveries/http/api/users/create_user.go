package users

import (
	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users/models"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) CreateUser(c *gin.Context) {
	req, err := models.NewCreateUserRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result, err := h.usersUseCase.CreateUser(c, &entities.User{Email: req.Email, Name: req.Name, Password: req.Password})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, models.UserResponseFromEntity(result))
}
