package zap

import (
	"context"

	"go.uber.org/zap"
	"ppharma/backend/internal/domain/common"
)

type Logger struct {
	base *zap.Logger
}

func New(level string) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	if level == "debug" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{base: l}, nil
}

func (l *Logger) fields(fields []common.Field) []zap.Field {
	out := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		out = append(out, zap.Any(f.Key, f.Value))
	}
	return out
}

func (l *Logger) Debug(_ context.Context, msg string, fields ...common.Field) {
	l.base.Debug(msg, l.fields(fields)...)
}
func (l *Logger) Info(_ context.Context, msg string, fields ...common.Field) {
	l.base.Info(msg, l.fields(fields)...)
}
func (l *Logger) Warn(_ context.Context, msg string, fields ...common.Field) {
	l.base.Warn(msg, l.fields(fields)...)
}
func (l *Logger) Error(_ context.Context, msg string, err error, fields ...common.Field) {
	fields = append(fields, common.Field{Key: "error", Value: err.Error()})
	l.base.Error(msg, l.fields(fields)...)
}
