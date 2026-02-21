package worker

import (
	"context"
	"time"

	"ppharma/backend/internal/domain/common"
)

type Service struct {
	log        common.Logger
	queue      common.Queue
	topic      string
	consumerID string
	pollEvery  time.Duration
}

func New(log common.Logger, queue common.Queue, topic, consumerID string, pollEvery time.Duration) *Service {
	if topic == "" {
		topic = "inventory.sync"
	}
	if consumerID == "" {
		consumerID = "default-worker"
	}
	if pollEvery <= 0 {
		pollEvery = 2 * time.Second
	}
	return &Service{log: log, queue: queue, topic: topic, consumerID: consumerID, pollEvery: pollEvery}
}

func (s *Service) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.pollEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info(ctx, "worker service stopped")
			return nil
		case <-ticker.C:
			msgs, err := s.queue.ConsumeBatch(ctx, s.topic, s.consumerID, 10)
			if err != nil {
				s.log.Error(ctx, "failed to consume messages", err, common.Field{Key: "topic", Value: s.topic})
				continue
			}
			for _, msg := range msgs {
				s.log.Info(ctx, "worker processed message",
					common.Field{Key: "topic", Value: msg.Topic},
					common.Field{Key: "message_id", Value: msg.ID},
					common.Field{Key: "payload", Value: string(msg.Payload)},
				)
			}
		}
	}
}
