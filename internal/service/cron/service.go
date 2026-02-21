package cron

import (
	"context"
	"encoding/json"
	"time"

	"ppharma/backend/internal/domain/common"
)

type Service struct {
	log      common.Logger
	queue    common.Queue
	topic    string
	interval time.Duration
}

func New(log common.Logger, queue common.Queue, topic string, interval time.Duration) *Service {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	if topic == "" {
		topic = "inventory.sync"
	}
	return &Service{log: log, queue: queue, topic: topic, interval: interval}
}

func (s *Service) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info(ctx, "cron service stopped")
			return nil
		case t := <-ticker.C:
			payload, _ := json.Marshal(map[string]any{
				"job":       "inventory.sync",
				"triggered": t.UTC().Format(time.RFC3339),
			})
			if err := s.queue.Publish(ctx, s.topic, payload); err != nil {
				s.log.Error(ctx, "failed to publish cron message", err, common.Field{Key: "topic", Value: s.topic})
				continue
			}
			s.log.Info(ctx, "cron message published", common.Field{Key: "topic", Value: s.topic})
		}
	}
}
