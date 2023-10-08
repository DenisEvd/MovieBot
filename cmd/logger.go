package main

import (
	"go.uber.org/zap"
	"log"
)

func initLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("error logger init")
	}

	return logger
}
