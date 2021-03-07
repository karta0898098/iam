package domain

import (
	"context"
	"time"
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

// Profile define user information
type Profile struct {
	ID            int64         // ID unique mock number
	Account       string        // Account user login mock account
	Password      string        // Password user login password
	Nickname      string        // Nickname user nickname
	FirstName     string        // FirstName user first name
	LastName      string        // LastName user last name
	Email         string        // Email user email address
	Avatar        string        // Avatar means user profile picture URL
	CreatedAt     time.Time     // CreatedAt this account create time
	UpdatedAt     time.Time     // UpdatedAt this account update time
	SuspendStatus SuspendStatus // SuspendStatus this account is suspend
}

// IdentityService define mock service
type IdentityService interface {
	// Login verify user account
	Login(ctx context.Context, account string, password string) (profile *Profile, err error)
	// Signup sign up a new user
	Signup(ctx context.Context, account string, password string, opts *SignupOption) (profile *Profile, err error)
	// UpdateProfile update user profile
	UpdateProfile(ctx context.Context, id int64, opts *UpdateProfileOption) (err error)
}
