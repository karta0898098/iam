package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoute(route *echo.Echo, handler *Handler) {
	route.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	api := route.Group("/api")
	{
		v1 := api.Group("/v1")
		v1.POST("/login", handler.LoginEndpoint)
		v1.POST("/signup", handler.SignupEndpoint)
	}
}
