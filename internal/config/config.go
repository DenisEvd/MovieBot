package config

import (
	"MovieBot/internal/logger"
	"MovieBot/internal/storage/postgres"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

type Config struct{}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.initConfig()

	return cfg
}

func (c *Config) initConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("error init config", zap.Error(err))
	}

	if err := godotenv.Load(); err != nil {
		logger.Fatal("error read .env", zap.Error(err))
	}
}

func (c *Config) DB() postgres.Config {
	return postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
}

func (c *Config) Tg() string {
	return os.Getenv("TG_TOKEN")
}

func (c *Config) Kinopoisk() string {
	return os.Getenv("KINOPOISK_TOKEN")
}
