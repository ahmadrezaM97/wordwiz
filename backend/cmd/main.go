package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"wordwiz/config"
	"wordwiz/internal/server"
	"wordwiz/internal/storage/postgres"
	"wordwiz/pkg/logger"

	"github.com/rs/zerolog"
)

func main() {
	cfg := config.Load()

	if err := logger.InitLogger(cfg.App.Environment, cfg.App.LogLevel); err != nil {
		log.Println(fmt.Errorf("init logger failed: %v", err))
	}

	l := logger.Get()

	if err := run(cfg, l); err != nil {
		l.Fatal().Msgf("main run faild: %v", err)
	}
}

func run(cfg config.Config, log *zerolog.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pgStorage, err := postgres.New(ctx, cfg.Postgres.MakeURL())
	if err != nil {
		return fmt.Errorf("pg init new client failed: %w", err)
	}

	defer func() {
		if err := pgStorage.Close(ctx); err != nil {
			log.Error().Msgf("close pg client failed: %v", err)
		}
	}()

	httpServer := server.New(cfg, pgStorage)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Info().Msg("httpServer started...")
		if err := httpServer.Serve(ctx); err != nil {
			log.Fatal().Err(err).Msg("http server failed to function")
		}
	}()

	wg.Wait()
	log.Info().Msg("DONE!")

	return nil
}
