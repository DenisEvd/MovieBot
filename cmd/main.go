package main

import (
	telegramClient "MovieBot/clients/telegram"
	eventConsumer "MovieBot/consumer/event-consumer"
	"MovieBot/events/telegram"
	postgresql "MovieBot/storage/postgres"
	"flag"
	"log"
)

const (
	host      = "api.telegram.org"
	batchSize = 100
)

func main() {
	tgClient := telegramClient.New(host, mustToken())

	cfg := postgresql.Config{}

	db, err := postgresql.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("can't sign in db: %s", err.Error())
	}

	postgres := postgresql.New(db)

	eventsProcessor := telegram.New(tgClient, postgres)

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}

}

func mustToken() string {
	token := flag.String("token-bot-token", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
