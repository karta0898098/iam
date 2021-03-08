package repository

import (
	"context"
	"time"

	"github.com/karta0898098/iam/pkg/access/domain"
	"github.com/karta0898098/iam/pkg/access/repository/model"
)

type AccessRepository interface {
	CreateUserAccess(ctx context.Context, userID int64, role []domain.Role) (access []*model.Access, err error)
	ListUserAccess(ctx context.Context, userID int64) (access []*model.Access, err error)
	StoreToken(ctx context.Context, tokenType domain.TokenType, tokenID string, value []byte, ttl time.Duration) (err error)
}
