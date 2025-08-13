package bot

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
	"github.com/mymmrac/telego"
)

type UserService interface {
	GetAllChatIDs(ctx context.Context) (*[]domain.User, error)
}

type PoemService interface {
	GetActualPoem(ctx context.Context) (*domain.Poem, error)
	UpdateStatusSentPoem(ctx context.Context, poemID int64, is_sent bool) error
}

const (
	sleepDuration time.Duration = 10 * time.Second
)

func (b *bot) WorkerStart(userService UserService, poemService PoemService) error {
	for {
		time.Sleep(sleepDuration)
		b.log.Debug("worker started")

		users, err := userService.GetAllChatIDs(context.TODO())
		if err != nil {
			b.log.Error("failed to get all chat ids", "error", err)
			continue
		}
		b.log.Debug("get all users", "users", users)

		poem, err := poemService.GetActualPoem(context.TODO())
		if err != nil {
			b.log.Error("failed to get actual poem", "error", err)
			continue
		}
		b.log.Debug("get poem", "poem", poem)

		err = poemService.UpdateStatusSentPoem(context.TODO(), poem.PoemID, true)

		var wg sync.WaitGroup
		for _, user := range *users {
			wg.Add(1)
			go func() {
				defer wg.Done()

				b.log.Debug("sending poem to chat", "chat_id", user.ChatID)
				b.sendPoemToUsers(telego.ChatID{
					ID: user.ChatID,
				}, *poem)
			}()
		}

		wg.Wait()
	}
}

func (b *bot) sendPoemToUsers(chatID telego.ChatID, poem domain.Poem) error {
	authorName := strings.TrimSpace(fmt.Sprintf("%s %s %s",
		poem.Author.FirstName,
		poem.Author.MiddleName,
		poem.Author.LastName))

	text := fmt.Sprintf(
		"📖 *%s*\n"+
			"✍️ _%s_\n\n"+
			"%s\n\n"+
			"🔍 *Темы / ключевые образы / мотивы:*\n"+
			"%s",
		poem.Name,
		authorName,
		poem.Text,
		poem.Meta,
	)

	_, err := b.b.SendMessage(context.TODO(), &telego.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: telego.ModeMarkdown,
	})

	if err != nil {
		b.log.Error("failed to send poem",
			"chatID", chatID,
			"poemID", poem.PoemID,
			"error", err)
		return fmt.Errorf("send poem: %w", err)
	}

	return nil
}