package middlewares

import (
	"strings"

	"github.com/chaihanij/backend-challenge/app/constants"
	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	_auth "github.com/chaihanij/backend-challenge/app/layers/usecases/auth"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"github.com/gin-gonic/gin"
)

type Auth interface {
	Auth(c *gin.Context)
}

type auth struct {
	authUseCase _auth.UseCase
}

func NewAuth(authUseCase _auth.UseCase) Auth {
	return &auth{
		authUseCase: authUseCase,
	}
}

func (a *auth) Auth(c *gin.Context) {

	authHeader := c.GetHeader(constants.AuthorizationHeaderKey)
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	if accessToken == "" {

		_ = c.Error(apperrors.NewUnauthorized(apperrors.ErrMissingAuthorizationMsg))
		c.Abort()
		return
	}

	userId, ok, err := a.authUseCase.ValidateToken(c, accessToken)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	if !ok {
		_ = c.Error(apperrors.NewUnauthorized("Invalid token"))
		c.Abort()
		return
	}

	c.Set(constants.UserIdKey, pointer.ToValue(userId))
	c.Next()
}
