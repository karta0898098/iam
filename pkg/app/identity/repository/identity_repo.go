package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/karta0898098/iam/pkg/app/identity/entity"
	"github.com/karta0898098/iam/pkg/db"
	"github.com/karta0898098/iam/pkg/errors"
)

// UserDAO define user information
type UserDAO struct {
	ID        string                   `gorm:"column:id"`         // ID unique mock number
	Username  string                   `gorm:"column:username"`   // Username user login identity account
	Password  string                   `gorm:"column:password"`   // Password user login identity password
	Nickname  string                   `gorm:"column:nickname"`   // Nickname user nickname
	FirstName string                   `gorm:"column:first_name"` // FirstName user first name
	LastName  string                   `gorm:"column:last_name"`  // LastName user last name
	Email     string                   `gorm:"column:email"`      // Email user email address
	Avatar    string                   `gorm:"column:avatar"`     // Avatar means user profile picture URL
	CreatedAt int64                    `gorm:"column:created_at"` // CreatedAt this account create time
	UpdatedAt int64                    `gorm:"column:updated_at"` // UpdatedAt this account update time
	Status    entity.UserAccountStatus `gorm:"column:status"`     // Status this account is suspend
}

func (u UserDAO) TableName() string {
	return "users"
}

func UnmarshalUserDAO(user *entity.User) *UserDAO {
	return &UserDAO{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.UnixMilli(),
		UpdatedAt: user.UpdatedAt.UnixMilli(),
		Status:    user.Status,
	}
}

func UnmarshalUser(dao *UserDAO) *entity.User {
	return &entity.User{
		ID:        dao.ID,
		Username:  dao.Username,
		Password:  dao.Password,
		Nickname:  dao.Nickname,
		FirstName: dao.FirstName,
		LastName:  dao.LastName,
		Email:     dao.Email,
		Avatar:    dao.Avatar,
		CreatedAt: time.UnixMilli(dao.CreatedAt),
		UpdatedAt: time.UnixMilli(dao.UpdatedAt),
		Status:    dao.Status,
	}
}

type SessionDAO struct {
	ID              string `grom:"column:id"`
	UserID          string `grom:"column:user_id"`
	CreatedAt       int64  `grom:"column:created_at"`
	UpdatedAt       int64  `grom:"column:updated_at"`
	ExpireAt        int64  `grom:"column:expire_at"`
	IPAddress       string `grom:"column:ip_address"`
	IdpProvider     string `grom:"column:idp_provider"`
	Platform        string `grom:"column:platform"`
	DeviceModel     string `grom:"column:device_model"`
	DeviceName      string `grom:"column:device_name"`
	DeviceOSVersion string `grom:"column:device_os_version"`
}

func (s SessionDAO) TableName() string {
	return "sessions"
}

func UnmarshalSessionDAO(session *entity.Session) *SessionDAO {
	return &SessionDAO{
		ID:              session.ID,
		UserID:          session.UserID,
		CreatedAt:       session.CreateAt.UnixMilli(),
		UpdatedAt:       session.UpdateAt.UnixMilli(),
		ExpireAt:        session.ExpireAt.UnixMilli(),
		IPAddress:       session.IPAddress,
		IdpProvider:     session.IdpProvider,
		Platform:        session.Platform,
		DeviceModel:     session.Device.Model,
		DeviceName:      session.Device.Name,
		DeviceOSVersion: session.Device.OSVersion,
	}
}

var _ Repository = &IdentityRepository{}

type Repository interface {
	StoreUser(ctx context.Context, user *entity.User) (err error)
	FindUserByUsername(ctx context.Context, username string) (profile *entity.User, err error)
	StoreSession(ctx context.Context, session *entity.Session) (err error)
}

type IdentityRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func New(conn db.Connection) Repository {
	return &IdentityRepository{
		readDB:  conn.ReadDB(),
		writeDB: conn.WriteDB(),
	}
}

func (repo *IdentityRepository) FindUserByUsername(ctx context.Context, username string) (profile *entity.User, err error) {
	var (
		user UserDAO
	)

	if username == "" {
		return nil, errors.Wrapf(errors.ErrInternal, "repo: input username is empty")
	}

	err = repo.readDB.
		WithContext(ctx).
		Model(user).
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(errors.ErrResourceNotFound, "cant not found profile username=%v", username)
		}
		return nil, errors.Wrapf(errors.ErrInternal, "reason : db occur error %v", err)
	}

	return UnmarshalUser(&user), nil
}

func (repo *IdentityRepository) StoreSession(ctx context.Context, session *entity.Session) (err error) {
	dao := UnmarshalSessionDAO(session)

	err = repo.writeDB.
		WithContext(ctx).
		Model(dao).
		Create(dao).
		Error
	if err != nil {
		return errors.Wrapf(errors.ErrInternal, "failed to create session %#v, err %v", session, err)
	}

	return nil
}

func (repo *IdentityRepository) StoreUser(ctx context.Context, user *entity.User) (err error) {
	dao := UnmarshalUserDAO(user)

	err = repo.writeDB.
		WithContext(ctx).
		Model(dao).
		Create(dao).
		Error
	if err != nil {
		return errors.Wrapf(errors.ErrInternal, "reason failed to create user %#v, err %v", user, err)
	}

	return nil
}
