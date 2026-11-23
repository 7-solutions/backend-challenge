package models

import (
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	httpvalidators "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email,min=3,max=255"`
	Name     string `json:"name" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewCreateUserRequest(c *gin.Context) (*CreateUserRequest, error) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err = apperrors.WrapValidationError(
			httpvalidators.FieldErrorItems(httpvalidators.MapValidationErrorsSlice(err)).GetFieldErrorItems(),
			err,
		)
		return nil, err
	}
	return &req, nil
}
