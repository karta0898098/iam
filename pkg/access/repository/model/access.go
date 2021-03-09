package model

import (
	"time"

	"github.com/karta0898098/iam/pkg/access/domain"
)

type Access struct {
	UserID    int64
	Role      domain.Role
	CreatedAt time.Time
}

func (a *Access) TableName() string {
	return "access"
}
