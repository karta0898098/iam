package identity

import (
	"github.com/google/wire"

	"github.com/karta0898098/iam/pkg/app/identity/endpoints"
	"github.com/karta0898098/iam/pkg/app/identity/repository"
	"github.com/karta0898098/iam/pkg/app/identity/service"
)

var DefaultProvider = wire.NewSet(
	endpoints.New,
	service.New,
	repository.New,
)
