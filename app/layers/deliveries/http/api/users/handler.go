package users

import (
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	_users "github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	"github.com/gin-gonic/gin"
)

type handler struct {
	usersUseCase _users.UseCase
}

func New(ginEngin *gin.Engine, usersUseCase _users.UseCase, authMiddleware middlewares.Auth) {
	h := &handler{
		usersUseCase: usersUseCase,
	}

	api := ginEngin.Group("/api/users").Use(authMiddleware.Auth)
	{
		api.GET("", h.FindUser)
		api.POST("", h.CreateUser)
		api.PUT("/:id", h.UpdateUser)
		api.GET("/:id", h.FindOneUser)
		api.DELETE("/:id", h.DeleteUser)
	}
}
