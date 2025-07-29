package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port              int           `mapstructure:"port"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
}

func newServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:              viper.GetInt("server.port"),
		ConnectionTimeout: viper.GetDuration("server.connection_timeout"),
	}
}
