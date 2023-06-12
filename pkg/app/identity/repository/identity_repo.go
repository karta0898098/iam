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
	ID        string                   // ID unique mock number
	Username  string                   // Username user login identity account
	Password  string                   // Password user login identity password
	Nickname  string                   // Nickname user nickname
	FirstName string                   // FirstName user first name
	LastName  string                   // LastName user last name
	Email     string                   // Email user email address
	Avatar    string                   // Avatar means user profile picture URL
	CreatedAt time.Time                // CreatedAt this account create time
	UpdatedAt time.Time                // UpdatedAt this account update time
	Status    entity.UserAccountStatus // Status this account is suspend
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
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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
		CreatedAt: dao.CreatedAt,
		UpdatedAt: dao.UpdatedAt,
		Status:    dao.Status,
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
		Model(&user).
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

	err = repo.writeDB.
		WithContext(ctx).
		Model(session).
		Create(session).
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
