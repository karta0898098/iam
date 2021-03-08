package configs

import (
	"github.com/karta0898098/kara/db/rw/db"
	"github.com/karta0898098/kara/http"
	"github.com/karta0898098/kara/redis"
	"github.com/karta0898098/kara/zlog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"go.uber.org/fx"
)

type Configurations struct {
	fx.Out
	Log      zlog.Config
	Database db.Config
	Http     http.Config
	Redis    redis.Config
}

// NewConfig read configs and create new instance
func NewConfig(path string) Configurations {
	// set file type toml or yaml
	viper.AutomaticEnv()
	viper.SetConfigType("toml")
	var config Configurations

	// user doesn't input any configs
	// then check env export configs path
	if path == "" {
		path = "./deployments/config"
	}

	// check user want setting other configs
	name := viper.GetString("CONFIG_NAME")
	if name == "" {
		name = "app"
	}
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msgf("Error reading configs file, %s", err)
		return config
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
		return config
	}
	return config
}
