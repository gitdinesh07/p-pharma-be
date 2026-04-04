package notification

import "context"

type EmailSender interface {
	SendEmail(ctx context.Context, to []string, subject string, body string, isHtml bool) error
}

type SMSSender interface {
	SendSMS(ctx context.Context, to string, message string) error
}
