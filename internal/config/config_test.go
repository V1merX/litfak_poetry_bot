package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

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
				MaxConnLifeTime:   61,
				MaxConnIdleTime:   60,
				HealthCheckPeriod: 23,
			},
			Logger: Logger{
				Level: "debug",
			},
		}

		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		data, err := json.Marshal(expected)
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
