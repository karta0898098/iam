package identity

import (
	"github.com/karta0898098/iam/pkg/identity/repository"
	"github.com/karta0898098/iam/pkg/identity/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(service.NewIdentityService),
	fx.Provide(repository.NewIdentityRepository),
)
