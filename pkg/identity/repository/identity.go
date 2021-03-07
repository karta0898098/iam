package repository

import (
	"context"

	"github.com/karta0898098/iam/pkg/identity/domain"
	"github.com/karta0898098/iam/pkg/identity/repository/model"

	"github.com/karta0898098/kara/errors"

	"gorm.io/gorm"
)

// IdentityRepository define mock repository
type IdentityRepository interface {
	// CreateProfile create a new user
	CreateProfile(ctx context.Context, params *model.Profile) (profile *domain.Profile, err error)
	// GetProfile get user profile
	GetProfile(ctx context.Context, opts *model.GetProfileOption) (profile *domain.Profile, err error)
	// UpdateProfile update user profile data by id
	UpdateProfile(ctx context.Context, id int64, opts *model.UpdateProfileOption) (err error)
}

var _ IdentityRepository = &identityRepository{}

type identityRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewIdentityRepository( /*conn *db.Connection*/) IdentityRepository {
	return &identityRepository{
		// readDB:  conn.ReadDB,
		// writeDB: conn.WriteDB,
	}
}

func (repo *identityRepository) CreateProfile(ctx context.Context, user *model.Profile) (profile *domain.Profile, err error) {
	err = repo.writeDB.
		WithContext(ctx).
		Model(user).
		Create(user).
		Error
	if err != nil {
		return nil, errors.ErrInternal.Build("reason : db occur error %v", err)
	}

	return user.ToEntity(), nil
}

func (repo *identityRepository) GetProfile(ctx context.Context, opts *model.GetProfileOption) (profile *domain.Profile, err error) {
	var (
		user model.Profile
	)
	err = repo.readDB.
		WithContext(ctx).
		Model(&user).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrResourceNotFound.Build("cant not found profile")
		}
		return nil, errors.ErrInternal.Build("reason : db occur error %v", err)
	}

	return user.ToEntity(), nil
}

func (repo *identityRepository) UpdateProfile(ctx context.Context, id int64, opts *model.UpdateProfileOption) (err error) {
	var (
		user model.Profile
	)

	err = repo.writeDB.
		WithContext(ctx).
		Model(&user).
		Where("id = ?", id).
		Updates(opts).
		Error
	if err != nil {
		return errors.ErrInternal.Build("reason : db occur error %v", err)
	}

	return nil
}
