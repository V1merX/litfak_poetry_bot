package bot

import (
	"context"
	"log/slog"

	"github.com/V1merX/litfak_poetry_bot/internal/repositories"
	"github.com/V1merX/litfak_poetry_bot/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type bot struct {
	b        *telego.Bot
	updates  <-chan telego.Update
	log      *slog.Logger
	botToken string
	h        *th.BotHandler
	storage  *pgxpool.Pool
}

func NewBot(log *slog.Logger, botToken string, storage *pgxpool.Pool) *bot {
	return &bot{
		botToken: botToken,
		log:      log,
		storage:  storage,
	}
}

func (b *bot) Start() error {
	const op = "internal.bot.Start"
	b.log = b.log.With("op", op)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b.log.Debug("initializing bot")

	tgBot, err := telego.NewBot(b.botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		b.log.Error("failed to create bot", "error", err)
		return err
	}
	b.b = tgBot

	b.log.Debug("setting up updates channel")
	updates, err := b.b.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		b.log.Error("failed to setup updates channel", "error", err)
		return err
	}
	b.updates = updates

	b.log.Debug("creating bot handler")
	handler, err := th.NewBotHandler(b.b, b.updates)
	if err != nil {
		b.log.Error("failed to create bot handler", "error", err)
		return err
	}
	b.h = handler

	userRepository := repositories.NewUserRepository(b.log, b.storage)
	userService := services.NewUserService(b.log, userRepository)

	poemRepository := repositories.NewPoemRepository(b.log, b.storage)
	poemService := services.NewPoemService(b.log, poemRepository)

	b.log.Debug("initializing handlers")
	b.InitHandlers()

	go func() {
		b.WorkerStart(userService, poemService)
	}()

	go func() {
		b.log.Info("starting bot handler")
		if err := b.h.Start(); err != nil {
			b.log.Error("bot handler stopped with error", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	b.log.Info("shutting down bot")

	return nil
}
