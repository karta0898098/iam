package middleware

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/karta0898098/iam/pkg/errors"
)

// NewErrorHandlingMiddleware handles panic error
func NewErrorHandlingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					trace := make([]byte, 4096)
					runtime.Stack(trace, true)
					traceID := c.Request().Header.Get(echo.HeaderXRequestID)
					customFields := map[string]interface{}{
						"url":         c.Request().RequestURI,
						"stack_error": string(trace),
						"trace_id":    traceID,
					}
					err, ok := r.(error)
					if !ok {
						if err == nil {
							err = fmt.Errorf("%v", r)
						} else {
							err = fmt.Errorf("%v", err)
						}
					}
					logger := log.With().Fields(customFields).Logger()
					logger.Error().Msgf("http: unknown error: %v", err)

					status, payload := errors.ErrInternal.ToViewModel()
					_ = c.JSON(status, payload)
				}
			}()
			return next(c)
		}
	}
}
