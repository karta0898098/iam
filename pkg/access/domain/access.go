package domain

import (
	"context"

	identity "github.com/karta0898098/iam/pkg/identity/domain"
)

type SuspendStatus int8

const (
	UnknownStatus = iota
	Unsuspend
	Suspend
)

func (s SuspendStatus) ToString() string {
	switch s {
	case Unsuspend:
		return "unsuspend"
	case Suspend:
		return "suspend"
	default:
		return "unknown"
	}
}

type Role int8

const (
	UnknownRole = iota
	SuperAdmin
	Admin
	User
)

func (r Role) ToString() string {
	switch r {
	case SuperAdmin:
		return "super admin"
	case Admin:
		return "admin"
	case User:
		return "user"
	default:
		return "unknown"
	}
}

type TokenType int8

const (
	UnknownTokenType = iota
	AccessToken
	RefreshToken
)

func (t TokenType) ToString() string {
	switch t {
	case AccessToken:
		return "accessToken"
	case RefreshToken:
		return "refreshToken"
	default:
		return "unknown"
	}
}

func (t TokenType) GetPrefix() string {
	switch t {
	case AccessToken:
		return "access:"
	case RefreshToken:
		return "refresh:"
	default:
		return "unknown:"
	}
}

type Claims struct {
	UserID        int64         // UserID  unique identity number
	Account       string        // Account user login mock account
	Nickname      string        // Nickname user nickname
	FirstName     string        // FirstName user first name
	LastName      string        // LastName user last name
	Email         string        // Email user email address
	Avatar        string        // Avatar means user profile picture URL
	SuspendStatus SuspendStatus // SuspendStatus this account is suspend
	Roles         []Role        // Roles user has role
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpireIn     int64
}

// AccessService define access service manager token and check user role
type AccessService interface {
	// CreateAccessTokens create user login or signup token
	CreateAccessTokens(ctx context.Context, user *identity.Profile) (tokens *Tokens, err error)
}
