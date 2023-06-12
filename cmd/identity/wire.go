//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"

	"github.com/karta0898098/iam/configs"
	"github.com/karta0898098/iam/pkg/app/identity"
	"github.com/karta0898098/iam/pkg/db"
)

func NewApp(logger zerolog.Logger, cfg configs.Configurations, conn db.Connection) *Application {
	wire.Build(
		identity.DefaultProvider,
		NewApplication,
	)
	return nil
}
