package service

import (
	"context"
	"testing"
	"time"

	"github.com/karta0898098/iam/pkg/identity/domain"
	"github.com/karta0898098/iam/pkg/identity/repository"
	"github.com/karta0898098/iam/pkg/identity/repository/mocks"

	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/zlog"

	"github.com/bwmarrin/snowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	zlog.Setup(zlog.Config{
		Env:   "uint_test",
		AppID: "identity",
		Debug: true,
	})
}

var (
	// mockIdUtils mock snowflake
	mockIdUtils, _ = snowflake.NewNode(1)
	// mockSuccessPassword -> a12345678 sha256 result
	mockSuccessPassword = "9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a"
)

func TestIdentityServiceLogin(t *testing.T) {
	type fields struct {
		repo    func() repository.IdentityRepository
		idUtils *snowflake.Node
	}
	type args struct {
		ctx      context.Context
		account  string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		profile *domain.Profile
		err     *errors.Exception
	}{
		{
			name: "InputAccountInvalid",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345",
				password: "",
			},
			profile: nil,
			err:     errors.ErrInvalidInput,
		},
		{
			name: "InputPasswordInvalid",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "A123",
			},
			profile: nil,
			err:     errors.ErrInvalidInput,
		},
		{
			name: "GetProfileOccurError",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(nil, errors.ErrInternal)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
			},
			profile: nil,
			err:     errors.ErrInternal,
		},
		{
			name: "GetProfileResourceNotFound",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(nil, errors.ErrResourceNotFound)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
			},
			profile: nil,
			err:     errors.ErrResourceNotFound,
		},
		{
			name: "PasswordNotCorrect",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(
							&domain.Profile{
								ID:            1,
								Account:       "A12345678",
								Password:      mockSuccessPassword,
								SuspendStatus: domain.Unsuspend,
							}, nil,
						)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a123456789",
			},
			profile: nil,
			err:     errors.ErrUnauthorized,
		},
		{
			name: "AccountIsSuspend",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(
							&domain.Profile{
								ID:            1,
								Account:       "A12345678",
								Password:      mockSuccessPassword,
								SuspendStatus: domain.Suspend,
							}, nil,
						)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
			},
			profile: nil,
			err:     errors.ErrUnauthorized,
		},
		{
			name: "Success",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(
							&domain.Profile{
								ID:            1,
								Account:       "A12345678",
								Password:      mockSuccessPassword,
								SuspendStatus: domain.Unsuspend,
							}, nil,
						)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
			},
			profile: &domain.Profile{
				ID:            1,
				Account:       "A12345678",
				Password:      mockSuccessPassword,
				SuspendStatus: domain.Unsuspend,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &identityService{
				repo:    tt.fields.repo(),
				idUtils: tt.fields.idUtils,
			}
			profile, err := srv.Login(tt.args.ctx, tt.args.account, tt.args.password)
			if err != nil {
				assert.True(t, tt.err.Is(err), "error type not equal")
				return
			}
			assert.Equal(t, tt.profile, profile)
		})
	}
}

func TestIdentityServiceSignup(t *testing.T) {
	type fields struct {
		repo    func() repository.IdentityRepository
		idUtils *snowflake.Node
	}
	type args struct {
		ctx      context.Context
		account  string
		password string
		opts     *domain.SignupOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		profile *domain.Profile
		err     *errors.Exception
	}{
		{
			name: "InputAccountInvalid",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A123",
				password: "",
				opts:     nil,
			},
			profile: nil,
			err:     errors.ErrInvalidInput,
		},
		{
			name: "InputPasswordInvalid",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a1234",
				opts:     nil,
			},
			profile: nil,
			err:     errors.ErrInvalidInput,
		},
		{
			name: "InputEmailFormatInvalid",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
				opts: &domain.SignupOption{
					Nickname:  "testNickname",
					FirstName: "testFirstName",
					LastName:  "testLastName",
					Email:     "invalidEmailFormat",
					Avatar:    "",
				},
			},
			profile: nil,
			err:     errors.ErrInvalidInput,
		},
		{
			name: "AccountDuplicate",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(
							&domain.Profile{
								ID:            1,
								Account:       "A12345678",
								Password:      "",
								Nickname:      "",
								FirstName:     "",
								LastName:      "",
								Email:         "",
								Avatar:        "",
								CreatedAt:     time.Time{},
								UpdatedAt:     time.Time{},
								SuspendStatus: 0,
							}, nil,
						)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
				opts: &domain.SignupOption{
					Nickname:  "testNickname",
					FirstName: "testFirstName",
					LastName:  "testLastName",
					Email:     "test@example.com",
					Avatar:    "",
				},
			},
			profile: nil,
			err:     errors.ErrConflict,
		},
		{
			name: "CreateProfileFailed",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).Return(nil, errors.ErrResourceNotFound).
						On("CreateProfile", mock.Anything, mock.Anything).Return(nil, errors.ErrInternal)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
				opts: &domain.SignupOption{
					Nickname:  "testNickname",
					FirstName: "testFirstName",
					LastName:  "testLastName",
					Email:     "test@example.com",
					Avatar:    "",
				},
			},
			profile: nil,
			err:     errors.ErrInternal,
		},
		{
			name: "Success",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).Return(nil, errors.ErrResourceNotFound).
						On("CreateProfile", mock.Anything, mock.Anything).
						Return(&domain.Profile{
							ID:            1,
							Account:       "A12345678",
							Password:      mockSuccessPassword,
							Nickname:      "testNickname",
							FirstName:     "testFirstName",
							LastName:      "testLastName",
							Email:         "test@example.com",
							Avatar:        "",
							SuspendStatus: domain.Unsuspend,
						}, nil,
						)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:      context.Background(),
				account:  "A12345678",
				password: "a12345678",
				opts: &domain.SignupOption{
					Nickname:  "testNickname",
					FirstName: "testFirstName",
					LastName:  "testLastName",
					Email:     "test@example.com",
					Avatar:    "",
				},
			},
			profile: &domain.Profile{
				ID:            1,
				Account:       "A12345678",
				Password:      mockSuccessPassword,
				Nickname:      "testNickname",
				FirstName:     "testFirstName",
				LastName:      "testLastName",
				Email:         "test@example.com",
				Avatar:        "",
				SuspendStatus: domain.Unsuspend,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &identityService{
				repo:    tt.fields.repo(),
				idUtils: tt.fields.idUtils,
			}
			profile, err := srv.Signup(tt.args.ctx, tt.args.account, tt.args.password, tt.args.opts)
			if err != nil {
				assert.True(t, tt.err.Is(err), "error type not equal")
				return
			}
			assert.Equal(t, tt.profile, profile)
		})
	}
}

func TestIdentityServiceUpdateProfile(t *testing.T) {
	var (
		invalidNickname = "nickname too many letters nickname ........"
		normalNickname  = "nickname"
		invalidEmail    = "invalidEmail"
		normalEmail     = "test@example.com"
	)
	type fields struct {
		repo    func() repository.IdentityRepository
		idUtils *snowflake.Node
	}
	type args struct {
		ctx  context.Context
		id   int64
		opts *domain.UpdateProfileOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		err    *errors.Exception
	}{
		{
			name: "InputIDIsZero",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:  context.Background(),
				id:   0,
				opts: nil,
			},
			err: errors.ErrInvalidInput,
		},
		{
			name: "DontNeedUpdate",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx:  context.Background(),
				id:   1,
				opts: nil,
			},
			err: nil,
		},
		{
			name: "InputEmailInvalid",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
				opts: &domain.UpdateProfileOption{
					Nickname:  nil,
					FirstName: nil,
					LastName:  nil,
					Email:     &invalidEmail,
					Avatar:    nil,
				},
			},
			err: errors.ErrInvalidInput,
		},
		{
			name: "InputNicknameTooMany",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
				opts: &domain.UpdateProfileOption{
					Nickname:  &invalidNickname,
					FirstName: nil,
					LastName:  nil,
					Email:     nil,
					Avatar:    nil,
				},
			},
			err: errors.ErrInvalidInput,
		},
		{
			name: "ProfileNotFound",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(nil, errors.ErrResourceNotFound)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
				opts: &domain.UpdateProfileOption{
					Nickname:  &normalNickname,
					FirstName: nil,
					LastName:  nil,
					Email:     &normalEmail,
					Avatar:    nil,
				},
			},
			err: errors.ErrResourceNotFound,
		},
		{
			name: "UpdateProfileFailed",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(&domain.Profile{ID: 1}, nil).
						On("UpdateProfile", mock.Anything, mock.Anything, mock.Anything).
						Return(errors.ErrInternal)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
				opts: &domain.UpdateProfileOption{
					Nickname:  &normalNickname,
					FirstName: nil,
					LastName:  nil,
					Email:     &normalEmail,
					Avatar:    nil,
				},
			},
			err: errors.ErrInternal,
		},
		{
			name: "Success",
			fields: fields{
				repo: func() repository.IdentityRepository {
					repo := new(mocks.IdentityRepository)
					repo.
						On("GetProfile", mock.Anything, mock.Anything).
						Return(&domain.Profile{ID: 1}, nil).
						On("UpdateProfile", mock.Anything, mock.Anything, mock.Anything).
						Return(nil)
					return repo
				},
				idUtils: mockIdUtils,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
				opts: &domain.UpdateProfileOption{
					Nickname:  &normalNickname,
					FirstName: nil,
					LastName:  nil,
					Email:     &normalEmail,
					Avatar:    nil,
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &identityService{
				repo:    tt.fields.repo(),
				idUtils: tt.fields.idUtils,
			}
			err := srv.UpdateProfile(tt.args.ctx, tt.args.id, tt.args.opts)
			if tt.err != nil {
				assert.True(t, tt.err.Is(err))
			}
		})
	}
}

func TestIdentityServiceValidatePassword(t *testing.T) {
	type args struct {
		password string
		min      int
		max      int
	}
	tests := []struct {
		name string
		args args
		exp  bool
	}{
		{
			name: "PasswordNotEnoughLength",
			args: args{
				password: "a123456",
				min:      8,
				max:      16,
			},
			exp: false,
		},
		{
			name: "PasswordLengthTooMany",
			args: args{
				password: "a123456789100000000000000",
				min:      8,
				max:      16,
			},
			exp: false,
		},
		{
			name: "PasswordHasSpecialLetter",
			args: args{
				password: "A12345678@",
				min:      8,
				max:      16,
			},
			exp: false,
		},
		{
			name: "PasswordHasSpaceCharacter",
			args: args{
				password: "A12345678 ",
				min:      8,
				max:      16,
			},
			exp: false,
		},
		{
			name: "Success",
			args: args{
				password: "A12345678",
				min:      8,
				max:      16,
			},
			exp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &identityService{}
			got := srv.ValidatePassword(tt.args.password, tt.args.min, tt.args.max)
			assert.Equal(t, tt.exp, got)
		})
	}
}

func TestIdentityServiceValidateEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		exp  bool
	}{
		{
			name: "InvalidEmail",
			args: args{
				email: "InvalidEmail",
			},
			exp: false,
		},
		{
			name: "InvalidEmail@",
			args: args{
				email: "InvalidEmail@",
			},
			exp: false,
		},
		{
			name: "Success",
			args: args{
				email: "test@example.com",
			},
			exp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &identityService{}
			got := srv.ValidateEmail(tt.args.email)
			assert.Equal(t, tt.exp, got)
		})
	}
}

func Test_identityService_ToSHA256(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		exp  string
	}{
		{
			name: "Success",
			args: args{
				password: "a12345678",
			},
			exp: mockSuccessPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &identityService{}
			password := srv.ToSHA256(tt.args.password)
			assert.Equal(t, tt.exp, password)
		})
	}
}
