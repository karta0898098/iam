package service

import (
	"context"
	"testing"
	"time"

	"github.com/karta0898098/iam/pkg/access/domain"
	"github.com/karta0898098/iam/pkg/access/repository"
	"github.com/karta0898098/iam/pkg/access/repository/mocks"
	"github.com/karta0898098/iam/pkg/access/repository/model"
	identity "github.com/karta0898098/iam/pkg/identity/domain"

	"github.com/karta0898098/kara/errors"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockIDGenerator struct {
}

func (m *mockIDGenerator) Generate() snowflake.ID {
	return snowflake.ParseInt64(1)
}

func TestAccessServiceCreateAccessTokens(t *testing.T) {
	var (
		mockTokenGenerator = func(tokenID string, exp int64) string {
			return "testAccessToken"
		}

		mockIDGenerator = &mockIDGenerator{}

		mockProfile = &identity.Profile{ID: 1}
	)

	type fields struct {
		idUtils        idUtilsAdapter
		repo           func() repository.AccessRepository
		tokenGenerator func(tokenID string, exp int64) (token string)
	}
	type args struct {
		ctx  context.Context
		user *identity.Profile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		tokens *domain.Tokens
		err    *errors.Exception
	}{
		{
			name: "ProfileIsNil",
			fields: fields{
				idUtils: mockIDGenerator,
				repo: func() repository.AccessRepository {
					return nil
				},
				tokenGenerator: mockTokenGenerator,
			},
			args: args{
				ctx:  context.Background(),
				user: nil,
			},
			tokens: nil,
			err:    errors.ErrInvalidInput,
		},
		{
			name: "ListUserAccessOccurError",
			fields: fields{
				idUtils: mockIDGenerator,
				repo: func() repository.AccessRepository {
					repo := new(mocks.AccessRepository)
					repo.
						On("ListUserAccess", mock.Anything, mock.Anything).
						Return(nil, errors.ErrInternal)
					return repo
				},
				tokenGenerator: mockTokenGenerator,
			},
			args: args{
				ctx:  context.Background(),
				user: mockProfile,
			},
			tokens: nil,
			err:    errors.ErrInternal,
		},
		{
			name: "CreateUserAccessOccurError",
			fields: fields{
				idUtils: mockIDGenerator,
				repo: func() repository.AccessRepository {
					repo := new(mocks.AccessRepository)
					repo.
						On("ListUserAccess", mock.Anything, mock.Anything).
						Return([]*model.Access{}, nil)
					repo.
						On("CreateUserAccess", mock.Anything, mock.Anything, mock.Anything).
						Return(nil, errors.ErrInternal)
					return repo
				},
				tokenGenerator: mockTokenGenerator,
			},
			args: args{
				ctx:  context.Background(),
				user: mockProfile,
			},
			tokens: nil,
			err:    errors.ErrInternal,
		},
		{
			name: "Success",
			fields: fields{
				idUtils: mockIDGenerator,
				repo: func() repository.AccessRepository {
					repo := new(mocks.AccessRepository)
					repo.
						On("ListUserAccess", mock.Anything, mock.Anything).
						Return([]*model.Access{{
							UserID:    1,
							Role:      domain.User,
							CreatedAt: time.Now(),
						}}, nil)
					repo.
						On("StoreToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(nil)
					return repo
				},
				tokenGenerator: mockTokenGenerator,
			},
			args: args{
				ctx:  context.Background(),
				user: mockProfile,
			},
			tokens: &domain.Tokens{
				AccessToken:  "testAccessToken",
				RefreshToken: "MQ==",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &accessService{
				idUtils:        tt.fields.idUtils,
				repo:           tt.fields.repo(),
				tokenGenerator: tt.fields.tokenGenerator,
			}
			tokens, err := srv.CreateAccessTokens(tt.args.ctx, tt.args.user)
			if err != nil {
				assert.True(t, tt.err.Is(err), "error type not equal")
				return
			}
			assert.Equal(t, tt.tokens.AccessToken, tokens.AccessToken)
			assert.Equal(t, tt.tokens.RefreshToken, tokens.RefreshToken)
		})
	}
}

func TestAccessServiceCreateJwtToken(t *testing.T) {
	type args struct {
		tokenID string
		exp     int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				tokenID: "1",
				exp:     time.Now().UTC().Unix(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &accessService{}
			token := srv.CreateJwtToken(tt.args.tokenID, tt.args.exp)
			log.Debug().Msg(token)
		})
	}
}
