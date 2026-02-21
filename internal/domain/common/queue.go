package common

import (
	"context"
	"time"
)

type QueueMessage struct {
	ID        string    `json:"id"`
	Topic     string    `json:"topic"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Queue interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	ConsumeBatch(ctx context.Context, topic, consumerID string, max int) ([]QueueMessage, error)
}
