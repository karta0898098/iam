package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/karta0898098/iam/pkg/identity/domain"
	"github.com/karta0898098/iam/pkg/identity/repository/mocks"
	"github.com/karta0898098/iam/pkg/identity/repository/model"

	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/zlog"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mockNewLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
		Colorful:      true,
	},
)

func TestIdentityRepositoryCreateProfile(t *testing.T) {
	zlog.Setup(zlog.Config{
		Debug: true,
	})

	var (
		// "2006/01/02 15:04:05"
		mockCreatedAtTime, _ = time.Parse(time.RFC3339, "2020-03-15T17:55:20+08:00")
		mockUpdatedAtTime, _ = time.Parse(time.RFC3339, "2020-03-15T17:55:20+08:00")
	)
	type fields struct {
		writeDB func(p *model.Profile) *gorm.DB
	}
	type args struct {
		ctx  context.Context
		user *model.Profile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		profile *domain.Profile
		err     *errors.Exception
	}{
		{
			name: "Success",
			fields: fields{
				writeDB: func(p *model.Profile) *gorm.DB {
					const insert = "INSERT INTO `profiles` (`account`,`password`,`nickname`,`first_name`,`last_name`,`email`,`avatar`,`created_at`,`updated_at`,`suspend_status`,`id`) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
					conn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.
						ExpectExec(insert).
						WithArgs(p.Account, p.Password, p.Nickname, p.FirstName, p.LastName, p.Email, p.Avatar, mocks.AnyTime{}, mocks.AnyTime{}, p.SuspendStatus, p.ID).
						WillReturnResult(sqlmock.NewResult(1, 1))

					db, _ := gorm.Open(mysql.New(mysql.Config{
						Conn:                      conn,
						SkipInitializeWithVersion: true,
					}), &gorm.Config{
						Logger: mockNewLogger,
					})
					db.Debug()

					return db
				},
			},
			args: args{
				ctx: context.Background(),
				user: &model.Profile{
					ID:            1,
					Account:       "A12345678",
					Password:      "9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a",
					Nickname:      "test",
					FirstName:     "test",
					LastName:      "test",
					Email:         "test@example.com",
					Avatar:        "",
					CreatedAt:     mockCreatedAtTime,
					UpdatedAt:     mockUpdatedAtTime,
					SuspendStatus: 1,
				},
			},
			profile: &domain.Profile{
				ID:            1,
				Account:       "A12345678",
				Password:      "9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a",
				Nickname:      "test",
				FirstName:     "test",
				LastName:      "test",
				Email:         "test@example.com",
				Avatar:        "",
				CreatedAt:     mockCreatedAtTime,
				UpdatedAt:     mockUpdatedAtTime,
				SuspendStatus: 1,
			},
			err: nil,
		},
		{
			name: "Failed",
			fields: fields{
				writeDB: func(p *model.Profile) *gorm.DB {
					const insert = "INSERT INTO `profiles` (`account`,`password`,`nickname`,`first_name`,`last_name`,`email`,`avatar`,`created_at`,`updated_at`,`suspend_status`,`id`) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
					conn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.
						ExpectExec(insert).
						WithArgs(p.Account, p.Password, p.Nickname, p.FirstName, p.LastName, p.Email, p.Avatar, mocks.AnyTime{}, mocks.AnyTime{}, p.SuspendStatus, p.ID).
						WillReturnError(sql.ErrConnDone)

					db, _ := gorm.Open(mysql.New(mysql.Config{
						Conn:                      conn,
						SkipInitializeWithVersion: true,
					}), &gorm.Config{
						Logger: mockNewLogger,
					})
					db.Debug()
					return db
				},
			},
			args: args{
				ctx: context.Background(),
				user: &model.Profile{
					ID:            1,
					Account:       "A12345678",
					Password:      "9fcefc0080d894e83ca7d360ce5ccd9ead2c5d8a80a10f9fa9698510aaba865a",
					Nickname:      "test",
					FirstName:     "test",
					LastName:      "test",
					Email:         "test@example.com",
					Avatar:        "",
					CreatedAt:     mockCreatedAtTime,
					UpdatedAt:     mockUpdatedAtTime,
					SuspendStatus: 1,
				},
			},
			profile: nil,
			err:     errors.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &identityRepository{
				writeDB: tt.fields.writeDB(tt.args.user),
			}
			profile, err := repo.CreateProfile(tt.args.ctx, tt.args.user)
			if err != nil {
				assert.True(t, tt.err.Is(err), err)
				return
			}
			assert.Equal(t, tt.profile, profile)
		})
	}
}

func TestIdentityRepositoryGetProfile(t *testing.T) {
	type fields struct {
		readDB func(*model.GetProfileOption, *domain.Profile) *gorm.DB
	}
	type args struct {
		ctx  context.Context
		opts *model.GetProfileOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		profile *domain.Profile
		err     *errors.Exception
	}{
		{
			name: "OptionIsNil",
			fields: fields{
				readDB: func(option *model.GetProfileOption, profile *domain.Profile) *gorm.DB {
					return nil
				},
			},
			args: args{
				ctx:  context.Background(),
				opts: nil,
			},
			profile: nil,
			err:     errors.ErrInternal,
		},
		{
			name: "ResourceNotFound",
			fields: fields{
				readDB: func(opt *model.GetProfileOption, result *domain.Profile) *gorm.DB {
					const querySQL = "SELECT * FROM `profiles` WHERE id = ? ORDER BY `profiles`.`id` LIMIT 1"
					conn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.
						ExpectQuery(querySQL).
						WithArgs(opt.ID).
						WillReturnRows(sqlmock.NewRows([]string{"id", "account", "password", "nickname", "first_name", "last_name", "email", "avatar", "suspend_status"}))

					db, _ := gorm.Open(mysql.New(mysql.Config{
						Conn:                      conn,
						SkipInitializeWithVersion: true,
					}), &gorm.Config{
						Logger: mockNewLogger,
					})
					db.Debug()
					return db
				},
			},
			args: args{
				ctx:  context.Background(),
				opts: &model.GetProfileOption{ID: 1},
			},
			profile: nil,
			err:     errors.ErrResourceNotFound,
		},
		{
			name: "ErrInternal",
			fields: fields{
				readDB: func(opt *model.GetProfileOption, result *domain.Profile) *gorm.DB {
					const querySQL = "SELECT * FROM `profiles` WHERE id = ? ORDER BY `profiles`.`id` LIMIT 1"
					conn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.
						ExpectQuery(querySQL).
						WithArgs(opt.ID).
						WillReturnError(sql.ErrConnDone)

					db, _ := gorm.Open(mysql.New(mysql.Config{
						Conn:                      conn,
						SkipInitializeWithVersion: true,
					}), &gorm.Config{
						Logger: mockNewLogger,
					})
					db.Debug()
					return db
				},
			},
			args: args{
				ctx:  context.Background(),
				opts: &model.GetProfileOption{ID: 1},
			},
			profile: nil,
			err:     errors.ErrInternal,
		},
		{
			name: "SuccessGetProfileByID",
			fields: fields{
				readDB: func(opt *model.GetProfileOption, result *domain.Profile) *gorm.DB {
					const querySQL = "SELECT * FROM `profiles` WHERE id = ? ORDER BY `profiles`.`id` LIMIT 1"
					conn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.
						ExpectQuery(querySQL).
						WithArgs(opt.ID).
						WillReturnRows(
							sqlmock.
								NewRows([]string{"id", "account", "password", "nickname", "first_name", "last_name", "email", "avatar", "suspend_status"}).
								AddRow(result.ID, result.Account, result.Password, result.Nickname, result.FirstName, result.LastName, result.Email, result.Avatar, result.SuspendStatus),
						)

					db, _ := gorm.Open(mysql.New(mysql.Config{
						Conn:                      conn,
						SkipInitializeWithVersion: true,
					}), &gorm.Config{
						Logger: mockNewLogger,
					})
					db.Debug()
					return db
				},
			},
			args: args{
				ctx: context.Background(),
				opts: &model.GetProfileOption{
					ID: 1,
				},
			},
			profile: &domain.Profile{
				ID:            1,
				Account:       "test",
				Password:      "test",
				Nickname:      "test",
				FirstName:     "test",
				LastName:      "test",
				Email:         "test@example.com",
				Avatar:        "",
				SuspendStatus: 1,
			},
			err: nil,
		},
		{
			name: "SuccessGetProfileByIDAndAccount",
			fields: fields{
				readDB: func(opt *model.GetProfileOption, result *domain.Profile) *gorm.DB {
					const querySQL = "SELECT * FROM `profiles` WHERE id = ? AND account = ? ORDER BY `profiles`.`id` LIMIT 1"
					conn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.
						ExpectQuery(querySQL).
						WithArgs(opt.ID, opt.Account).
						WillReturnRows(
							sqlmock.
								NewRows([]string{"id", "account", "password", "nickname", "first_name", "last_name", "email", "avatar", "suspend_status"}).
								AddRow(result.ID, result.Account, result.Password, result.Nickname, result.FirstName, result.LastName, result.Email, result.Avatar, result.SuspendStatus),
						)

					db, _ := gorm.Open(mysql.New(mysql.Config{
						Conn:                      conn,
						SkipInitializeWithVersion: true,
					}), &gorm.Config{
						Logger: mockNewLogger,
					})
					db.Debug()
					return db
				},
			},
			args: args{
				ctx: context.Background(),
				opts: &model.GetProfileOption{
					ID:      1,
					Account: "test",
				},
			},
			profile: &domain.Profile{
				ID:            1,
				Account:       "test",
				Password:      "test",
				Nickname:      "test",
				FirstName:     "test",
				LastName:      "test",
				Email:         "test@example.com",
				Avatar:        "",
				SuspendStatus: 1,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &identityRepository{
				readDB: tt.fields.readDB(tt.args.opts, tt.profile),
			}
			profile, err := repo.GetProfile(tt.args.ctx, tt.args.opts)
			if err != nil {
				assert.True(t, tt.err.Is(err), err)
				return
			}
			assert.Equal(t, tt.profile, profile)
		})
	}
}

func TestIdentityRepositoryUpdateProfile(t *testing.T) {
	var (
		mockNickname         = "test"
		mockUpdatedAtTime, _ = time.Parse(time.RFC3339, "2020-03-15T17:55:20+08:00")
	)

	type fields struct {
		writeDB func(id int64, opts *model.UpdateProfileOption) *gorm.DB
	}
	type args struct {
		ctx  context.Context
		id   int64
		opts *model.UpdateProfileOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		err    *errors.Exception
	}{
		{
			name: "Success",
			fields: fields{
				writeDB: func(id int64, opts *model.UpdateProfileOption) *gorm.DB {
					const updateSQL = "UPDATE `profiles` SET `nickname`=?,`updated_at`=? WHERE id = ?"
					conn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.
						ExpectExec(updateSQL).
						WithArgs(*opts.Nickname, mocks.AnyTime{}, id).
						WillReturnResult(sqlmock.NewResult(1, 1))

					db, _ := gorm.Open(mysql.New(mysql.Config{
						Conn:                      conn,
						SkipInitializeWithVersion: true,
					}), &gorm.Config{
						Logger: mockNewLogger,
					})
					return db
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
				opts: &model.UpdateProfileOption{
					Nickname:  &mockNickname,
					UpdatedAt: mockUpdatedAtTime,
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &identityRepository{
				writeDB: tt.fields.writeDB(tt.args.id, tt.args.opts),
			}
			err := repo.UpdateProfile(tt.args.ctx, tt.args.id, tt.args.opts)
			if tt.err != nil {
				assert.True(t, tt.err.Is(err), err)
			}
		})
	}
}
