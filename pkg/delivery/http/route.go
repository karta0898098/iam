package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoute(route *echo.Echo, handler *Handler) {
	route.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
}
