package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"TZshka/internal/auth/repository"
	"TZshka/internal/auth/route"
	"TZshka/internal/auth/service"
	"TZshka/internal/config"
	"TZshka/internal/migrations"
	"TZshka/pkg/postgres"
)

func main() {
	ctx := context.Background()

	zapLog, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer zapLog.Sync()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		zapLog.Fatal("Failed load config", zap.Error(err))
	}
	zapLog.Info("Succesfully load config")

	pgDB, err := postgres.NewPostgres(ctx, &cfg.Postgres)
	if err != nil {
		zapLog.Fatal("Failed connect to postgres DB", zap.Error(err))
	}
	if err := pgDB.Ping(ctx); err != nil {
		zapLog.Fatal("Failed ping pgDB", zap.Error(err))
	}
	zapLog.Info("Succesfully connected to pgDB")

	runMigrations(ctx, pgDB, zapLog)

	authRepo := repository.New(pgDB, zapLog)
	authSvc := service.New(authRepo, cfg.Secret, zapLog)
	authHandler := route.New(authSvc, zapLog)

	r := gin.Default()

	authHandler.RegisterRoutes(r)

	zapLog.Info("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		zapLog.Fatal("Server error", zap.Error(err))
	}
}

func runMigrations(ctx context.Context, db *pgxpool.Pool, log *zap.Logger) {
	log.Info("Running migrations...")

	if err := migrations.Run(ctx, db); err != nil {
		log.Error("Migration failed", zap.Error(err))
	}

	log.Info("Migrations completed")
}
