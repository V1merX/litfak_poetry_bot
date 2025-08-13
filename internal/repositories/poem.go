package repositories

import (
	"context"
	"log/slog"
	"time"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PoemRepository struct {
	client *pgxpool.Pool
	log    *slog.Logger
}

func NewPoemRepository(log *slog.Logger, client *pgxpool.Pool) *PoemRepository {
	return &PoemRepository{
		client: client,
		log:    log,
	}
}

func (r *PoemRepository) GetActualPoem(ctx context.Context) (*domain.Poem, error) {
	const op = "internal.repositories.PoemRepository.GetActualPoem"

	nCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	r.log = r.log.With(
		"op", op,
	)

	var poem domain.Poem

	query := `
        SELECT 
            poems.poem_id, 
            poems.name, 
            poems.author_id, 
            poems.text, 
			poems.meta,
            poems.is_sent, 
            poems.created_at, 
            poems.updated_at,
			authors.first_name,
			authors.middle_name,
			authors.last_name
        FROM poems 
		INNER JOIN authors USING(author_id)
        WHERE is_sent = false 
        ORDER BY created_at ASC 
        LIMIT 1
    `

	err := r.client.QueryRow(nCtx, query).Scan(
		&poem.PoemID,
		&poem.Name,
		&poem.AuthorID,
		&poem.Text,
		&poem.Meta,
		&poem.IsSent,
		&poem.CreatedAt,
		&poem.UpdatedAt,
		&poem.Author.FirstName,
		&poem.Author.MiddleName,
		&poem.Author.LastName,
	)

	if err != nil {
		r.log.Error("failed to get actual poem", "error", err)
		return nil, err
	}

	return &poem, err
}

func (r *PoemRepository) UpdateStatusSentPoem(ctx context.Context, poemID int64, is_sent bool) error {
	const op = "internal.repositories.PoemRepository.UpdateStatusSentPoem"

	nCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	r.log = r.log.With(
		"op", op,
	)

	query := `
        UPDATE 
            poems
        SET 
			is_sent = $1
		WHERE
			poem_id = $2
    `

	err := r.client.QueryRow(nCtx, query, is_sent, poemID).Scan(
		is_sent,
		poemID,
	)

	if err != nil {
		r.log.Error("failed to update status sent poem", "error", err)
		return err
	}

	return nil
}	