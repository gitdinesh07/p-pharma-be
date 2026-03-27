package common

import "context"

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}

type Cache interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key, value string)
	Delete(ctx context.Context, key string)
	Exists(ctx context.Context, key string) bool
}

type TokenProvider interface {
	ParseAccessToken(token string) (*Principal, error)
}

type Principal struct {
	ID     string
	Role   string
	Email  string
	Mobile string
}

type APIKeyPrincipal struct {
	KeyID  string
	Scopes map[string]struct{}
}

type APIKeyAuthenticator interface {
	Authenticate(ctx context.Context, rawKey string) (*APIKeyPrincipal, error)
}

type SecretProvider interface {
	InternalAPIKeys(ctx context.Context) ([]InternalAPIKeySecret, error)
}

type InternalAPIKeySecret struct {
	KeyID  string
	RawKey string
	Scopes []string
}
