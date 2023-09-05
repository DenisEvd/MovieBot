package main

import (
	telegramClient "MovieBot/clients/telegram"
	eventConsumer "MovieBot/consumer/event-consumer"
	"MovieBot/events/telegram"
	postgresql "MovieBot/storage/postgres"
	"flag"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	host      = "api.telegram.org"
	batchSize = 100
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("can't read config: %s", err.Error())
	}

	tgClient := telegramClient.New(host, mustToken())

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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
