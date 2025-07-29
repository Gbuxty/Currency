package app

import (
	"GetCurrency/internal/adapter/grinex"
	"GetCurrency/internal/adapter/repository"
	"GetCurrency/internal/config"
	"GetCurrency/internal/service"
	"GetCurrency/internal/transport/grpc/handlers"
	"GetCurrency/internal/transport/grpc/server"
	"GetCurrency/pkg/logger"
	"GetCurrency/pkg/pg"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx,cancel := context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()

	db, err := pg.NewClient(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	repo := repository.NewStorage(db)

	grinex := grinex.NewGrinex(cfg.Grinex, logger)

	svc := service.NewRateService(repo, grinex)

	handlers := handlers.NewRateHandler(svc, logger)

	srv := server.NewGrpcServer(logger, cfg.Server, handlers)


	if err := srv.Start(ctx); err != nil {
		logger.Errorf("error server :%w", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	srv.Stop()

	return nil
}
