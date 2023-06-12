// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/karta0898098/iam/configs"
	"github.com/karta0898098/iam/pkg/app/identity/endpoints"
	"github.com/karta0898098/iam/pkg/app/identity/repository"
	"github.com/karta0898098/iam/pkg/app/identity/service"
	"github.com/karta0898098/iam/pkg/db"
	"github.com/rs/zerolog"
)

// Injectors from wire.go:

func NewApp(logger zerolog.Logger, cfg configs.Configurations, conn db.Connection) *Application {
	repositoryRepository := repository.New(conn)
	identityService := service.New(repositoryRepository)
	endpointsEndpoints := endpoints.New(identityService)
	application := NewApplication(logger, cfg, endpointsEndpoints)
	return application
}
