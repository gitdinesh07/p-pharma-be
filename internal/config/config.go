package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type InternalAPIKeyConfig struct {
	ID     string
	Key    string
	Scopes []string
}

type Config struct {
	Port           string
	JWTSecret      string
	InternalAPIKey []InternalAPIKeyConfig
}

func Load() (Config, error) {
	// Load .env values if file exists. Existing environment variables are not overridden.
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("failed to load .env: %w", err)
	}

	cfg := Config{
		Port:      getenv("PORT", "4545"),
		JWTSecret: getenv("JWT_SECRET", "dev-secret"),
	}
	keysRaw := strings.TrimSpace(os.Getenv("INTERNAL_API_KEYS"))
	if keysRaw != "" {
		entries := strings.Split(keysRaw, ",")
		for _, e := range entries {
			parts := strings.Split(e, ":")
			if len(parts) < 3 {
				return Config{}, fmt.Errorf("invalid INTERNAL_API_KEYS format, expected id:key:scope1|scope2")
			}
			cfg.InternalAPIKey = append(cfg.InternalAPIKey, InternalAPIKeyConfig{
				ID:     parts[0],
				Key:    parts[1],
				Scopes: strings.Split(parts[2], "|"),
			})
		}
	}
	return cfg, nil
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
