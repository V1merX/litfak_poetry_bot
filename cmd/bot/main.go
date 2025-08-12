package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/V1merX/litfak_poetry_bot/internal/bot"
	"github.com/V1merX/litfak_poetry_bot/internal/config"
	"github.com/V1merX/litfak_poetry_bot/internal/storage/postgresql"
)

var log *slog.Logger

func init() {
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	if err := run(); err != nil {
		log.Error("failed to start application", slog.String("error", err.Error()))
	}
}

func run() error {
	log.Info("starting application")

	configPath := os.Getenv("CONFIG_PATH")

	cfg := config.MustLoad(log, configPath)

	pool, err := postgresql.NewPool(context.Background(), log, &cfg.PostgreSQL)
	if err != nil {
		log.Error("failed to connect to database")
		return err
	}

	_ = pool

	bot := bot.NewBot(log, cfg.Telegram.BotToken)
	if err := bot.Start(); err != nil {
		return err
	}

	return nil
}
