package users

import (
	"github.com/chaihanij/backend-challenge/app/entities"
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users/models"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	httpvalidators "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/gin-gonic/gin"
)

func (h *handler) FindUser(c *gin.Context) {
	var req struct {
		Email *string `form:"email" binding:"omitempty,email,max=255,min=3"`
		Name  *string `form:"name" binding:"omitempty,min=3"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		err = apperrors.WrapValidationError(
			httpvalidators.FieldErrorItems(httpvalidators.MapValidationErrorsSlice(err)).GetFieldErrorItems(),
			err,
		)
		_ = c.Error(err)
		return
	}

	results, err := h.usersUseCase.FindUser(c, &entities.UserFilter{Email: req.Email, Name: req.Name})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, models.UsersResponseFromEntity(results))
}
