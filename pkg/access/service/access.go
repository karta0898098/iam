package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/karta0898098/iam/pkg/access/domain"
	"github.com/karta0898098/iam/pkg/access/repository"
	identity "github.com/karta0898098/iam/pkg/identity/domain"

	"github.com/karta0898098/kara/errors"

	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
)

var _ domain.AccessService = &accessService{}

const (
	jwtTokenSignature    = "w6ADHx4Mbb8ww3KA"
	jwtAccessTokenExpire = 7 * 24 * time.Hour

	refreshTokenExpire = 14 * 24 * time.Hour
)

type accessService struct {
	idUtils *snowflake.Node
	repo    repository.AccessRepository
}

func NewAccessService(idUtils *snowflake.Node, repo repository.AccessRepository) domain.AccessService {
	return &accessService{
		idUtils: idUtils,
		repo:    repo,
	}
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

	// create access token
	accessTokenExpireTime := time.Now().UTC().Add(jwtAccessTokenExpire).Unix()
	accessTokenID := srv.idUtils.Generate().Base58()
	accessClaimsValue, _ := json.Marshal(&claims)
	accessToken := srv.CreateJwtToken(accessTokenID, accessTokenExpireTime)
	err = srv.repo.StoreToken(ctx, domain.AccessToken, accessTokenID, accessClaimsValue, jwtAccessTokenExpire)
	if err != nil {
		return nil, err
	}

	// create refresh token
	refreshToken := srv.idUtils.Generate().Base64()
	refreshTokenClaims := []byte(accessToken)
	err = srv.repo.StoreToken(ctx, domain.RefreshToken, refreshToken, refreshTokenClaims, refreshTokenExpire)
	if err != nil {
		return nil, err
	}

	tokens = &domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireIn:     accessTokenExpireTime,
	}

	return tokens, nil
}

func (srv *accessService) CreateJwtToken(tokenID string, exp int64) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["jti"] = tokenID
	claims["exp"] = exp

	// Generate encoded token and send it as response.
	// why not catch error
	// because impossible occur error
	// except SignedString input type not []byte
	tokenString, _ := token.SignedString([]byte(jwtTokenSignature))

	return tokenString
}
