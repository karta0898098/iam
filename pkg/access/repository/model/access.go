package model

import (
	"time"

	"github.com/karta0898098/iam/pkg/access/domain"
)

type Access struct {
	UserID    int64
	Roles     domain.Role
	CreatedAt time.Time
}
