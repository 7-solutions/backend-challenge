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

	api := ginEngin.Group("/api/users")
	api.POST("", h.CreateUser)

	protected := api.Group("").Use(authMiddleware.Auth)
	{
		protected.GET("", h.FindUser)
		protected.PUT("/:id", h.UpdateUser)
		protected.GET("/:id", h.FindOneUser)
		protected.DELETE("/:id", h.DeleteUser)
	}
}
