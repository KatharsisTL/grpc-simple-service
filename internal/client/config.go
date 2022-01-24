package client

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"time"
)

type Config struct {
	LogLevel                  string        `envconfig:"LOG_LEVEL" default:"debug"`
	Listen                    string        `envconfig:"LISTEN" default:":8081"`
	ListenInternal            string        `envconfig:"LISTEN_INTERNAL" default:":8001"`
	HelloClientUrl            string        `envconfig:"HELLO_CLIENT_URL" default:"localhost:8080"`
	HelloClientConnectTimeout time.Duration `envconfig:"HELLO_CLIENT_CONNECT_TIMEOUT" default:"10s"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Error().Err(err).Msg("failed to reading config")

		return nil, err
	}

	return &config, nil
}
