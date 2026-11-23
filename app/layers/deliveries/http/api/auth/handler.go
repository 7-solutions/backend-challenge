package auth

import (
	_auth "github.com/chaihanij/backend-challenge/app/layers/usecases/auth"
	_users "github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	"github.com/gin-gonic/gin"
)

type handler struct {
	authUseCase  _auth.UseCase
	usersUseCase _users.UseCase
}

func New(ginEngin *gin.Engine, authUseCase _auth.UseCase, usersUseCase _users.UseCase) {
	h := &handler{
		authUseCase:  authUseCase,
		usersUseCase: usersUseCase,
	}

	auth := ginEngin.Group("/api/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}

}
