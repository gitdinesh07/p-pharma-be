package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type InternalAPIKeyConfig struct {
	ID     string
	Key    string
	Scopes []string
}

type DBConfig struct {
	MongoURI    string
	MongoDBName string
}

type Config struct {
	AppEnv            string
	Port              string
	JWTSecret         string
	DB                DBConfig
	InternalAPIKey    []InternalAPIKeyConfig
	QueueDir          string
	QueueTopic        string
	CronTickSeconds   int
	WorkerPollSeconds int
	WorkerConsumerID  string
}

func Load() (Config, error) {
	// Load .env values if file exists. Existing environment variables are not overridden.
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("failed to load .env: %w", err)
	}

	cfg := defaultConfig()

	jsonPath := strings.TrimSpace(getenv("CONFIG_FILE", "config.json"))
	if jsonPath != "" {
		if err := loadJSONConfig(jsonPath, &cfg); err != nil {
			return Config{}, err
		}
	}

	// Environment variables have highest priority and override JSON/default values.
	cfg.AppEnv = getenv("APP_ENV", cfg.AppEnv)
	cfg.Port = getenv("PORT", cfg.Port)
	cfg.JWTSecret = getenv("JWT_SECRET", cfg.JWTSecret)
	if dbRaw := strings.TrimSpace(getenv("DB", "")); dbRaw != "" {
		dbCfg, err := parseDBConfigJSON(dbRaw)
		if err != nil {
			return Config{}, err
		}
		if dbCfg.MongoURI != "" {
			cfg.DB.MongoURI = dbCfg.MongoURI
		}
		if dbCfg.MongoDBName != "" {
			cfg.DB.MongoDBName = dbCfg.MongoDBName
		}
	}
	// Keep these as explicit overrides for environments that prefer flat vars.
	cfg.DB.MongoURI = getenv("MONGO_URI", cfg.DB.MongoURI)
	cfg.DB.MongoDBName = getenv("MONGO_DB_NAME", cfg.DB.MongoDBName)
	cfg.QueueDir = getenv("QUEUE_DIR", cfg.QueueDir)
	cfg.QueueTopic = getenv("QUEUE_TOPIC", cfg.QueueTopic)
	cfg.CronTickSeconds = getenvInt("CRON_TICK_SECONDS", cfg.CronTickSeconds)
	cfg.WorkerPollSeconds = getenvInt("WORKER_POLL_SECONDS", cfg.WorkerPollSeconds)
	cfg.WorkerConsumerID = getenv("WORKER_CONSUMER_ID", cfg.WorkerConsumerID)

	keysRaw := strings.TrimSpace(getenv("INTERNAL_API_KEYS", ""))
	if keysRaw != "" {
		keys, err := parseInternalAPIKeys(keysRaw)
		if err != nil {
			return Config{}, err
		}
		cfg.InternalAPIKey = keys
	}
	return cfg, nil
}

func defaultConfig() Config {
	return Config{
		AppEnv:            "development",
		Port:              "4545",
		JWTSecret:         "dev-secret",
		DB:                DBConfig{MongoURI: "mongodb://127.0.0.1:27017", MongoDBName: "ppharma"},
		QueueDir:          "/tmp/ppharma-queue",
		QueueTopic:        "inventory.sync",
		CronTickSeconds:   30,
		WorkerPollSeconds: 2,
		WorkerConsumerID:  "worker-1",
	}
}

func loadJSONConfig(path string, cfg *Config) error {
	cleanPath := filepath.Clean(path)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read %s: %w", cleanPath, err)
	}

	type legacyFileConfig struct {
		AppEnv            string                 `json:"app_env"`
		Port              string                 `json:"port"`
		JWTSecret         string                 `json:"jwt_secret"`
		MongoURI          string                 `json:"mongo_uri"`
		MongoDBName       string                 `json:"mongo_db_name"`
		InternalAPIKeys   []InternalAPIKeyConfig `json:"internal_api_keys"`
		QueueDir          string                 `json:"queue_dir"`
		QueueTopic        string                 `json:"queue_topic"`
		CronTickSeconds   int                    `json:"cron_tick_seconds"`
		WorkerPollSeconds int                    `json:"worker_poll_seconds"`
		WorkerConsumerID  string                 `json:"worker_consumer_id"`
	}
	type fileDBConfig struct {
		MongoURI    string `json:"MONGO_URI"`
		MongoDBName string `json:"MONGO_DB_NAME"`
	}
	type fileConfig struct {
		AppEnv            string                 `json:"APP_ENV"`
		Port              string                 `json:"PORT"`
		JWTSecret         string                 `json:"JWT_SECRET"`
		DB                fileDBConfig           `json:"DB"`
		InternalAPIKeys   []InternalAPIKeyConfig `json:"INTERNAL_API_KEYS"`
		QueueDir          string                 `json:"QUEUE_DIR"`
		QueueTopic        string                 `json:"QUEUE_TOPIC"`
		CronTickSeconds   int                    `json:"CRON_TICK_SECONDS"`
		WorkerPollSeconds int                    `json:"WORKER_POLL_SECONDS"`
		WorkerConsumerID  string                 `json:"WORKER_CONSUMER_ID"`
	}
	var fc fileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return fmt.Errorf("invalid JSON config in %s: %w", cleanPath, err)
	}

	if fc.AppEnv != "" {
		cfg.AppEnv = fc.AppEnv
	}
	if fc.Port != "" {
		cfg.Port = fc.Port
	}
	if fc.JWTSecret != "" {
		cfg.JWTSecret = fc.JWTSecret
	}
	if fc.DB.MongoURI != "" {
		cfg.DB.MongoURI = fc.DB.MongoURI
	}
	if fc.DB.MongoDBName != "" {
		cfg.DB.MongoDBName = fc.DB.MongoDBName
	}
	if len(fc.InternalAPIKeys) > 0 {
		cfg.InternalAPIKey = fc.InternalAPIKeys
	}
	if fc.QueueDir != "" {
		cfg.QueueDir = fc.QueueDir
	}
	if fc.QueueTopic != "" {
		cfg.QueueTopic = fc.QueueTopic
	}
	if fc.CronTickSeconds > 0 {
		cfg.CronTickSeconds = fc.CronTickSeconds
	}
	if fc.WorkerPollSeconds > 0 {
		cfg.WorkerPollSeconds = fc.WorkerPollSeconds
	}
	if fc.WorkerConsumerID != "" {
		cfg.WorkerConsumerID = fc.WorkerConsumerID
	}

	// Backward compatibility: accept legacy lowercase keys if present.
	var legacy legacyFileConfig
	if err := json.Unmarshal(data, &legacy); err == nil {
		if legacy.AppEnv != "" {
			cfg.AppEnv = legacy.AppEnv
		}
		if legacy.Port != "" {
			cfg.Port = legacy.Port
		}
		if legacy.JWTSecret != "" {
			cfg.JWTSecret = legacy.JWTSecret
		}
		if legacy.MongoURI != "" {
			cfg.DB.MongoURI = legacy.MongoURI
		}
		if legacy.MongoDBName != "" {
			cfg.DB.MongoDBName = legacy.MongoDBName
		}
		if len(legacy.InternalAPIKeys) > 0 {
			cfg.InternalAPIKey = legacy.InternalAPIKeys
		}
		if legacy.QueueDir != "" {
			cfg.QueueDir = legacy.QueueDir
		}
		if legacy.QueueTopic != "" {
			cfg.QueueTopic = legacy.QueueTopic
		}
		if legacy.CronTickSeconds > 0 {
			cfg.CronTickSeconds = legacy.CronTickSeconds
		}
		if legacy.WorkerPollSeconds > 0 {
			cfg.WorkerPollSeconds = legacy.WorkerPollSeconds
		}
		if legacy.WorkerConsumerID != "" {
			cfg.WorkerConsumerID = legacy.WorkerConsumerID
		}
	}
	return nil
}

func parseInternalAPIKeys(keysRaw string) ([]InternalAPIKeyConfig, error) {
	var keys []InternalAPIKeyConfig
	entries := strings.Split(keysRaw, ",")
	for _, e := range entries {
		parts := strings.Split(e, ":")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid INTERNAL_API_KEYS format, expected id:key:scope1|scope2")
		}
		keys = append(keys, InternalAPIKeyConfig{
			ID:     parts[0],
			Key:    parts[1],
			Scopes: strings.Split(parts[2], "|"),
		})
	}
	return keys, nil
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func getenvInt(k string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func parseDBConfigJSON(raw string) (DBConfig, error) {
	type envDBConfig struct {
		MongoURI    string `json:"MONGO_URI"`
		MongoDBName string `json:"MONGO_DB_NAME"`
	}
	var parsed envDBConfig
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return DBConfig{}, fmt.Errorf("invalid DB env JSON: %w", err)
	}
	return DBConfig{
		MongoURI:    parsed.MongoURI,
		MongoDBName: parsed.MongoDBName,
	}, nil
}
