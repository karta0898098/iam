package model

import (
	"time"

	"github.com/karta0898098/iam/pkg/identity/domain"
)

// Profile define user information
type Profile struct {
	ID            int64                // ID unique mock number
	Account       string               // Account user login mock account
	Password      string               // Password user login password
	Nickname      string               // Nickname user nickname
	FirstName     string               // FirstName user first name
	LastName      string               // LastName user last name
	Email         string               // Email user email address
	Avatar        string               // Avatar means user profile picture URL
	CreatedAt     time.Time            // CreatedAt this account create time
	UpdatedAt     time.Time            // UpdatedAt this account update time
	SuspendStatus domain.SuspendStatus // SuspendStatus this account is suspend
}

func (user *Profile) ToEntity() *domain.Profile {
	return &domain.Profile{
		ID:            user.ID,
		Account:       user.Account,
		Password:      user.Password,
		Nickname:      user.Nickname,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Avatar:        user.Avatar,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		SuspendStatus: user.SuspendStatus,
	}
}

type GetProfileOption struct {
	ID        int64     // ID unique mock number
	Account   string    // Account user login mock account
	Nickname  string    // Nickname user nickname
	Email     string    // Email user email address
	CreatedAt time.Time // CreatedAt this account create time
	UpdatedAt time.Time // UpdatedAt this account update time
}

type UpdateProfileOption struct {
	Password  *string   // Password user login password
	Nickname  *string   // Nickname user nickname
	FirstName *string   // FirstName user first name
	LastName  *string   // LastName user last name
	Email     *string   // Email user email address
	Avatar    *string   // Avatar means user profile picture URL
	UpdatedAt time.Time // UpdatedAt this account update time
}
