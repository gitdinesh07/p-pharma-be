package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ppharma/backend/internal/app"
	"ppharma/backend/internal/config"
	cronservice "ppharma/backend/internal/service/cron"
	workerservice "ppharma/backend/internal/service/worker"
	"ppharma/backend/support-pkg/logger/zap"
	"ppharma/backend/support-pkg/queue/filequeue"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	service := os.Getenv("APP_SERVICE")
	if service == "" {
		service = "api"
	}

	switch service {
	case "api":
		runAPI(cfg)
	case "cron":
		runCron(cfg)
	case "worker":
		runWorker(cfg)
	default:
		log.Fatalf("invalid APP_SERVICE=%q (allowed: api|cron|worker)", service)
	}
}

func runAPI(cfg config.Config) {
	a, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("app build error: %v", err)
	}
	if err := a.Engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runCron(cfg config.Config) {
	logger, err := zap.New("debug")
	if err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	queue := filequeue.New(cfg.QueueDir)
	service := cronservice.New(logger, queue, cfg.QueueTopic, time.Duration(cfg.CronTickSeconds)*time.Second)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger.Info(ctx, "cron service started")
	if err := service.Run(ctx); err != nil {
		log.Fatalf("cron service error: %v", err)
	}
}

func runWorker(cfg config.Config) {
	logger, err := zap.New("debug")
	if err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	queue := filequeue.New(cfg.QueueDir)
	service := workerservice.New(logger, queue, cfg.QueueTopic, cfg.WorkerConsumerID, time.Duration(cfg.WorkerPollSeconds)*time.Second)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger.Info(ctx, "worker service started")
	if err := service.Run(ctx); err != nil {
		log.Fatalf("worker service error: %v", err)
	}
}
