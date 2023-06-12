package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog/log"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger := log.Ctx(ctx)

			defer func(begin time.Time) {
				if err == nil {
					logger.Info().
						Str("method", method).
						Dur("took", time.Since(begin)).
						Msg("endpoint metrics")
				} else {
					logger.Error().
						Str("method", method).
						Dur("took", time.Since(begin)).
						Err(err).
						Msg("endpoint metrics")
				}
			}(time.Now())
			return next(ctx, request)
		}
	}
}
