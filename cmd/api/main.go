package main

import (
	"log"

	"ppharma/backend/internal/app"
	"ppharma/backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
	a, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("app build error: %v", err)
	}
	if err := a.Engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
