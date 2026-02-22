package main

import (
	"go.uber.org/zap"
	"log/slog"
	"os"
)

var password = "secretPassword"
var apiKey = "1234567890"
var token = "abcdef12345"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	slogLogger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("Starting server on port 8080")           // want "log message must start with a lowercase letter"
	slogLogger.Error("Failed to connect to database")     // want "log message must start with a lowercase letter"

	logger.Info("starting server on port 8080")
	slogLogger.Error("failed to connect to database")

	logger.Info("запуск сервера")                         // want "log message must be only in english"
	slogLogger.Error("ошибка подключения к базе данных") // want "log message must be only in english"

	logger.Info("starting server")
	slogLogger.Error("failed to connect to database")

	logger.Info("server started!🚀")                      // want "log message must not contain special symbols or emoji"
	slogLogger.Error("connection failed!!!")             // want "log message must not contain special symbols or emoji"
	slogLogger.Warn("warning: something went wrong...")  // want "log message must not contain special symbols or emoji"

	logger.Info("server started")
	slogLogger.Error("connection failed")
	slogLogger.Warn("something went wrong")

	logger.Info("user password: " + password) // want "log message must not contain potentially sensitive data"
	logger.Debug("api_key=" + apiKey)         // want "log message must not contain potentially sensitive data"
	logger.Info("token: " + token)            // want "log message must not contain potentially sensitive data"

	logger.Info("user authenticated successfully")
	logger.Debug("api request completed")
	logger.Info("token validated")
}