package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

// Enabled implements slog.Handler.
func (d *DiscardHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

// Handle implements slog.Handler.
func (d *DiscardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

// WithAttrs implements slog.Handler.
func (d *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return d
}

// WithGroup implements slog.Handler.
func (d *DiscardHandler) WithGroup(name string) slog.Handler {
	return d
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}
