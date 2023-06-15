package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/karta0898098/iam/cmd/identity/configs"
	pb "github.com/karta0898098/iam/pb/identity"
	"github.com/karta0898098/iam/pkg/app/identity/endpoints"
	transportgrpc "github.com/karta0898098/iam/pkg/app/identity/transports/grpc"
	transportshttp "github.com/karta0898098/iam/pkg/app/identity/transports/http"
	"github.com/karta0898098/iam/pkg/db"
	pkggrpc "github.com/karta0898098/iam/pkg/grpc"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// read all config
	config := configs.NewConfig("")

	// setup default logger
	logger := logging.Setup(config.Log)

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
	app.MakeRouter()

	wg := &sync.WaitGroup{}

	go app.startHttpServer(ctx, wg)
	go app.startGRPCServer(ctx, wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cancel()
	wg.Wait()
}

// startHttpServer wrap http start
func (app *Application) startHttpServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	app.logger.Info().Msgf("start http server on %v", app.config.HTTP.Port)

	go func() {
		err := app.httpServer.Start(app.config.HTTP.Port)
		if err != nil {
			app.logger.Error().Err(err).Msg("http server shutdown ...")
		}
	}()

	<-ctx.Done()

	// gracefulShutdown wrap graceful shutdown http server
	// timeout 5 sec will direct close server
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

func (app *Application) startGRPCServer(ctx context.Context, wg *sync.WaitGroup) {
	var (
		server *grpc.Server
	)

	wg.Add(1)
	defer wg.Done()

	port := ":9091"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		app.logger.Panic().Err(err).Msgf("failed to listen on prot=%v", port)
	}

	server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			kitgrpc.Interceptor,
			pkggrpc.UnaryServerLoggerInterceptor(app.logger),
		),
	)
	pb.RegisterIdentityServiceServer(server, transportgrpc.MakeGRPCServer(app.endpoints))
	reflection.Register(server)

	app.logger.Info().Msgf("start grpc server on %v", port)

	go func() {
		// service connections
		err = server.Serve(listener)
		if err != nil {
			app.logger.Error().Msgf("grpc serve : %s\n", err)
		}
	}()

	<-ctx.Done()

	// ignore error since it will be "Err shutting down server : context canceled"
	server.GracefulStop()

	app.logger.Info().Msgf("grpc server gracefully stopped")
}
