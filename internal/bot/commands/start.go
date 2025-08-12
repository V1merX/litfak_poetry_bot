package commands

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const (
	startMessage = `
	📚 Готовься к ЕГЭ по литературе с умом!
Привет! Здесь ты сможешь легко и эффективно запомнить все нужные стихотворения для ЕГЭ.

Как это работает?
✨ Каждые 4 дня ты получаешь новое стихотворение.
📝 Затем 3 дня подряд — задания к нему, чтобы закрепить материал.

Важно!
Если зашел, а стихотворения еще нет — не паникуй! Возможно, оно уже в пути и появится в течение 3-х дней. Просто наберись терпения — литература любит тех, кто умеет ждать. 😉

Готов покорять ЕГЭ? Тогда начинаем! 🚀
	`
)

func Start(ctx *th.Context, update telego.Update) error {
	// TODO: Add context
	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		startMessage,
	))
	return nil
}
