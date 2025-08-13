package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gookit/goutil/testutil/assert"
)

func TestLoad(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		expected := &Config{
			Telegram: Telegram{
				BotToken: "default_token",
			},
			PostgreSQL: PostgreSQL{
				Hostname:          "localhost",
				Port:              9999,
				User:              "test_user",
				Password:          "test_password",
				Database:          "test_database",
				MaxConns:          22,
				MinConns:          23,
				MaxConnLifeTime:   61 * time.Nanosecond, // Changed to time.Duration
				MaxConnIdleTime:   60 * time.Nanosecond, // Changed to time.Duration
				HealthCheckPeriod: 23 * time.Nanosecond, // Changed to time.Duration
			},
			Logger: Logger{
				Level: "debug",
			},
		}

		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		// Create a temporary struct that matches the JSON format
		type jsonConfig struct {
			Telegram   Telegram `json:"telegram"`
			PostgreSQL struct {
				Hostname          string `json:"hostname"`
				Port              int    `json:"port"`
				User              string `json:"user"`
				Password          string `json:"password"`
				Database          string `json:"database"`
				MaxConns          int32  `json:"max_conns"`
				MinConns          int32  `json:"min_conns"`
				MaxConnLifeTime   string `json:"max_conn_lifetime"`
				MaxConnIdleTime   string `json:"max_conn_idle_time"`
				HealthCheckPeriod string `json:"health_check_period"`
			} `json:"postgresql"`
			Logger Logger `json:"logger"`
		}

		jsonCfg := jsonConfig{
			Telegram: expected.Telegram,
			PostgreSQL: struct {
				Hostname          string `json:"hostname"`
				Port              int    `json:"port"`
				User              string `json:"user"`
				Password          string `json:"password"`
				Database          string `json:"database"`
				MaxConns          int32  `json:"max_conns"`
				MinConns          int32  `json:"min_conns"`
				MaxConnLifeTime   string `json:"max_conn_lifetime"`
				MaxConnIdleTime   string `json:"max_conn_idle_time"`
				HealthCheckPeriod string `json:"health_check_period"`
			}{
				Hostname:          expected.PostgreSQL.Hostname,
				Port:              expected.PostgreSQL.Port,
				User:              expected.PostgreSQL.User,
				Password:          expected.PostgreSQL.Password,
				Database:          expected.PostgreSQL.Database,
				MaxConns:          expected.PostgreSQL.MaxConns,
				MinConns:          expected.PostgreSQL.MinConns,
				MaxConnLifeTime:   "61ns",
				MaxConnIdleTime:   "60ns",
				HealthCheckPeriod: "23ns",
			},
			Logger: expected.Logger,
		}

		data, err := json.Marshal(jsonCfg)
		assert.NoErr(t, err)
		assert.NoErr(t, os.WriteFile(configPath, data, 0644))

		got, err := Load(slog.Default(), configPath)
		assert.NoErr(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("file not exists", func(t *testing.T) {
		_, err := Load(slog.Default(), "non_existing_file.json")
		assert.Err(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "bad_config.json")

		assert.NoErr(t, os.WriteFile(configPath, []byte("{invalid json}"), 0644))

		_, err := Load(slog.Default(), configPath)
		assert.Err(t, err)
	})
}
