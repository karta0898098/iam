package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/karta0898098/iam/pkg/errors"
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

// ValidateMiddleware returns an endpoint middleware that validate each invocation,
// and the resulting error, if any.
func ValidateMiddleware(validate *validator.Validate, trans ut.Translator) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = validate.StructCtx(ctx, request)
			if err != nil {

				var details []errors.Detail
				errs := err.(validator.ValidationErrors)
				for _, e := range errs {

					details = append(details, errors.Detail{
						Reason: e.(validator.FieldError).Translate(trans),
					})
				}

				return nil, errors.Wrap(errors.ErrInvalidInput.WithDetails(details...), "invalid input")
			}

			return next(ctx, request)
		}
	}
}
