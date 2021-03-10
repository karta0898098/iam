package serve

import (
	"context"
	"flag"
	"time"

	"github.com/karta0898098/iam/configs"
	"github.com/karta0898098/iam/pkg/access"
	deliveryHttp "github.com/karta0898098/iam/pkg/delivery/http"
	"github.com/karta0898098/iam/pkg/identity"
	"github.com/karta0898098/kara/redis"

	"github.com/karta0898098/kara/db/rw/db"
	"github.com/karta0898098/kara/http"
	"github.com/karta0898098/kara/zlog"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func Run() {
	var (
		path    string
		idUtils *snowflake.Node
	)
	flag.StringVar(&path, "p", "", "serve -p ./deployments/config")
	flag.Parse()

	idUtils, _ = snowflake.NewNode(1)

	config := configs.NewConfig(path)

	app := fx.New(
		fx.Supply(config),
		fx.Supply(idUtils),
		fx.Provide(db.NewConnection),
		fx.Provide(http.NewEcho),
		fx.Provide(redis.NewRedis),
		fx.Provide(deliveryHttp.NewHandler),
		identity.Module,
		access.Module,
		fx.Invoke(zlog.Setup),
		fx.Invoke(deliveryHttp.SetupRoute),
		fx.Invoke(http.RunEcho),
	)
	app.Run()

	log.Info().Msg("Graceful Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Info().Msgf("Server Shutdown: %v", err)
	}

	log.Info().Msg("Server exiting")
}
