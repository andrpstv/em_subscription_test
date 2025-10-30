package main

import (
	"log"

	"em_subscription_test/config"
	"em_subscription_test/internal/app"

	_ "em_subscription_test/docs"

	"github.com/sirupsen/logrus"
)

// @title Subscription API
// @version 1.0
// @description REST API for managing user subscriptions
// @host localhost:8080
// @BasePath /api/v1/

func main() {
	g, err := app.InitializeApp()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	cfg := config.Load()

	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Printf("Invalid log level %s, using info", cfg.LogLevel)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.WithField("port", cfg.ServerPort).Info("Starting server")
	if err := g.Run(":" + cfg.ServerPort); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
