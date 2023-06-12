package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/karta0898098/iam/pkg/app/identity/entity"
	"github.com/karta0898098/iam/pkg/app/identity/service"
)

// Endpoints contain all identity endpoint
type Endpoints struct {
	SigninEndpoint endpoint.Endpoint
	SignupEndpoint endpoint.Endpoint
}

// New endpoints
func New(svc service.IdentityService) (ep Endpoints) {
	signinEndpoint := MakeSigninEndpoint(svc)
	signinEndpoint = LoggingMiddleware("Signin")(signinEndpoint)
	ep.SigninEndpoint = signinEndpoint

	signupEndpoint := MakeSignupEndpoint(svc)
	signupEndpoint = LoggingMiddleware("Signup")(signupEndpoint)
	ep.SignupEndpoint = signupEndpoint

	return ep
}

// SigninRequest define signin request
type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	IPAddress   string        `json:"ip_address"`
	Platform    string        `json:"platform"`
	IdpProvider string        `json:"idp_provider"`
	Device      entity.Device `json:"device"`
}

// SigninResponse define signup response
type SigninResponse struct {
	IDToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// MakeSigninEndpoint make signin endpoint
func MakeSigninEndpoint(svc service.IdentityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SigninRequest)

		_, err = svc.Signin(ctx, req.Username, req.Password, &service.SigninOption{
			IPAddress:   req.IPAddress,
			Platform:    req.Platform,
			Device:      req.Device,
			IdpProvider: req.IdpProvider,
		})
		if err != nil {
			return nil, err
		}

		return &SigninResponse{
			// AccessToken:  identity.NewAccessToken(),
			// RefreshToken: identity.NewRefreshToken(),
		}, nil
	}
}

// SignupRequest define signup response
type SignupRequest struct {
	Username  string        `json:"username,omitempty"`
	Password  string        `json:"password,omitempty"`
	Nickname  string        `json:"nickname,omitempty"`
	FirstName string        `json:"first_name,omitempty"`
	LastName  string        `json:"last_name,omitempty"`
	Email     string        `json:"email,omitempty"`
	Platform  string        `json:"platform,omitempty"`
	Device    entity.Device `json:"device"`
}

// SignupResponse define signup response
type SignupResponse struct {
	IDToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// MakeSignupEndpoint make signup endpoint
func MakeSignupEndpoint(svc service.IdentityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SignupRequest)

		identity, err := svc.Signup(ctx, req.Username, req.Password, &service.SignupOption{
			Nickname:  req.Nickname,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			IPAddress: "", // TODO
			Platform:  req.Platform,
			Device:    req.Device,
		})
		if err != nil {
			return nil, err
		}

		return &SignupResponse{
			AccessToken:  identity.NewAccessToken(),
			RefreshToken: identity.NewRefreshToken(),
		}, nil
	}
}