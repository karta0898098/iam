package http

import (
	"context"
	"net/http"

	access "github.com/karta0898098/iam/pkg/access/domain"
	identity "github.com/karta0898098/iam/pkg/identity/domain"

	"github.com/karta0898098/kara/errors"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	identityService identity.IdentityService
	accessService   access.AccessService
}

func NewHandler(identityService identity.IdentityService, accessService access.AccessService) *Handler {
	return &Handler{
		identityService: identityService,
		accessService:   accessService,
	}
}

func (h *Handler) LoginEndpoint(c echo.Context) error {
	type (
		LoginRequest struct {
			Account  string `json:"account"`
			Password string `json:"password"`
		}

		LoginResponse struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			TokenType    string `json:"tokenType,omitempty"`
			ExpireIn     int64  `json:"expireIn"`
		}
	)
	var (
		ctx      context.Context
		request  LoginRequest
		response LoginResponse
	)

	ctx = c.Request().Context()

	if err := c.Bind(&request); err != nil {
		return errors.ErrInvalidInput.BuildWithError(err)
	}

	user, err := h.identityService.Login(ctx, request.Account, request.Password)
	if err != nil {
		return err
	}

	claims, err := h.accessService.CreateAccessTokens(ctx, user)
	if err != nil {
		return err
	}

	response = LoginResponse{
		AccessToken:  claims.AccessToken,
		RefreshToken: claims.RefreshToken,
		ExpireIn:     claims.ExpireIn,
	}

	return c.JSON(http.StatusOK, response)
}
