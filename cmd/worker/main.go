package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"ppharma/backend/internal/config"
	workerservice "ppharma/backend/internal/service/worker"
	"ppharma/backend/support-pkg/logger/zap"
	"ppharma/backend/support-pkg/queue/filequeue"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
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
