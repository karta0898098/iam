package service

import (
	"context"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/karta0898098/kara/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	entity2 "github.com/karta0898098/iam/pkg/app/identity/entity"
	"github.com/karta0898098/iam/pkg/app/identity/mocks"
	"github.com/karta0898098/iam/pkg/app/identity/repository"
)

var (
	mockIdUtils, _ = snowflake.NewNode(1)
)

func Test_identityService_Signup(t *testing.T) {
	type fields struct {
		repo    repository.Repository
		idUtils *snowflake.Node
	}
	type args struct {
		username string
		password string
		opt      *SignupOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect *entity2.Identity
		err    error
	}{
		{
			name: "Success",
			fields: fields{
				repo: func() repository.Repository {
					repo := mocks.NewRepository(t)
					repo.
						On("FindUserByUsername", mock.Anything, mock.Anything).Return(nil, errors.ErrResourceNotFound).
						On("StoreUser", mock.Anything, mock.Anything).Return(nil).
						On("StoreSession", mock.Anything, mock.Anything).Return(nil)
					return repo
				}(),
				idUtils: mockIdUtils,
			},
			args: args{
				username: "karta0898098",
				password: "X12345678",
				opt: &SignupOption{
					Nickname:  "",
					FirstName: "Ray",
					LastName:  "Change",
					Email:     "karta0898098@gmail.com",
					IPAddress: "192.168.1.1",
					Platform:  "web",
					Device: entity2.Device{
						Model:      "Chrome",
						DeviceName: "Chrome",
						OSVersion:  "Mac Apple Silicon 12.2.1",
					},
				},
			},
			expect: &entity2.Identity{
				User: &entity2.User{
					ID:       0,
					Username: "karta0898098",
				},
				Session: nil,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			srv := &Impl{
				repo:    tt.fields.repo,
				idUtils: tt.fields.idUtils,
			}
			actual, err := srv.Signup(ctx, tt.args.username, tt.args.password, tt.args.opt)
			if err != nil || tt.err != nil {
				if !assert.True(t, errors.Is(err, tt.err)) {
					t.Errorf("the error should be %v", err)
				}
				return
			}
			assert.Equal(t, tt.expect.Username, actual.Username)
		})
	}
}
