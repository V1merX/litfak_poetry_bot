package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
)

type UserRepository interface {
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	GetByUserID(ctx context.Context, userID int64) (*domain.User, error)
	NewUser(ctx context.Context, user *domain.User) (int64, error)
	GetAllChatIDs(ctx context.Context) (*[]domain.User, error)
}

type UserService struct {
	rep UserRepository
	log *slog.Logger
}

func NewUserService(log *slog.Logger, rep UserRepository) *UserService {
	return &UserService{
		log: log,
		rep: rep,
	}
}

func (s *UserService) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	const op = "internal.services.UserService.GetByTelegramID"

	nCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	s.log = s.log.With(
		"op", op,
		"telegram_id", telegramID,
	)

	user, err := s.rep.GetByTelegramID(nCtx, telegramID)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *UserService) GetByUserID(ctx context.Context, userID int64) (*domain.User, error) {
	const op = "internal.services.UserService.GetByUserID"

	nCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	s.log = s.log.With(
		"op", op,
		"user_id", userID,
	)

	user, err := s.rep.GetByUserID(nCtx, userID)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *UserService) NewUser(ctx context.Context, user *domain.User) (int64, error) {
	const op = "internal.services.UserService.NewUser"

	nCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	s.log = s.log.With(
		"op", op,
		"telegram_id", user.TelegramID,
		"first_name", user.FirstName,
	)

	userID, err := s.rep.NewUser(nCtx, user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *UserService) GetAllChatIDs(ctx context.Context) (*[]domain.User, error) {
	const op = "internal.services.UserService.GetAllChatIDs"

	nCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	s.log = s.log.With(
		"op", op,
	)

	users, err := s.rep.GetAllChatIDs(nCtx)
	if err != nil {
		return nil, err
	}

	return users, nil
}