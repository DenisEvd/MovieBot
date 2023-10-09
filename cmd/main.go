package main

import (
	"MovieBot/internal/clients/kinopoisk"
	telegramClient "MovieBot/internal/clients/telegram"
	eventConsumer "MovieBot/internal/consumer/event-consumer"
	"MovieBot/internal/events/movie_fetcher"
	"MovieBot/internal/events/processor"
	"MovieBot/internal/events/tg_fetcher"
	"MovieBot/internal/logger"
	"MovieBot/internal/storage"
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

	db, err := storage.NewPostgresDB(storage.Config{
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

	store := storage.NewStorage(db)

	movieFetcher := movie_fetcher.NewKpFetcher(movieApi, 4)
	eventsFetcher := tg_fetcher.NewFetcher(tgClient)
	eventsProcessor := processor.NewTgProcessor(tgClient, movieFetcher, store)

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
