package main

import (
	"MovieBot/internal/pkg/clients/kinopoisk"
	telegramClient "MovieBot/internal/pkg/clients/telegram"
	eventConsumer "MovieBot/internal/pkg/consumer/event-consumer"
	"MovieBot/internal/pkg/events/telegram"
	postgresql "MovieBot/internal/pkg/storage/postgres"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	hostTg    = "api.telegram.org"
	hostKp    = "api.kinopoisk.dev"
	batchSize = 100
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("can't read config: %s", err.Error())
	}

	fmt.Println(os.Getenv("TG_TOKEN"))
	tgClient := telegramClient.New(hostTg, mustToken())

	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't read .env: %s", err.Error())
	}

	db, err := postgresql.NewPostgresDB(postgresql.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("can't sign in db: %s", err.Error())
	}

	movieApi := kinopoisk.NewKp(hostKp, os.Getenv("KINOPISK_TOKEN"))

	postgres := postgresql.New(db)

	eventsProcessor := telegram.New(tgClient, movieApi, postgres)

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
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
