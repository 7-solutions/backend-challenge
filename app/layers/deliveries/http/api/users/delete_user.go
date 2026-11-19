package users

import (
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	httpvalidators "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/gin-gonic/gin"
)

func (h *handler) DeleteUser(c *gin.Context) {
	var req struct {
		ID string `uri:"id,required" binding:"required,hexadecimal"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		err = apperrors.WrapValidationError(
			httpvalidators.FieldErrorItems(httpvalidators.MapValidationErrorsSlice(err)).GetFieldErrorItems(),
			err,
		)
		_ = c.Error(err)
		return
	}

	err := h.usersUseCase.DeleteUser(c, req.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, nil)
}
