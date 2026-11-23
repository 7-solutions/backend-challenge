package models

import (
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	httpvalidators "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/gin-gonic/gin"
)

type UpdateUserRequest struct {
	Email string `json:"email,required" binding:"required,email,max=255,min=3"`
	Name  string `json:"name,required" binding:"required,min=3"`
}

func NewUpdateUserRequest(c *gin.Context) (*UpdateUserRequest, error) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err = apperrors.WrapValidationError(
			httpvalidators.FieldErrorItems(httpvalidators.MapValidationErrorsSlice(err)).GetFieldErrorItems(),
			err,
		)
		return nil, err
	}
	return &req, nil
}
