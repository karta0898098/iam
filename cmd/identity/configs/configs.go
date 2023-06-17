package configs

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/karta0898098/iam/pkg/db"
	"github.com/karta0898098/iam/pkg/http"
	"github.com/karta0898098/iam/pkg/logging"
)

// Configurations define this application need configs
type Configurations struct {
	Database db.Config      `mapstructure:"database"`
	HTTP     http.Config    `mapstructure:"http"`
	Log      logging.Config `mapstructure:"log"`
	GRPC     GRPC           `mapstructure:"grpc"`
}

type GRPC struct {
	Port string `mapstructure:"port"`
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
