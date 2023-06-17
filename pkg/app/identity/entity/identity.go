package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Identity aggregate user and session
type Identity struct {
	*User
	*Session
}

// NewAccessToken new access token
func (i *Identity) NewAccessToken() string {
	return newSignedToken(
		i.Session.ID,
		i.User.ID,
		i.Session.CreateAt,
		i.Session.ExpireAt,
	)
}

// NewRefreshToken new refresh token
func (i *Identity) NewRefreshToken() string {
	expiresAt := i.Session.CreateAt.Add(time.Second * 60 * 60 * 24 * 30)
	return newSignedToken(
		i.Session.ID,
		i.User.ID,
		i.Session.CreateAt,
		expiresAt,
	)
}

// newSignedToken signing token string
func newSignedToken(sessionID, userID string, now time.Time, expiresAt time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        sessionID,
	})

	accessToken, _ := token.SignedString([]byte(SigninKey))

	return accessToken
}
