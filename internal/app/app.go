package app

import (
	"database/sql"
	"fmt"

	"em_subscription_test/config"
	"em_subscription_test/db"
	"em_subscription_test/handlers"
	"em_subscription_test/internal/repository"
	"em_subscription_test/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeApp() (*gin.Engine, error) {
	cfg := config.Load()

	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Printf("Invalid log level %s, using info", cfg.LogLevel)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	database, err := db.NewDB(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
		return nil, err
	}

	if err := runMigrations(database.DB.DB, logger); err != nil {
		logger.WithError(err).Fatal("Failed to run migrations")
		return nil, err
	}

	repo := repository.NewSubscriptionRepository(database.DB)

	svc := service.NewSubscriptionService(repo, logger)

	h := handlers.NewHandler(svc, logger)

	g := gin.Default()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := g.Group("/api/v1")
	subscriptions := api.Group("/subscriptions")
	{
		subscriptions.POST("", h.CreateSubscription)
		subscriptions.GET("", h.ListSubscriptions)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.DELETE("/:id", h.DeleteSubscription)
		subscriptions.POST("/total-cost", h.GetTotalCost)
	}

	return g, nil
}

func runMigrations(db *sql.DB, logger *logrus.Logger) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Migrations completed successfully")
	return nil
}
