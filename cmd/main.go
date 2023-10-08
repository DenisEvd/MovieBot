package main

import (
	"MovieBot/internal/pkg/clients/kinopoisk"
	telegramClient "MovieBot/internal/pkg/clients/telegram"
	eventConsumer "MovieBot/internal/pkg/consumer/event-consumer"
	kinopoiskFetch "MovieBot/internal/pkg/events/kinopoisk"
	"MovieBot/internal/pkg/events/telegram"
	"MovieBot/internal/pkg/storage"
	"flag"
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
	logger := initLogger()

	if err := initConfig(); err != nil {
		logger.Fatal("error read config", zap.Error(err))
	}

	tgClient := telegramClient.New(hostTg, mustToken())

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

	movieApi := kinopoisk.NewKp(hostKp, os.Getenv("KINOPOISK_TOKEN"))

	store := storage.NewStorage(db)

	movieFetcher := kinopoiskFetch.NewKpFetcher(movieApi, 4)
	eventsFetcher := telegram.NewFetcher(tgClient)
	eventsProcessor := telegram.NewProcessor(tgClient, movieFetcher, store)

	consumer := eventConsumer.New(logger, eventsFetcher, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatal("error starting consumer", zap.Error(err))
	}

}

func mustToken() string {
	token := flag.String("tg-token", "", "token for access to telegram bot")

	flag.Parse()

	return *token
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
