package main

import (
	"MovieBot/internal/clients/kinopoisk"
	telegramClient "MovieBot/internal/clients/telegram"
	"MovieBot/internal/config"
	eventConsumer "MovieBot/internal/consumer/event-consumer"
	"MovieBot/internal/events/movie_fetcher"
	"MovieBot/internal/events/processor/telegram"
	"MovieBot/internal/events/tg_fetcher"
	"MovieBot/internal/logger"
	"MovieBot/internal/storage/postgres"
	"go.uber.org/zap"
)

const (
	hostTg    = "api.telegram.org"
	hostKp    = "api.kinopoisk.dev"
	batchSize = 100
)

func main() {
	conf := config.NewConfig()

	db, err := postgres.NewPostgresDB(conf.DB())
	if err != nil {
		logger.Fatal("error sign in db", zap.Error(err))
	}

	tgClient := telegramClient.New(hostTg, conf.Tg())

	movieApi := kinopoisk.NewKp(hostKp, conf.Kinopoisk())

	store := postgres.New(db)

	movieFetcher := movie_fetcher.NewKpFetcher(movieApi, 4)
	eventsFetcher := tg_fetcher.NewFetcher(tgClient)
	eventsProcessor := telegram.NewTgProcessor(tgClient, movieFetcher, store)

	consumer := eventConsumer.New(eventsFetcher, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatal("error starting consumer", zap.Error(err))
	}
}
