package main

import (
	"MovieBot/internal/clients/kinopoisk"
	telegramClient "MovieBot/internal/clients/telegram"
	eventConsumer "MovieBot/internal/consumer/event-consumer"
	kinopoiskFetch "MovieBot/internal/events/kinopoisk"
	telegram2 "MovieBot/internal/events/telegram"
	"MovieBot/internal/logger"
	storage2 "MovieBot/internal/storage"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

const (
	hostTg    = "api.telegram.org"
	hostKp    = "api.kinopoisk.dev"
	batchSize = 100
)

func main() {
	if err := initConfig(); err != nil {
		logger.Fatal("error read config", zap.Error(err))
	}

	if err := godotenv.Load(); err != nil {
		logger.Fatal("error read .env", zap.Error(err))
	}

	db, err := storage2.NewPostgresDB(storage2.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatal("error sign in db", zap.Error(err))
	}

	tgClient := telegramClient.New(hostTg, os.Getenv("TG_TOKEN"))

	movieApi := kinopoisk.NewKp(hostKp, os.Getenv("KINOPOISK_TOKEN"))

	store := storage2.NewStorage(db)

	movieFetcher := kinopoiskFetch.NewKpFetcher(movieApi, 4)
	eventsFetcher := telegram2.NewFetcher(tgClient)
	eventsProcessor := telegram2.NewProcessor(tgClient, movieFetcher, store)

	consumer := eventConsumer.New(eventsFetcher, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatal("error starting consumer", zap.Error(err))
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
