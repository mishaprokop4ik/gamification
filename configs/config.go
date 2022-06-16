package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host     string `envconfig:"POSTGRES_HOST"          default:"localhost"`
	Port     string `envconfig:"POSTGRES_PORT"          default:"5432"`
	Username string `envconfig:"POSTGRES_USER"          default:"postgres"`
	Password string `envconfig:"POSTGRES_PASSWORD"      default:"rootroot"`
	DBName   string `envconfig:"POSTGRES_DATABASE"      default:"postgres"`
	SSLMode  string
}

func Init(path string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	return viper.ReadInConfig()
}
