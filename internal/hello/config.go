package hello

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel       string `envconfig:"LOG_LEVEL" default:"debug"`
	Listen         string `envconfig:"LISTEN" default:":8080"`
	ListenInternal string `envconfig:"LISTEN_INTERNAL" default:":8000"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Error().Err(err).Msg("failed to reading config")

		return nil, err
	}

	return &config, nil
}
