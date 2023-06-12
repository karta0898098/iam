package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLoggerMiddleware make request context with logger
// then all context can get logger from context
func NewLoggerMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			// set context trace id
			ctx := c.Request().Context()

			ctx = logger.WithContext(ctx)

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func NewLoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceID := c.Request().Header.Get(echo.HeaderXRequestID)

			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			var logger *zerolog.Event

			status := c.Response().Status
			if status >= 500 {
				logger = log.Error()
			} else if status >= 400 {
				logger = log.Info()
			} else {
				logger = log.Info()
			}

			logger.
				Str("method", c.Request().Method).
				Str("uri", c.Request().RequestURI).
				Str("trace_id", traceID).
				Str("latency_human", stop.Sub(start).String()).
				Int("status", status).
				Msg("http access log.")

			return nil
		}
	}
}

// RecordErrorMiddleware provide error middleware
// useful to record error occur
func RecordErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logFields := map[string]interface{}{}

				// record request data
				req := c.Request()
				{
					logFields["method"] = req.Method
					logFields["uri"] = req.RequestURI
				}
				ctx := req.Context()

				// record response data
				resp := c.Response()
				resp.After(func() {
					logFields["status"] = resp.Status
					// according http status to decide log level
					logger := log.Ctx(ctx).With().Fields(logFields).Logger()
					if resp.Status >= http.StatusInternalServerError {
						logger.Error().Msgf("%+v", err)
					} else if resp.Status >= http.StatusBadRequest {
						logger.Debug().Msgf("%+v", err)
					} else {
						logger.Debug().Msgf("%+v", err)
					}
				})
			}
			return err
		}
	}
}
