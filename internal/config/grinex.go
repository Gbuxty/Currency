package config

import "github.com/spf13/viper"

type GrinexConfig struct {
	Url string `mapstructure:"url"`
}

func newGrinexConfig() *GrinexConfig {
	return &GrinexConfig{
		Url: viper.GetString("api_grinex.url"),
	}
}
