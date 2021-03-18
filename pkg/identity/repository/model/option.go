package model

import (
	"time"

	"github.com/karta0898098/iam/pkg/identity/domain"
	"gorm.io/gorm"
)

// Profile define user information
type Profile struct {
	ID            int64                // ID unique mock number
	Account       string               // Account user login identity account
	Password      string               // Password user login identity password
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
	ID       int64  // ID unique mock number
	Account  string // Account user login mock account
	Nickname string // Nickname user nickname
	Email    string // Email user email address
}

func (opt *GetProfileOption) Query(db *gorm.DB) *gorm.DB {
	if opt.ID != 0 {
		db = db.Where("id = ?", opt.ID)
	}

	if opt.Account != "" {
		db = db.Where("account = ?", opt.Account)
	}

	if opt.Nickname != "" {
		db = db.Where("nickname = ?", opt.Nickname)
	}

	if opt.Email != "" {
		db = db.Where("email = ?", opt.Email)
	}

	return db
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

func (opt *UpdateProfileOption) ToMap() (updateField map[string]interface{}) {
	updateField = make(map[string]interface{})
	updateField["updated_at"] = opt.UpdatedAt

	if opt.Password != nil {
		updateField["password"] = *opt.Password
	}

	if opt.Nickname != nil {
		updateField["nickname"] = *opt.Nickname
	}

	if opt.FirstName != nil {
		updateField["first_name"] = *opt.FirstName
	}

	if opt.LastName != nil {
		updateField["last_name"] = *opt.LastName
	}

	if opt.Email != nil {
		updateField["email"] = *opt.Email
	}

	if opt.Avatar != nil {
		updateField["avatar"] = *opt.Avatar
	}

	return updateField
}
