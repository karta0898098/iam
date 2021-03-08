package access

import (
	"github.com/karta0898098/iam/pkg/access/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(service.NewAccessService),
)
