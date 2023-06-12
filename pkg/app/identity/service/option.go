package service

import (
	"github.com/karta0898098/iam/pkg/app/identity/entity"
)

type SigninOption struct {
	IPAddress   string
	Platform    string
	Device      entity.Device
	IdpProvider string
}

type SignupOption struct {
	Nickname  string        // Nickname user nickname
	FirstName string        // FirstName user first name
	LastName  string        // LastName user last name
	Email     string        // Email user email address
	IPAddress string        // IPAddress from user login
	Platform  string        // Login platform
	Device    entity.Device // Device login information
}
