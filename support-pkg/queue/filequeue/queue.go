package filequeue

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"ppharma/backend/internal/domain/common"
)

type Queue struct {
	baseDir string
	mu      sync.Mutex
}

type fileMessage struct {
	ID        string    `json:"id"`
	Topic     string    `json:"topic"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

func New(baseDir string) *Queue {
	if strings.TrimSpace(baseDir) == "" {
		baseDir = "/tmp/ppharma-queue"
	}
	return &Queue{baseDir: baseDir}
}

func (q *Queue) Publish(_ context.Context, topic string, payload []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if err := os.MkdirAll(q.baseDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(q.baseDir, topic+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	msg := fileMessage{
		ID:        uuid.NewString(),
		Topic:     topic,
		Payload:   base64.StdEncoding.EncodeToString(payload),
		CreatedAt: time.Now().UTC(),
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

func (q *Queue) ConsumeBatch(_ context.Context, topic, consumerID string, max int) ([]common.QueueMessage, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if max <= 0 {
		max = 10
	}
	if err := os.MkdirAll(q.baseDir, 0o755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(q.baseDir, topic+".log")
	offsetPath := filepath.Join(q.baseDir, fmt.Sprintf("%s.%s.offset", topic, consumerID))

	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	offset, err := readOffset(offsetPath)
	if err != nil {
		return nil, err
	}
	if offset >= int64(len(data)) {
		return nil, nil
	}

	lines := strings.Split(string(data[offset:]), "\n")
	msgs := make([]common.QueueMessage, 0, max)
	advanced := offset
	for _, line := range lines {
		advanced += int64(len(line) + 1)
		if strings.TrimSpace(line) == "" {
			continue
		}
		var fm fileMessage
		if err := json.Unmarshal([]byte(line), &fm); err != nil {
			continue
		}
		payload, err := base64.StdEncoding.DecodeString(fm.Payload)
		if err != nil {
			continue
		}
		msgs = append(msgs, common.QueueMessage{
			ID:        fm.ID,
			Topic:     fm.Topic,
			Payload:   payload,
			CreatedAt: fm.CreatedAt,
		})
		if len(msgs) == max {
			break
		}
	}
	if err := os.WriteFile(offsetPath, []byte(strconv.FormatInt(advanced, 10)), 0o644); err != nil {
		return nil, err
	}
	return msgs, nil
}

func readOffset(path string) (int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	v := strings.TrimSpace(string(data))
	if v == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, nil
	}
	if n < 0 {
		return 0, nil
	}
	return n, nil
}
