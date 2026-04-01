package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type InternalAPIKeyConfig struct {
	ID     string `json:"ID"`
	Key    string `json:"KEY"`
	Scopes []string `json:"SCOPES"`
}

type DBConfig struct {
	DBURI  string `json:"DB_URI"`
	DBName string `json:"DB_NAME"`
}

type Config struct {
	AppEnv            string                 `json:"APP_ENV"`
	Port              string                 `json:"PORT"`
	JWTSecret         string                 `json:"JWT_SECRET"`
	DB                DBConfig               `json:"DB"`
	InternalAPIKey    []InternalAPIKeyConfig `json:"INTERNAL_API_KEYS"`
	QueueDir          string                 `json:"QUEUE_DIR"`
	QueueTopic        string                 `json:"QUEUE_TOPIC"`
	CronTickSeconds   int                    `json:"CRON_TICK_SECONDS"`
	WorkerPollSeconds int                    `json:"WORKER_POLL_SECONDS"`
	WorkerConsumerID  string                 `json:"WORKER_CONSUMER_ID"`
}

func defaultConfig() Config {
	return Config{
		AppEnv:            "development",
		Port:              "4545",
		JWTSecret:         "dev-secret",
		DB:                DBConfig{DBURI: "mongodb://127.0.0.1:27017", DBName: "p-care"},
		QueueDir:          "/tmp/p-care-queue",
		QueueTopic:        "inventory.sync",
		CronTickSeconds:   30,
		WorkerPollSeconds: 2,
		WorkerConsumerID:  "worker-1",
	}
}

func Load() (Config, error) {
	cfg := defaultConfig()
	
	cleanPath := filepath.Clean("config.json")
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Fallback natively to default configuration 
		}
		return Config{}, fmt.Errorf("failed to read %s: %w", cleanPath, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config inside %s: %w", cleanPath, err)
	}

	return cfg, nil
}
