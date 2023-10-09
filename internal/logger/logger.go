package logger

import (
	"go.uber.org/zap"
	"log"
)

var logger *zap.Logger

func init() {
	locLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("error logger init")
	}

	logger = locLogger
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
