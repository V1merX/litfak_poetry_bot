package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/V1merX/litfak_poetry_bot/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, log *slog.Logger, dbConfig *config.PostgreSQL) (*pgxpool.Pool, error) {
	var db *pgxpool.Pool
	var pgOnce sync.Once

	connString := strings.TrimSpace(fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d",
		dbConfig.User, dbConfig.Password, dbConfig.Database,
		dbConfig.Hostname, dbConfig.Port))

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error("error parsing pool config", slog.String("error", err.Error()))
		return nil, err
	}

	config.MaxConns = dbConfig.MaxConns
	config.MinConns = dbConfig.MinConns
	config.MaxConnLifetime = dbConfig.MaxConnLifeTime
	config.MaxConnIdleTime = dbConfig.MaxConnIdleTime
	config.HealthCheckPeriod = dbConfig.HealthCheckPeriod

	pgOnce.Do(func() {
		db, err = pgxpool.NewWithConfig(ctx, config)
	})

	if err = db.Ping(ctx); err != nil {
		log.Error("unable to ping database", slog.String("error", err.Error()))
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Info("successfully connected to database")
	return db, nil
}
