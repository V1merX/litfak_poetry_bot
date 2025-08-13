package repositories

import (
	"context"
	"log/slog"
	"time"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	client *pgxpool.Pool
	log    *slog.Logger
}

func NewUserRepository(log *slog.Logger, client *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		client: client,
		log: log,
	}
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	const op = "internal.repositories.UserRepository.GetByTelegramID"

	nCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	r.log = r.log.With(
		"op", op,
		"telegram_id", telegramID,
	)

	var user domain.User

	err := r.client.QueryRow(nCtx, "SELECT * FROM users WHERE telegram_id = $1", telegramID).Scan(&user)
	if err != nil {
		r.log.Error("failed to get user by telegram id", "error", err)
		return nil, err
	}

	return &user, err
}

func (r *UserRepository) GetByUserID(ctx context.Context, userID int64) (*domain.User, error) {
	const op = "internal.repositories.UserRepository.GetByUserID"

	nCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	r.log = r.log.With(
		"op", op,
		"user_id", userID,
	)

	var user domain.User

	err := r.client.QueryRow(nCtx, "SELECT * FROM users WHERE user_id = $1", userID).Scan(&user)
	if err != nil {
		r.log.Error("failed to get user by id", "error", err)
		return nil, err
	}

	return &user, err
}

func (r *UserRepository) NewUser(ctx context.Context, user *domain.User) (int64, error) {
	const op = "internal.repositories.UserRepository.NewUser"

	nCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	r.log = r.log.With(
		"op", op,
		"telegram_id", user.TelegramID,
		"first_name", user.FirstName,
	)

	query := `
		INSERT INTO users (
			telegram_id, 
			chat_id,
			username,
			first_name, 
			last_name, 
			created_at, 
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		ON CONFLICT (telegram_id) 
		DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			updated_at = EXCLUDED.updated_at
		RETURNING user_id
	`

	now := time.Now()
	var userID int64

	err := r.client.QueryRow(nCtx, query,
		user.TelegramID,
		user.ChatID,
		user.UserName,
		user.FirstName,
		user.LastName,
		now,
		now,
	).Scan(&userID)

	if err != nil {
		r.log.Error("failed to create/update user",
			"error", err,
			"query", query,
		)
		return 0, err
	}

	return userID, nil
}

func (r *UserRepository) GetAllChatIDs(ctx context.Context) (*[]domain.User, error) {
	const op = "internal.repositories.UserRepository.GetAll"

	nCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	r.log = r.log.With(
		"op", op,
	)

	users := make([]domain.User, 0, 5000)

	rows, err := r.client.Query(nCtx, "SELECT chat_id FROM users")
	if err != nil {
		r.log.Error("failed to get all users", "error", err)
		return nil, err
	}

	for rows.Next() {
		var chatID int64
		err := rows.Scan(&chatID)
		if err != nil {
			r.log.Error("failed to scan users", "error", err)
			continue
		}
		users = append(users, domain.User{
			ChatID: chatID,
		})
	}

	return &users, err
}