package auth

import (
	_auth "github.com/chaihanij/backend-challenge/app/layers/usecases/auth"
	"github.com/gin-gonic/gin"
)

type handler struct {
	authUseCase _auth.UseCase
}

func New(ginEngin *gin.Engine, authUseCase _auth.UseCase) {
	h := &handler{
		authUseCase: authUseCase,
	}

	auth := ginEngin.Group("/auth")
	{
		auth.POST("/login", h.Login)
	}

}
