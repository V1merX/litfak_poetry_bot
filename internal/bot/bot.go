package bot

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type bot struct {
	b        *telego.Bot
	updates  <-chan telego.Update
	log      *slog.Logger
	botToken string
	h        *th.BotHandler
}

func NewBot(log *slog.Logger, botToken string) *bot {
	return &bot{
		botToken: botToken,
		log:      log,
	}
}

func (b *bot) Start() error {
	const op = "internal.bot.Start"
	b.log = b.log.With("op", op)

	ctx := context.Background()

	b.log.Debug("creating new bot")

	bot, err := telego.NewBot(b.botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		b.log.Error("failed to create a new bot", "error", err)
		return err
	}

	b.log.Debug("reciving updates")

	b.updates, err = bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		b.log.Error("failed to recieve updates in chan", "error", err)
		return err
	}

	b.log.Debug("creating bot handler")

	b.h, err = th.NewBotHandler(bot, b.updates)
	if err != nil {
		b.log.Error("failed to create a new bot handler", "error", err)
		return err
	}

	defer func() {
		if err = b.h.Stop(); err != nil {
			b.log.Error("failed to stop handling of updates")
		}
	}()

	b.InitHandlers()

	b.WorkerStart()

	b.log.Debug("starting handling of updates")

	if err := b.h.Start(); err != nil {
		b.log.Error("failed to start handling of updates")
		return err
	}

	return nil
}
