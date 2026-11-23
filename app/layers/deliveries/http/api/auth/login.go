package auth

import (
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/auth/models"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) Login(c *gin.Context) {

	req, err := models.NewLoginRequest(c)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	token, err := h.authUseCase.Login(c, req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	response.Success(c, token)
}
