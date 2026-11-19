package auth

import (
	"github.com/chaihanij/backend-challenge/app/constants"
	middleware "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type Auth struct {
	mock.Mock
	userId string
}

func (a *Auth) Auth(c *gin.Context) {
	c.Set(constants.UserIdKey, a.userId)
	c.Next()
}

func NewMockAuth() middleware.Auth {
	return &Auth{
		userId: "mocked-user-id",
	}
}
