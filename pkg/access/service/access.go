package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/karta0898098/iam/pkg/access/domain"
	"github.com/karta0898098/iam/pkg/access/repository"
	identity "github.com/karta0898098/iam/pkg/identity/domain"

	"github.com/karta0898098/kara/errors"
)

var _ domain.AccessService = &accessService{}

const (
	jwtTokenSignature    = "w6ADHx4Mbb8ww3KA"
	jwtAccessTokenExpire = 7 * 24 * time.Hour
)

type accessService struct {
	idUtils *snowflake.Node
	repo    repository.AccessRepository
}

func NewAccessService() domain.AccessService {
	return &accessService{

	}
}

func (srv *accessService) AssignUserRole(ctx context.Context, userID int64, role domain.Role) (err error) {
	panic("implement me")
}

func (srv *accessService) CreateAccessTokens(ctx context.Context, user *identity.Profile) (tokens *domain.Tokens, err error) {
	var (
		defaultUserRoles = []domain.Role{domain.User}
		claims           *domain.Claims
	)
	if user == nil {
		return nil, errors.ErrInvalidInput.Build("input user is nil")
	}

	access, err := srv.repo.ListUserAccess(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if len(access) == 0 {
		access, err = srv.repo.CreateUserAccess(ctx, user.ID, defaultUserRoles)
		if err != nil {
			return nil, err
		}
	}

	claims = &domain.Claims{
		UserID:        user.ID,
		Account:       user.Account,
		Nickname:      user.Nickname,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Avatar:        user.Avatar,
		SuspendStatus: domain.SuspendStatus(user.SuspendStatus),
		Roles:         nil,
	}

	for i := 0; i < len(access); i++ {
		claims.Roles = append(claims.Roles, access[i].Roles)
	}

	accessTokenID := srv.idUtils.Generate().Base58()
	accessClaimsValue, _ := json.Marshal(&claims)
	err = srv.repo.StoreToken(ctx, domain.AccessToken, accessTokenID, accessClaimsValue, jwtAccessTokenExpire)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
