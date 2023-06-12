package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/karta0898098/iam/configs"
	"github.com/karta0898098/iam/pkg/app/identity/endpoints"
	transportshttp "github.com/karta0898098/iam/pkg/app/identity/transports/http"
	"github.com/karta0898098/iam/pkg/db"
	"github.com/karta0898098/iam/pkg/http"
	"github.com/karta0898098/iam/pkg/http/middleware"
	"github.com/karta0898098/iam/pkg/logging"
)

// Application define application
type Application struct {
	logger     zerolog.Logger
	config     configs.Configurations
	httpServer *echo.Echo
	endpoints  endpoints.Endpoints
}

// NewApplication new application
func NewApplication(
	logger zerolog.Logger,
	config configs.Configurations,
	endpoints endpoints.Endpoints,
) *Application {
	return &Application{
		logger:     logger,
		config:     config,
		httpServer: http.NewEcho(config.HTTP),
		endpoints:  endpoints,
	}
}

func main() {

	// read all config
	config := configs.NewConfig("")

	// setup default logger
	logger := logging.Setup(config.Log)
	logger.Printf("call")

	// init database
	dbConn, err := db.NewConnection(config.Database)
	if err != nil {
		logger.
			Panic().
			Err(err).
			Msg("failed to connection database")
	}

	// make app
	app := NewApp(logger, config, dbConn)
	app.httpServer.Pre(middleware.NewLoggerMiddleware(logger))
	app.httpServer.Use(middleware.NewLoggingMiddleware())
	app.httpServer.Use(middleware.NewErrorHandlingMiddleware())
	app.httpServer.Use(middleware.RecordErrorMiddleware())
	app.MakeRouter()

	go app.startServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	app.gracefulShutdown()
}

// startServer wrap http start
func (app *Application) startServer() {
	app.logger.Info().Msgf("start http server on %v", app.config.HTTP.Port)

	err := app.httpServer.Start(app.config.HTTP.Port)
	if err != nil {
		app.logger.Error().Err(err).Msg("http server shutdown ...")
	}
}

// gracefulShutdown wrap graceful shutdown http server
// timeout 5 sec will direct close server
func (app *Application) gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.httpServer.Shutdown(ctx)
	if err != nil {
		app.logger.
			Error().
			Err(err).
			Msg("failed to http shutdown")
	}
}

// MakeRouter register router into echo http server
func (app *Application) MakeRouter() *Application {
	app.httpServer.POST("/signin", echo.WrapHandler(transportshttp.MakeSignin(app.endpoints)))
	app.httpServer.POST("/signup", echo.WrapHandler(transportshttp.MakeSignup(app.endpoints)))

	return app
}
