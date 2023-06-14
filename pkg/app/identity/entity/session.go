package entity

import (
	"time"
)

const (
	SigninKey = ""
)

type Session struct {
	ID          string
	UserID      string
	CreateAt    time.Time
	UpdateAt    time.Time
	ExpireAt    time.Time
	IPAddress   string
	IdpProvider string
	Platform    string
	Device      Device
}

type Device struct {
	Model     string
	Name      string
	OSVersion string
}

type NewSessionOption func(p *Session)

func NewSession(id string, userID string, ip string, platform string, opts ...NewSessionOption) *Session {
	now := time.Now()
	accessTokenExpiresAt := now.Add(10 * 60 * time.Second)

	s := &Session{
		ID:        id,
		UserID:    userID,
		CreateAt:  now,
		UpdateAt:  now,
		ExpireAt:  accessTokenExpiresAt,
		IPAddress: ip,
		Platform:  platform,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithIdpProvider(idpProvider string) NewSessionOption {
	return func(p *Session) {
		p.IdpProvider = idpProvider
	}
}

func WithDevice(device Device) NewSessionOption {
	return func(p *Session) {
		p.Device = device
	}
}
