package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"log/slog"
)

type Config struct {
	Telegram   Telegram   `json:"telegram"`
	PostgreSQL PostgreSQL `json:"postgresql"`
	Logger     Logger     `json:"logger"`
}

type Logger struct {
	Level string `json:"level"`
}

type Telegram struct {
	BotToken string `json:"bot_token"`
}

type PostgreSQL struct {
	Hostname          string        `json:"hostname"`
	Port              int           `json:"port"`
	User              string        `json:"user"`
	Password          string        `json:"password"`
	Database          string        `json:"database"`
	MaxConns          int32         `json:"max_conns"`
	MinConns          int32         `json:"min_conns"`
	MaxConnLifeTime   time.Duration `json:"max_conn_lifetime"`
	MaxConnIdleTime   time.Duration `json:"max_conn_idle_time"`
	HealthCheckPeriod time.Duration `json:"health_check_period"`
}

func (p *PostgreSQL) UnmarshalJSON(data []byte) error {
	type Alias PostgreSQL
	aux := &struct {
		MaxConnLifeTime   string `json:"max_conn_lifetime"`
		MaxConnIdleTime   string `json:"max_conn_idle_time"`
		HealthCheckPeriod string `json:"health_check_period"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if p.MaxConnLifeTime, err = time.ParseDuration(aux.MaxConnLifeTime); err != nil {
		return err
	}
	if p.MaxConnIdleTime, err = time.ParseDuration(aux.MaxConnIdleTime); err != nil {
		return err
	}
	if p.HealthCheckPeriod, err = time.ParseDuration(aux.HealthCheckPeriod); err != nil {
		return err
	}

	return nil
}

func Load(log *slog.Logger, path string) (*Config, error) {
	const op = "internal.config.Load"
	log = log.With("op", op)

	log.Debug("opening config file", "path", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error("failed to open config file", "path", path, "error", err)
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Error("failed to read config file", "path", path, "error", err)
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Error("failed to parse config", "error", err)
		return nil, err
	}

	if err := validateConfig(&cfg); err != nil {
		log.Error("invalid config", "error", err)
		return nil, err
	}

	return &cfg, nil
}

func MustLoad(log *slog.Logger, path string) *Config {
	cfg, err := Load(log, path)
	if err != nil {
		panic(err)
	}
	return cfg
}

func validateConfig(cfg *Config) error {
	if cfg.Telegram.BotToken == "" {
		return errors.New("telegram bot token is required")
	}
	if cfg.PostgreSQL.Hostname == "" {
		return errors.New("postgresql hostname is required")
	}
	if cfg.PostgreSQL.MaxConnLifeTime <= 0 {
		return errors.New("max_conn_lifetime must be positive")
	}
	if cfg.PostgreSQL.MaxConnIdleTime <= 0 {
		return errors.New("max_conn_idle_time must be positive")
	}
	if cfg.PostgreSQL.HealthCheckPeriod <= 0 {
		return errors.New("health_check_period must be positive")
	}
	return nil
}
