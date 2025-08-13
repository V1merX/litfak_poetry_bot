package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
)

type PoemRepository interface {
	GetActualPoem(ctx context.Context) (*domain.Poem, error)
	UpdateStatusSentPoem(ctx context.Context, poemID int64, is_sent bool) error
}

type PoemService struct {
	rep PoemRepository
	log *slog.Logger
}

func NewPoemService(log *slog.Logger, rep PoemRepository) *PoemService {
	return &PoemService{
		log: log,
		rep: rep,
	}
}

func (s *PoemService) GetActualPoem(ctx context.Context) (*domain.Poem, error) {
	const op = "internal.services.PoemService.GetActualPoem"

	nCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	s.log = s.log.With(
		"op", op,
	)

	poem, err := s.rep.GetActualPoem(nCtx)
	if err != nil {
		return nil, err
	}

	return poem, nil
}

func (s *PoemService) UpdateStatusSentPoem(ctx context.Context, poemID int64, is_sent bool) error {
	const op = "internal.services.PoemService.UpdateStatusSentPoem"

	nCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	s.log = s.log.With(
		"op", op,
	)

	err := s.rep.UpdateStatusSentPoem(nCtx, poemID, is_sent)
	if err != nil {
		return err
	}

	return nil
}
