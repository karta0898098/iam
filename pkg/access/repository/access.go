package repository

import (
	"context"
	"time"

	"github.com/karta0898098/iam/pkg/access/domain"
	"github.com/karta0898098/iam/pkg/access/repository/model"

	"github.com/karta0898098/kara/db/rw/db"
	"github.com/karta0898098/kara/errors"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var _ AccessRepository = &accessRepository{}

type AccessRepository interface {
	CreateUserAccess(ctx context.Context, userID int64, role []domain.Role) (access []*model.Access, err error)
	ListUserAccess(ctx context.Context, userID int64) (access []*model.Access, err error)
	StoreToken(ctx context.Context, tokenType domain.TokenType, tokenID string, value []byte, ttl time.Duration) (err error)
}

func NewAccessRepository(redisClient *redis.Client, conn db.Connection) AccessRepository {
	return &accessRepository{
		redisClient: redisClient,
		readDB:      conn.ReadDB,
		writeDB:     conn.WriteDB,
	}
}

type accessRepository struct {
	redisClient *redis.Client
	readDB      *gorm.DB
	writeDB     *gorm.DB
}

func (repo *accessRepository) CreateUserAccess(ctx context.Context, userID int64, roles []domain.Role) (access []*model.Access, err error) {
	for i := 0; i < len(roles); i++ {
		access = append(access, &model.Access{
			UserID:    userID,
			Role:      roles[i],
			CreatedAt: time.Now().UTC(),
		})
	}

	err = repo.writeDB.
		WithContext(ctx).
		Model(&model.Access{}).
		Create(&access).
		Error

	if err != nil {
		return nil, errors.ErrInternal.BuildWithError(err)
	}

	return access, nil
}

func (repo *accessRepository) ListUserAccess(ctx context.Context, userID int64) (access []*model.Access, err error) {
	err = repo.readDB.
		WithContext(ctx).
		Model(&model.Access{}).
		Where("user_id = ?", userID).
		Find(&access).
		Error

	if err != nil {
		return nil, errors.ErrInternal.BuildWithError(err)
	}

	return access, nil
}

func (repo *accessRepository) StoreToken(ctx context.Context, tokenType domain.TokenType, tokenID string, value []byte, ttl time.Duration) (err error) {
	var (
		key string
	)

	key = tokenType.GetPrefix() + tokenID

	err = repo.redisClient.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return errors.ErrInternal.BuildWithError(err)
	}

	return nil
}
