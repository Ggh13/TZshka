package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	authRepo "TZshka/internal/auth/repository"
	authRoute "TZshka/internal/auth/route"
	authService "TZshka/internal/auth/service"
	"TZshka/internal/config"
	historyRepo "TZshka/internal/history/repository"
	historyRoute "TZshka/internal/history/route"
	historyService "TZshka/internal/history/service"
<<<<<<< HEAD
=======
	llmRoute "TZshka/internal/llm/route"
	llmService "TZshka/internal/llm/service"
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
	"TZshka/internal/migrations"
	sessionRepo "TZshka/internal/session/repository"
	sessionRoute "TZshka/internal/session/route"
	sessionService "TZshka/internal/session/service"
	"TZshka/pkg/postgres"
	"TZshka/pkg/registr"
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

	authRepoInstance := authRepo.New(pgDB, zapLog)
	authSvc := authService.New(authRepoInstance, cfg.Secret, zapLog)
	authHandler := authRoute.New(authSvc, zapLog)

	tokenSvc := registr.NewTokenService(cfg.Secret)

	sessionRepoInstance := sessionRepo.New(pgDB, zapLog)
	sessionSvc := sessionService.New(sessionRepoInstance, zapLog)
	sessionHandler := sessionRoute.New(sessionSvc, tokenSvc, zapLog)

	httpClient := &http.Client{}
	llmSvc := llmService.New(httpClient, cfg.LLMURL, zapLog)
	llmHandler := llmRoute.New(llmSvc, zapLog)

	historyRepoInstance := historyRepo.New(pgDB, zapLog)
	historySvc := historyService.New(historyRepoInstance, zapLog)
	historyHandler := historyRoute.New(historySvc, sessionRepoInstance, tokenSvc, zapLog)

	tokenSvc := registr.NewTokenService(cfg.Secret)

	sessionRepoInstance := sessionRepo.New(pgDB, zapLog)
	sessionSvc := sessionService.New(sessionRepoInstance, zapLog)
	sessionHandler := sessionRoute.New(sessionSvc, tokenSvc, zapLog)

	historyRepoInstance := historyRepo.New(pgDB, zapLog)
	historySvc := historyService.New(historyRepoInstance, zapLog)
	historyHandler := historyRoute.New(historySvc, zapLog)

	r := gin.Default()
	r.Use(corsMiddleware())

	authHandler.RegisterRoutes(r)
	sessionHandler.RegisterRoutes(r)
<<<<<<< HEAD
=======
	llmHandler.RegisterRoutes(r)
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
	historyHandler.RegisterRoutes(r)

	zapLog.Info("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		zapLog.Fatal("Server error", zap.Error(err))
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func runMigrations(ctx context.Context, db *pgxpool.Pool, log *zap.Logger) {
	log.Info("Running migrations...")

	if err := migrations.Run(ctx, db); err != nil {
		log.Error("Migration failed", zap.Error(err))
	}

	log.Info("Migrations completed")
}
