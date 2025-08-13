package commands

import (
	"context"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type UserService interface {
	GetByUserID(ctx context.Context, userID int64) (*domain.User, error)
	NewUser(ctx context.Context, user *domain.User) (int64, error)
}

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

	errorMessage = `
⚠️ Упс! Что-то пошло не так...

Кажется, у нас небольшая техническая заминка. Но не переживай!

Попробуй:
1. Подождать немного и повторить попытку
2. Проверить соединение с интернетом
3. Написать нам, если проблема повторится

Мы уже работаем над решением! Скоро всё заработает как надо. 💫

А пока можешь перечитать последнее стихотворение или освежить в памяти прошлые задания. 📖

Спасибо за понимание! Терпение — лучший друг литератора. ✨
`
)

func Start(ctx *th.Context, update telego.Update, userService UserService) error {
	_, err := userService.NewUser(ctx, &domain.User{
		TelegramID: update.Message.From.ID,
		ChatID:     update.Message.Chat.ID,
		UserName:   update.Message.From.Username,
		FirstName:  update.Message.From.FirstName,
		LastName:   update.Message.From.LastName,
	})
	if err != nil {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			errorMessage,
		))
		return err
	}

	// TODO: Add context
	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		startMessage,
	))
	return nil
}
