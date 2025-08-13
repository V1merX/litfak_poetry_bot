package bot

import (
	"github.com/V1merX/litfak_poetry_bot/internal/bot/commands"
	"github.com/V1merX/litfak_poetry_bot/internal/repositories"
	"github.com/V1merX/litfak_poetry_bot/internal/services"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func (b *bot) InitHandlers() {
	b.log.Info("init routes")
	userRepository := repositories.NewUserRepository(b.log, b.storage)
	userService := services.NewUserService(b.log, userRepository)

	b.h.Handle(
		func(ctx *th.Context, update telego.Update) error {
			return commands.Start(ctx, update, userService)
		},
		th.CommandEqual("start"),
	)
	b.h.Handle(commands.Unknown, th.AnyCommand())
}
