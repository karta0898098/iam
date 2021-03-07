package http

import (
	"github.com/karta0898098/iam/pkg/identity/domain"
)

type Handler struct {
	identityService domain.IdentityService
}

func NewHandler(identityService domain.IdentityService) *Handler {
	return &Handler{identityService: identityService}
}
