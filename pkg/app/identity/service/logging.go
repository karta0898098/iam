package service

import (
	"context"

	"github.com/karta0898098/iam/pkg/app/identity/entity"
)

type loggingMiddleware struct {
	next IdentityService `json:"-"`
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware() Middleware {
	return func(next IdentityService) IdentityService {
		return loggingMiddleware{next}
	}
}

func (lm loggingMiddleware) Signin(ctx context.Context, username string, password string, opt *SigninOption) (identity *entity.Identity, err error) {
	defer func() {
		// lm.logger.Log(
		// 	"method", "Signin",
		// 	"username", username,
		// 	"err", err,
		// )
	}()

	return lm.next.Signin(ctx, username, password, opt)
}

func (lm loggingMiddleware) Signup(ctx context.Context, username string, password string, opts *SignupOption) (identity *entity.Identity, err error) {
	defer func() {
		// lm.logger.Log(
		// 	"method", "Signup",
		// 	"username", username,
		// 	"err", err,
		// )
	}()
	return lm.next.Signup(ctx, username, password, opts)
}
