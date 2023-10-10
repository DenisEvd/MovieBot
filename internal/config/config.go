package config

import (
	"MovieBot/internal/logger"
	"MovieBot/internal/storage/postgres"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configPath     = "configs"
	configFileName = "config"
)

type Config struct{}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.initConfig()

	return cfg
}

func (c *Config) initConfig() {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFileName)

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("error init config", zap.Error(err))
	}

	if err := viper.BindEnv("db_password"); err != nil {
		logger.Fatal("error bind env", zap.Error(err))
	}

	if err := viper.BindEnv("kinopoisk_token"); err != nil {
		logger.Fatal("error bind env", zap.Error(err))
	}

	if err := viper.BindEnv("tg_token"); err != nil {
		logger.Fatal("error bind env", zap.Error(err))
	}
}

func (c *Config) DB() *postgres.Config {
	return &postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db_password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
}

func (c *Config) Tg() string {
	return viper.GetString("tg_token")
}

func (c *Config) Kinopoisk() string {
	return viper.GetString("kinopoisk_token")
}
