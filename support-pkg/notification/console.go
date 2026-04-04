package notification

import (
	"context"
	"log"
)

type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

func (s *ConsoleSender) SendEmail(ctx context.Context, to []string, subject string, body string, isHtml bool) error {
	log.Printf("[EMAIL] To: %s | Subject: %s\nBody: %s\nIsHtml: %t\n", to, subject, body, isHtml)
	return nil
}

func (s *ConsoleSender) SendSMS(ctx context.Context, to string, message string) error {
	log.Printf("[SMS] To: %s | Message: %s\n", to, message)
	return nil
}
