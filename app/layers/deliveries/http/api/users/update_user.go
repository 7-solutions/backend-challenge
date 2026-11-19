package users

import (
	"github.com/chaihanij/backend-challenge/app/entities"
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users/models"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	httpvalidators "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/gin-gonic/gin"
)

func (h *handler) UpdateUser(c *gin.Context) {
	var reqUri struct {
		ID string `uri:"id,required" binding:"required,hexadecimal"`
	}

	if err := c.ShouldBindUri(&reqUri); err != nil {
		err = apperrors.WrapValidationError(
			httpvalidators.FieldErrorItems(httpvalidators.MapValidationErrorsSlice(err)).GetFieldErrorItems(),
			err,
		)
		_ = c.Error(err)
		return
	}

	reqBody, err := models.NewUpdateUserRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result, err := h.usersUseCase.UpdateUser(c.Request.Context(), &entities.User{ID: reqUri.ID, Email: reqBody.Email, Name: reqBody.Name})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, models.UserResponseFromEntity(result))
}
