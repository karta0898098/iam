package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"time"
	"unicode"

	"github.com/karta0898098/iam/pkg/identity/domain"
	"github.com/karta0898098/iam/pkg/identity/repository"
	"github.com/karta0898098/iam/pkg/identity/repository/model"

	"github.com/karta0898098/kara/errors"

	"github.com/bwmarrin/snowflake"
)

var _ domain.IdentityService = &identityService{}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type identityService struct {
	repo    repository.IdentityRepository
	idUtils *snowflake.Node
}

func NewIdentityService(repo repository.IdentityRepository, idUtils *snowflake.Node) domain.IdentityService {
	return &identityService{
		repo:    repo,
		idUtils: idUtils,
	}
}

const (
	accountLengthMax = 16
	accountLengthMin = 8

	passwordLengthMax = 16
	passwordLengthMin = 8

	nameLengthMax = 20

	avatarURLLengthMax = 100
)

func (srv *identityService) Login(ctx context.Context, account string, password string) (profile *domain.Profile, err error) {
	var (
		opts *model.GetProfileOption
	)

	if len(account) < accountLengthMin || len(account) > accountLengthMax {
		return nil, errors.ErrInvalidInput.Build("input account length not equal rule")
	}

	if !srv.ValidatePassword(password, passwordLengthMin, passwordLengthMax) {
		return nil, errors.ErrInvalidInput.Build("input password length not equal rule")
	}

	opts = &model.GetProfileOption{
		Account: account,
	}
	profile, err = srv.repo.GetProfile(ctx, opts)
	if err != nil {
		return nil, err
	}

	if profile.Password != srv.ToSHA256(password) {
		return nil, errors.ErrUnauthorized.Build("input password not equal")
	}

	if profile.SuspendStatus == domain.Suspend || profile.SuspendStatus == domain.UnknownStatus {
		return nil, errors.ErrUnauthorized.Build("account status is %s", profile.SuspendStatus.ToString())
	}

	return profile, nil
}

func (srv *identityService) Signup(ctx context.Context, account string, password string, opts *domain.SignupOption) (*domain.Profile, error) {
	var (
		queryOpts    *model.GetProfileOption
		createParams *model.Profile
	)

	createParams = &model.Profile{}

	if len(account) < accountLengthMin || len(account) > accountLengthMax {
		return nil, errors.ErrInvalidInput.Build("input account length not equal rule")
	}

	if !srv.ValidatePassword(password, passwordLengthMin, passwordLengthMax) {
		return nil, errors.ErrInvalidInput.Build("input password length not equal rule")
	}

	if opts != nil {
		if opts.Email != "" {
			if !srv.ValidateEmail(opts.Email) {
				return nil, errors.ErrInvalidInput.Build("input email validate not correct")
			}
		}
		if len(opts.Nickname) > nameLengthMax {
			return nil, errors.ErrInvalidInput.Build("input nickname length too many")
		}
		if len(opts.FirstName) > nameLengthMax {
			return nil, errors.ErrInvalidInput.Build("input first name length too many")
		}
		if len(opts.LastName) > nameLengthMax {
			return nil, errors.ErrInvalidInput.Build("input last name length too many")
		}
		if len(opts.Avatar) > avatarURLLengthMax {
			return nil, errors.ErrInvalidInput.Build("input avatar url length too many")
		}

		createParams.Nickname = opts.Nickname
		createParams.FirstName = opts.FirstName
		createParams.LastName = opts.LastName
		createParams.Email = opts.Email
		createParams.Avatar = opts.Avatar
	}

	queryOpts = &model.GetProfileOption{
		Account: account,
	}
	profile, err := srv.repo.GetProfile(ctx, queryOpts)
	if err != nil {
		if !errors.Is(errors.ErrResourceNotFound, err) {
			return nil, err
		}
	}

	if profile != nil {
		return nil, errors.ErrConflict.Build("account is duplicate")
	}

	createParams.ID = srv.idUtils.Generate().Int64()
	createParams.Account = account
	createParams.Password = srv.ToSHA256(password)
	createParams.CreatedAt = time.Now().UTC()
	createParams.UpdatedAt = time.Now().UTC()
	createParams.SuspendStatus = domain.Unsuspend

	profile, err = srv.repo.CreateProfile(ctx, createParams)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (srv *identityService) UpdateProfile(ctx context.Context, id int64, opts *domain.UpdateProfileOption) error {
	var (
		queryOpts    *model.GetProfileOption
		updateParams *model.UpdateProfileOption
	)

	if id == 0 {
		return errors.ErrInvalidInput.Build("input profile id is zero")
	}

	if opts == nil {
		return nil
	}

	if opts.Email != nil {
		if !srv.ValidateEmail(*opts.Email) {
			return errors.ErrInvalidInput.Build("input email validate not correct")
		}
	}
	if opts.Nickname != nil {
		if len(*opts.Nickname) > nameLengthMax {
			return errors.ErrInvalidInput.Build("input nickname length too many")
		}
	}
	if opts.FirstName != nil {
		if len(*opts.FirstName) > nameLengthMax {
			return errors.ErrInvalidInput.Build("input first name length too many")
		}
	}
	if opts.LastName != nil {
		if len(*opts.LastName) > nameLengthMax {
			return errors.ErrInvalidInput.Build("input last name length too many")
		}
	}
	if opts.Avatar != nil {
		if len(*opts.Avatar) > avatarURLLengthMax {
			return errors.ErrInvalidInput.Build("input avatar url length too many")
		}
	}

	queryOpts = &model.GetProfileOption{
		ID: id,
	}
	profile, err := srv.repo.GetProfile(ctx, queryOpts)
	if err != nil {
		return err
	}

	updateParams = &model.UpdateProfileOption{
		Nickname:  opts.Nickname,
		FirstName: opts.FirstName,
		LastName:  opts.LastName,
		Email:     opts.Email,
		Avatar:    opts.Avatar,
		UpdatedAt: time.Now().UTC(),
	}

	err = srv.repo.UpdateProfile(ctx, profile.ID, updateParams)
	if err != nil {
		return err
	}

	return nil
}

// ValidatePassword check input password match rule
// password length min <= len <= max
// password need a alphabet
// password don't have any space or special symbol
func (srv *identityService) ValidatePassword(password string, min, max int) bool {
	var (
		number  bool
		special bool
	)
	letters := 0
	for _, c := range password {
		switch {
		case unicode.IsSpace(c):
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			letters++
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsLetter(c):
			letters++
		}
	}
	return (number) && (!special) && letters <= max && letters >= min
}

// ValidateEmail validate email format
func (srv *identityService) ValidateEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// ToSHA256 for password encrypt
func (srv *identityService) ToSHA256(password string) string {
	var (
		result string
	)
	h := sha256.New()
	h.Write([]byte(password))
	result = hex.EncodeToString(h.Sum(nil))
	return result
}
