package health

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler(ginEngine *gin.Engine) *Handler {
	h := &Handler{}

	router := ginEngine
	{
		router.Handle("GET", "/api/health", h.Health)
	}

	return h
}
