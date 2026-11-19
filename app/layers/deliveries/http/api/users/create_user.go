package users

import (
	"github.com/chaihanij/backend-challenge/app/entities"
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users/models"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	httpvalidators "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/gin-gonic/gin"
)

func (h *handler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email,required" validate:"required,email,max=255,min=3,format=email"`
		Name     string `json:"name,required" validate:"required,min=3"`
		Password string `json:"password,required" validate:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		err = apperrors.WrapValidationError(
			httpvalidators.FieldErrorItems(httpvalidators.MapValidationErrorsSlice(err)).GetFieldErrorItems(),
			err,
		)
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
