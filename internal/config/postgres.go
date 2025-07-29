package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

func newPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetInt("postgres.port"),
		User:     viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		DBName:   viper.GetString("postgres.database"),
		SSLMode:  viper.GetString("postgres.sslmode"),
	}
}

func (p *PostgresConfig) ToDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode)
}
