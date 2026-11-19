package models

import (
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	httpvalidators "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,min=3,max=255"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewLoginRequest(c *gin.Context) (*LoginRequest, error) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err = apperrors.WrapValidationError(
			httpvalidators.FieldErrorItems(httpvalidators.MapValidationErrorsSlice(err)).GetFieldErrorItems(),
			err,
		)
		return nil, err
	}
	return &req, nil
}
