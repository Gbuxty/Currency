package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres *PostgresConfig
	Server   *ServerConfig
	Grinex   *GrinexConfig
}

func NewConfig() (*Config, error) {
	configPath := flag.String("config", "configs/local.yml", "path to configuration file")
	flag.Parse()

	viper.SetConfigFile(*configPath)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()


	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	return &Config{
		Postgres: newPostgresConfig(),
		Server:   newServerConfig(),
		Grinex:   newGrinexConfig(),
	}, nil
}
