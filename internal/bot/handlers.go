package bot

import (
	"github.com/V1merX/litfak_poetry_bot/internal/bot/commands"
	th "github.com/mymmrac/telego/telegohandler"
)

func (b *bot) InitHandlers() {
	b.h.Handle(commands.Start, th.CommandEqual("start"))
	b.h.Handle(commands.Unknown, th.AnyCommand())
}
