package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/karta0898098/iam/pkg/app/identity/entity"
	"github.com/karta0898098/iam/pkg/app/identity/mocks"
	"github.com/karta0898098/iam/pkg/app/identity/repository"
	"github.com/karta0898098/iam/pkg/app/identity/service"
	"github.com/karta0898098/iam/pkg/errors"
)

func TestImpl_Signin(t *testing.T) {
	type args struct {
		username string
		password string
		opt      *service.SigninOption
	}
	tests := []struct {
		name     string
		repo     repository.Repository
		args     args
		expected *entity.Identity
		err      error
	}{
		{
			name: "Success",
			repo: func() repository.Repository {
				user, _ := entity.NewUser(
					"MOCK-USER-ID",
					"Username",
					"A12345678",
				)
				user.Status = entity.UserAccountStatusActive
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindUserByUsername(mock.Anything, mock.Anything).
					Return(user, nil)

				repo.EXPECT().
					StoreSession(mock.Anything, mock.Anything).
					Return(nil)
				return repo
			}(),
			args: args{
				username: "Username",
				password: "A12345678",
				opt: &service.SigninOption{
					IPAddress: "127.0.0.1",
					Platform:  "web",
					Device: entity.Device{
						Model:     "Intel Mac OS X",
						Name:      "Macintosh",
						OSVersion: "13_4",
					},
				},
			},
			expected: &entity.Identity{
				User: &entity.User{
					ID:       "MOCK-USER-ID",
					Username: "Username",
				},
			},
			err: nil,
		},
		{
			name: "Wrong Password",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindUserByUsername(mock.Anything, mock.Anything).
					Return(entity.NewUser(
						"MOCK-USER-ID",
						"Username",
						"A12345678",
					))
				return repo
			}(),
			args: args{
				username: "Username",
				password: "WrongPassword",
				opt: &service.SigninOption{
					IPAddress: "127.0.0.1",
					Platform:  "web",
					Device: entity.Device{
						Model:     "Intel Mac OS X",
						Name:      "Macintosh",
						OSVersion: "13_4",
					},
				},
			},
			expected: nil,
			err:      errors.ErrUnauthorized,
		},
		{
			name: "User is not active",
			repo: func() repository.Repository {
				user, _ := entity.NewUser(
					"MOCK-USER-ID",
					"Username",
					"A12345678",
				)
				user.Status = entity.UserAccountStatusSuspend
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindUserByUsername(mock.Anything, mock.Anything).
					Return(user, nil)
				return repo
			}(),
			args: args{
				username: "Username",
				password: "A12345678",
				opt: &service.SigninOption{
					IPAddress: "127.0.0.1",
					Platform:  "web",
					Device: entity.Device{
						Model:     "Intel Mac OS X",
						Name:      "Macintosh",
						OSVersion: "13_4",
					},
				},
			},
			expected: nil,
			err:      errors.ErrUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			srv := service.New(tt.repo)
			actual, err := srv.Signin(
				ctx,
				tt.args.username,
				tt.args.password,
				tt.args.opt,
			)
			if tt.err != nil || err != nil {
				if !errors.Is(tt.err, err) {
					t.Errorf("Signin() error = %v, wantErr %v", err, tt.err)
				}
				return
			}

			assert.Equal(t, tt.expected.User.ID, actual.User.ID)
			assert.Equal(t, tt.expected.User.Username, actual.User.Username)
			assert.NotEmpty(t, actual.Session.ID)
		})
	}
}

func TestImpl_Signup(t *testing.T) {
	type args struct {
		username string
		password string
		opt      *service.SignupOption
	}
	tests := []struct {
		name     string
		repo     repository.Repository
		args     args
		expected *entity.Identity
		err      error
	}{
		{
			name: "Success",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindUserByUsername(mock.Anything, mock.Anything).
					Return(nil, errors.ErrResourceNotFound)

				repo.EXPECT().
					StoreUser(mock.Anything, mock.Anything).
					Return(nil)

				repo.EXPECT().
					StoreSession(mock.Anything, mock.Anything).
					Return(nil)
				return repo
			}(),
			args: args{
				username: "Username",
				password: "A12345678",
				opt: &service.SignupOption{
					Nickname:  "mock-name",
					Email:     "mock@gmail.com",
					IPAddress: "127.0.0.1",
					Platform:  "web",
					Device: entity.Device{
						Model:     "Intel Mac OS X",
						Name:      "Macintosh",
						OSVersion: "13_4",
					},
				},
			},
			expected: func() *entity.Identity {
				user, _ := entity.NewUser(
					"MOCK-USER-ID",
					"Username",
					"A12345678",
					entity.WithEmail("mock@gmail.com"),
					entity.WithNickname("mock-name"),
				)
				return &entity.Identity{
					User: user,
				}
			}(),
			err: nil,
		},
		{
			name: "User Already Exist",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindUserByUsername(mock.Anything, mock.Anything).
					Return(&entity.User{
						ID:       "MOCK_ID",
						Username: "Username",
					}, nil)
				return repo
			}(),
			args: args{
				username: "Username",
				password: "A12345678",
				opt: &service.SignupOption{
					Nickname:  "mock-name",
					Email:     "mock@gmail.com",
					IPAddress: "127.0.0.1",
					Platform:  "web",
					Device: entity.Device{
						Model:     "Intel Mac OS X",
						Name:      "Macintosh",
						OSVersion: "13_4",
					},
				},
			},
			expected: nil,
			err:      errors.ErrConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			srv := service.New(tt.repo)
			actual, err := srv.Signup(
				ctx,
				tt.args.username,
				tt.args.password,
				tt.args.opt,
			)
			if tt.err != nil || err != nil {
				if !errors.Is(tt.err, err) {
					t.Errorf("Signup() error = %v, wantErr %v", err, tt.err)
				}
				return
			}

			assert.NotEmpty(t, actual.User.ID)
			assert.Equal(t, tt.expected.User.Username, actual.User.Username)
			assert.Equal(t, tt.expected.User.Email, actual.User.Email)
			assert.NotEmpty(t, actual.Session.ID)
		})
	}
}
