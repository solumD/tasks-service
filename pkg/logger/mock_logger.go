package logger

import (
	"context"
	"log/slog"
)

// NewMockLogger возвращает логгер с мок-обработчиком
func NewMockLogger() *slog.Logger {
	return slog.New(NewMockHandler())
}

// MockHandler мок-обработчик
type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *MockHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *MockHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *MockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
