package service

import (
	"context"

	"github.com/rs/xid"

	"github.com/karta0898098/iam/pkg/app/identity/entity"
	"github.com/karta0898098/iam/pkg/app/identity/repository"
	"github.com/karta0898098/iam/pkg/errors"
)

var _ IdentityService = &Impl{}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(IdentityService) IdentityService

// IdentityService define identity service
type IdentityService interface {
	// Signin verify user account
	Signin(
		ctx context.Context,
		username string,
		password string,
		opt *SigninOption,
	) (identity *entity.Identity, err error)

	// Signup sign up a new user
	Signup(
		ctx context.Context,
		username string,
		password string,
		opts *SignupOption,
	) (identity *entity.Identity, err error)
}

type Impl struct {
	repo repository.Repository
}

func New(repo repository.Repository) IdentityService {
	var svc IdentityService
	svc = &Impl{
		repo: repo,
	}
	svc = LoggingMiddleware()(svc)

	return svc
}

func (srv *Impl) Signin(
	ctx context.Context,
	username string,
	password string,
	opt *SigninOption,
) (identity *entity.Identity, err error) {
	user, err := srv.repo.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(password) {
		return nil, errors.Wrapf(
			errors.ErrUnauthorized,
			"user=%s validate password failed",
			user.ID,
		)
	}

	if !user.IsActive() {
		return nil, errors.Wrapf(
			errors.ErrUnauthorized,
			"user=%s is not active status=%v",
			user.ID, user.Status,
		)
	}

	session := entity.NewSession(
		xid.New().String(),
		user.ID,
		opt.IPAddress,
		opt.Platform,
		entity.WithDevice(opt.Device),
	)

	err = srv.repo.StoreSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &entity.Identity{
		User:    user,
		Session: session,
	}, nil
}

func (srv *Impl) Signup(
	ctx context.Context,
	username string,
	password string,
	opt *SignupOption,
) (identity *entity.Identity, err error) {

	// build a new user
	newUser, err := entity.NewUser(
		xid.New().String(),
		username,
		password,
		entity.WithEmail(opt.Email),
		entity.WithNickname(opt.Nickname),
	)
	if err != nil {
		return nil, err
	}

	// check user already exist or not
	user, err := srv.repo.FindUserByUsername(ctx, username)
	if err != nil {
		if !errors.Is(errors.ErrResourceNotFound, err) {
			return nil, err
		}
	}

	if user != nil {
		return nil, errors.Wrapf(
			errors.ErrConflict,
			"username=%v already exist",
			username,
		)
	}

	err = srv.repo.StoreUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// create session to record
	session := entity.NewSession(
		xid.New().String(),
		newUser.ID,
		opt.IPAddress,
		opt.Platform,
		entity.WithDevice(opt.Device),
	)

	err = srv.repo.StoreSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &entity.Identity{
		User:    newUser,
		Session: session,
	}, nil
}
