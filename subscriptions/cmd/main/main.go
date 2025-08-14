package main

import (
	"context"
	"github.com/ekkserapopova/subscriptions/internal/config"
	"github.com/ekkserapopova/subscriptions/internal/pkg/db"
	"github.com/ekkserapopova/subscriptions/internal/pkg/migrations"
	"github.com/ekkserapopova/subscriptions/internal/pkg/server"
	subscriptionHandler "github.com/ekkserapopova/subscriptions/internal/services/subscriptions/delivery/http"
	subscriptionRepository "github.com/ekkserapopova/subscriptions/internal/services/subscriptions/repo"
	subscriptionUseCase "github.com/ekkserapopova/subscriptions/internal/services/subscriptions/usecase"
	"github.com/ekkserapopova/subscriptions/pkg/builder"
	"github.com/ekkserapopova/subscriptions/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title Subscriptions API
// @version 1.0
// @description API для управления подписками

// @host localhost:8080
// @BasePath /api/v1
func main() {
	app := fx.New(
		fx.Provide(
			logger.SetupLogger,
			builder.SetupBuilder,
			config.MustLoad,

			db.NewPostgresPool,
			db.NewPostgresConnect,

			server.NewRouter,

			subscriptionHandler.NewHandler,
			subscriptionUseCase.NewUseCase,
			subscriptionRepository.NewRepository,
		),

		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),

		fx.Invoke(
			server.RunServer,
			migrations.RunMigrations,
		),
	)

	ctx := context.Background()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	<-stop
	app.Stop(ctx)
}
