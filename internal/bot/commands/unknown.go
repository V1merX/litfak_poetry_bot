package commands

import (
	"context"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const (
	unknownCommandMessage = `📚 Ой, кажется, я не знаю такой команды!

Но не переживай! Я помогу тебе подготовиться к ЕГЭ по литературе. Вот что я умею:

✨ Присылать стихотворения для запоминания
📝 Давать задания для закрепления материала
⏳ Напоминать о повторении в оптимальные сроки

Просто начни общение со мной, и я помогу тебе покорить ЕГЭ по литературе! 🚀

Если нужна помощь, напиши /start — я расскажу подробнее.`
)

func Unknown(ctx *th.Context, update telego.Update) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := ctx.Bot().SendMessage(nCtx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		unknownCommandMessage,
	).WithParseMode(telego.ModeMarkdown))

	return err
}
