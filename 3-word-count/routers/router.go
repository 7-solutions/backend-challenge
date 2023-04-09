package routers

import (
	"net/http"
	"word-count/handlers"

	"github.com/labstack/echo/v4"
)

func InitRouter() *echo.Echo {
	// Init echo
	e := echo.New()

	e.GET("/health", healthCheck)

	beef := e.Group("/beef")

	beef.GET("/summary", handlers.GetSummary)

	return e

}

func healthCheck(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, echo.Map{"message": "Service is Running !!"}, "	")
}
